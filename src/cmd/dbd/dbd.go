package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"net"
	"sync"

	"local/db/rpc/api"
	"local/lib/transport"
)

var store = make(map[string]string)
var mu sync.Mutex

func main() {
	ctx := context.Background()

	addr := "127.0.0.1:0"
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		panic(err)
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		panic(err)
	}
	fmt.Printf("listening on %s\n", conn.LocalAddr().String())

	ctx = transport.WithUDPListenerContext(ctx, conn)
	transport.Listen(ctx, handleRequest)
}

func handleRequest(req *bytes.Buffer) []byte {
	dec := gob.NewDecoder(req)
	encBuf := new(bytes.Buffer)
	enc := gob.NewEncoder(encBuf)

	var method string
	if err := dec.Decode(&method); err != nil {
		fmt.Println("error decoding method:", err)
		enc.Encode("error: bad method")
		return encBuf.Bytes()
	}

	switch method {
	case "Put":
		var args api.PutArgs
		if err := dec.Decode(&args); err != nil {
			fmt.Println("error decoding PutArgs:", err)
			enc.Encode("error decoding args")
			return encBuf.Bytes()
		}
		mu.Lock()
		store[args.Key] = args.Value
		mu.Unlock()

		fmt.Printf("Put request: key=%s, value=%s\n", args.Key, args.Value)
		enc.Encode(struct{}{}) // Acknowledge with empty struct

	case "Get":
		var key string
		if err := dec.Decode(&key); err != nil {
			enc.Encode("error decoding key")
			return encBuf.Bytes()
		}
		mu.Lock()
		value := store[key]
		mu.Unlock()
		enc.Encode(value)

	default:
		fmt.Println("unknown method:", method)
		enc.Encode("unknown method")
	}
	return encBuf.Bytes()
}
