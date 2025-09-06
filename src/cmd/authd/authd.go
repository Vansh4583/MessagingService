package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"local/auth/rpc/serverStub"
	"local/db/rpc/clientStub"
	"local/lib/transport"
)

func main() {
	ctx := context.Background()

	if len(os.Args) != 2 {
		log.Fatal("usage: authd <db address>")
	}
	dbAddr := os.Args[1]

	// Start listening on a random UDP port
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
	fmt.Printf("authd listening on %s\n", localAddr)
	fmt.Printf("Registering authd with db at %s as %s\n", dbAddr, localAddr)

	// Register authd address in db
	err = clientStub.Put("auth", localAddr, dbAddr)
	if err != nil {
		log.Fatalf("failed to register auth service: %v", err)
	}
	fmt.Println("authd registered with db")

	// Add conn to context and start listening
	ctx = transport.WithUDPListenerContext(ctx, conn)
	transport.Listen(ctx, serverStub.Dispatch)
}
