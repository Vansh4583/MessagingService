package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"local/db/rpc/api"
	"local/lib/transport"
	serverStub "local/message/rpc/serverStub"
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

	serverStub.Initialize(dbAddr)

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

	// Register with DB using STRUCT-BASED RPC
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode("Put"); err != nil {
		log.Fatalf("failed to encode method: %v", err)
	}

	// USE STRUCT INSTEAD OF SEPARATE STRINGS
	if err := enc.Encode(api.PutArgs{Key: "messaged", Value: localAddr}); err != nil {
		log.Fatalf("failed to encode PutArgs: %v", err)
	}

	if _, err := transport.Call(&buf, dbAddr); err != nil {
		log.Fatalf("failed to register messaged: %v", err)
	}

	fmt.Println("messaged registered with db")

	// Add connection to context and listen for requests
	ctx = transport.WithUDPListenerContext(ctx, conn)
	transport.Listen(ctx, serverStub.Dispatch)
}
