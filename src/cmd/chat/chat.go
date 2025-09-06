package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"io"
	authClientStub "local/auth/rpc/clientStub"
	authTypes "local/auth/types"
	messageClientStub "local/message/rpc/clientStub"
	"log"
	"net"
	"os"
	"strings"

	"local/cmd/chat/listener"
	"local/lib/finalizer"
	"local/lib/transport"
	"local/message/types"
)

var chatServerAddr string // Global variable to store chat server address

func main() {
	if len(os.Args) != 5 {
		fmt.Println("usage: chat <db_address> s|l <user> <password>")
		return
	}

	dbAddr := os.Args[1]
	signupOrLogin := os.Args[2]
	userId := os.Args[3]
	password := os.Args[4]

	// Get server addresses
	authAddr := getAuthAddr(dbAddr)
	messageAddr := getMessageAddr(dbAddr)

	var cap authTypes.UserCap

	// Setup finalizer context
	ctx, cancel := finalizer.WithCancel(context.Background())
	defer func() { cancel(); <-ctx.Done() }()

	// Perform signup or login
	if signupOrLogin == "s" {
		if !authClientStub.Signup(userId, password, authAddr) {
			fmt.Println("signup failure")
			return
		}
		fmt.Println("Signup success")
	}

	capInt := authClientStub.Login(userId, password, authAddr)
	if capInt == 0 {
		fmt.Println("login failure")
		return
	}

	cap = authTypes.UserCap(capInt)
	fmt.Printf("Login success. Capability: %d\n", cap)

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
				messageClientStub.SetSendingAllowed(cap, toks[1], true, messageAddr)
				fmt.Printf("User %s allowed to send messages\n", toks[1])

			case "block":
				if len(toks) != 2 {
					fmt.Println(usage)
					continue
				}
				messageClientStub.SetSendingAllowed(cap, toks[1], false, messageAddr)
				fmt.Printf("User %s blocked from sending messages\n", toks[1])

			case "s":
				if len(toks) < 3 {
					fmt.Println(usage)
					continue
				}
				sent := messageClientStub.Send(cap, toks[1], strings.Join(toks[2:], " "), messageAddr)
				if !sent {
					fmt.Println("Failed to send message")
				}

			case "read":
				readAll(cap, messageAddr)

			case "notify":
				if chatServerAddr == "" {
					startChatServer(ctx)
				}
				success := messageClientStub.SetReceiver(cap, chatServerAddr, true, messageAddr)
				if success {
					fmt.Println("Push notifications enabled")
				} else {
					fmt.Println("Failed to enable push notifications")
				}

			default:
				fmt.Println(usage)
			}
		}
	}
}

func startChatServer(ctx context.Context) {
	// Start UDP server for receiving push notifications
	addr := "127.0.0.1:0"
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Fatalf("failed to resolve UDP addr: %v", err)
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatalf("failed to listen on UDP: %v", err)
	}

	chatServerAddr = conn.LocalAddr().String()
	fmt.Printf("Chat server listening on %s\n", chatServerAddr)

	// Start listening for incoming RPCs in a goroutine
	go func() {
		ctx = transport.WithUDPListenerContext(ctx, conn)
		transport.Listen(ctx, listener.Dispatch)
	}()
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

func getMessageAddr(dbAddr string) string {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode("Get")
	enc.Encode("messaged")

	resp, err := transport.Call(&buf, dbAddr)
	if err != nil {
		log.Fatalf("failed to get message address from DB: %v", err)
	}

	var messageAddr string
	if err := gob.NewDecoder(resp).Decode(&messageAddr); err != nil {
		log.Fatalf("failed to decode message address: %v", err)
	}

	return messageAddr
}

func readAll(cap authTypes.UserCap, messageAddr string) {
	for {
		msg := messageClientStub.Receive(cap, messageAddr)
		if msg == nil {
			break
		}
		printMsg(msg)
	}
}

func printMsg(msg *types.Message) {
	fmt.Printf("%s: %s\n", msg.From, msg.Text)
}

var usage = "usage: s <user> <message> | read | allow <user> | block <user> | notify"
