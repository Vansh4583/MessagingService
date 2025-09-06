package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"local/db/rpc/api"
	"local/lib/transport"
	"local/message"
	"local/message/rpc/serverStub"
	"log"
	"net"
	"os"
)

func main() {
	ctx := context.Background()
	if len(os.Args) != 2 {
		log.Fatalf("usage: messaged <db_address>")
	}

	dbAddr := os.Args[1]

	// Listen on random UDP port
	addr := "127.0.0.1:0"
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Fatalf("failed to resolve UDP addr: %v", err)
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatalf("failed to listen on UDP: %v", err)
	}

	localAddr := conn.LocalAddr().String()
	fmt.Printf("messaged listening on %s\n", localAddr)

	// Register with DB
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode("Put")
	enc.Encode(api.PutArgs{Key: "messaged", Value: localAddr})

	if _, err := transport.Call(&buf, dbAddr); err != nil {
		log.Fatalf("failed to register messaged: %v", err)
	}
	fmt.Println("messaged registered with db")

	// Get auth server address and set it in message package
	authAddr := getAuthAddr(dbAddr)
	message.SetAuthServerAddr(authAddr)

	// Start listening
	ctx = transport.WithUDPListenerContext(ctx, conn)
	transport.Listen(ctx, serverStub.Dispatch)
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
