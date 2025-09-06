package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	clientStub "local/auth/rpc/clientStub"
	authTypes "local/auth/types"
	"local/cmd/chat/listener"
	"local/lib/finalizer"
	"local/lib/transport"
	"local/message"

	messageClient "local/message/rpc/clientStub"
	"local/message/types"
)

func main() {
	if len(os.Args) != 5 {
		fmt.Println("usage: chat <dbAddr> s|l user password")
		return
	}

	dbAddr := os.Args[1]
	signupOrLogin := os.Args[2]
	userId := os.Args[3]
	password := os.Args[4]

	// Get authd address via DB lookup
	authAddr := getAuthAddr(dbAddr)

	messageClient.Initialize(dbAddr)

	var cap authTypes.UserCap

	// Setup finalizer context
	ctx, cancel := finalizer.WithCancel(context.Background())
	defer func() { cancel(); <-ctx.Done() }()
	finalizer.AfterFunc(ctx, func() {
		if cap != 0 {
			message.SetReceiver(cap, listener.MessageListener(cap), false)
		}
	})

	// Perform signup or login via RPC
	if signupOrLogin == "s" {
		if !clientStub.Signup(userId, password, authAddr) {
			fmt.Println("signup failure")
			return
		}
		fmt.Println("Signup success")
	}

	capInt := clientStub.Login(userId, password, authAddr)
	if capInt == 0 {
		fmt.Println("login failure")
		return
	}
	cap = authTypes.UserCap(capInt)
	fmt.Printf("Login success. Capability: %d\n", cap)

	// Register push receiver so messages can be delivered
	message.SetReceiver(cap, listener.MessageListener(cap), true)

	// Enter chat prompt
	reader := bufio.NewReader(os.Stdin)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			fmt.Print("> ")
			cmd, err := reader.ReadString('\n')
			if err == io.EOF {
				return
			}
			if err != nil {
				panic(err)
			}
			toks := strings.Fields(strings.TrimSpace(cmd))
			if len(toks) == 0 {
				fmt.Println(usage)
				continue
			}
			switch toks[0] {
			case "allow":
				if len(toks) != 2 {
					fmt.Println(usage)
					continue
				}
				messageClient.SetSendingAllowed(cap, toks[1], true)
				fmt.Printf("User %s allowed to send messages\n", toks[1])
			case "block":
				if len(toks) != 2 {
					fmt.Println(usage)
					continue
				}
				messageClient.SetSendingAllowed(cap, toks[1], false)
				fmt.Printf("User %s blocked from sending messages\n", toks[1])
			case "s":
				if len(toks) < 3 {
					fmt.Println(usage)
					continue
				}
				sent := messageClient.Send(cap, toks[1], strings.Join(toks[2:], " "))
				if !sent {
					fmt.Println("Invalid receiver")
				}
			case "read":
				readAll(cap)
			case "push":
				fmt.Println("New messages will be pushed automatically")
				message.SetReceiver(cap, listener.MessageListener(cap), true)
				readAll(cap)
			case "pull":
				message.SetReceiver(cap, listener.MessageListener(cap), false)
				fmt.Println("You must use 'read' to check for new messages")
			default:
				fmt.Println(usage)
			}
		}
	}
}

func getAuthAddr(dbAddr string) string {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode("Get")
	enc.Encode("auth")

	resp, err := transport.Call(&buf, dbAddr)
	if err != nil {
		log.Fatalf("failed to get auth address from DB: %v", err)
	}

	var authAddr string
	if err := gob.NewDecoder(resp).Decode(&authAddr); err != nil {
		log.Fatalf("failed to decode auth address: %v", err)
	}
	return authAddr
}

func readAll(cap authTypes.UserCap) {
	for {
		msg := messageClient.Receive(cap)
		if msg == nil {
			break
		}
		printMsg(msg)
	}
}

func printMsg(msg *types.Message) {
	fmt.Printf("%s: %s\n", msg.From, msg.Text)
}

var usage = "usage: s <user> <message> | read | push | pull | allow <user> | block <user>"
