package clientStub

import (
	"bytes"
	"encoding/gob"
	"fmt"
	authTypes "local/auth/types"
	"local/lib/transport"
	"local/message/rpc/api"
	"local/message/types"
	"log"
)

var messageServerAddr string

func Initialize(dbAddr string) {
	// Get message server address from DB
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode("Get")
	enc.Encode("messaged")

	resp, err := transport.Call(&buf, dbAddr)
	if err != nil {
		log.Fatalf("failed to get messaged address: %v", err)
	}

	if err := gob.NewDecoder(resp).Decode(&messageServerAddr); err != nil {
		log.Fatalf("failed to decode messaged address: %v", err)
	}
	fmt.Printf("DEBUG: messageServerAddr set to %s\n", messageServerAddr)
}

func Send(cap authTypes.UserCap, to string, text string) bool {
	fmt.Printf("DEBUG: Send called - cap: %d, to: %s, text: %s\n", cap, to, text)

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode("Send")
	enc.Encode(api.SendArgs{Cap: int(cap), ToID: to, Message: text})

	resp, err := transport.Call(&buf, messageServerAddr)
	if err != nil {
		fmt.Printf("ERROR: Send RPC failed: %v\n", err)
		return false
	}

	var reply api.SendReply
	if err := gob.NewDecoder(resp).Decode(&reply); err != nil {
		fmt.Printf("ERROR: Failed to decode Send reply: %v\n", err)
		return false
	}

	fmt.Printf("DEBUG: Send reply success: %t\n", reply.Success)
	return reply.Success
}

func Receive(cap authTypes.UserCap) *types.Message {
	fmt.Printf("DEBUG: Receive called - cap: %d\n", cap)

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode("Receive")
	enc.Encode(api.ReceiveArgs{Cap: int(cap)})

	resp, err := transport.Call(&buf, messageServerAddr)
	if err != nil {
		fmt.Printf("ERROR: Receive RPC failed: %v\n", err)
		return nil
	}

	var reply api.ReceiveReply
	if err := gob.NewDecoder(resp).Decode(&reply); err != nil {
		fmt.Printf("ERROR: Failed to decode Receive reply: %v\n", err)
		return nil
	}

	if !reply.Ok {
		fmt.Printf("DEBUG: Receive reply Ok is false, no message\n")
		return nil
	}

	fmt.Printf("DEBUG: Receive reply From: %s, Text: %s\n", reply.From, reply.Text)
	return &types.Message{From: reply.From, Text: reply.Text}
}

func SetSendingAllowed(cap authTypes.UserCap, from string, allowed bool) {
	fmt.Printf("DEBUG: SetSendingAllowed called - cap: %d, from: %s, allowed: %t\n", cap, from, allowed)
	fmt.Printf("DEBUG: messageServerAddr: %s\n", messageServerAddr)

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode("SetSendAllowed")
	enc.Encode(api.SetSendAllowedArgs{Cap: int(cap), Target: from, Allowed: allowed})

	resp, err := transport.Call(&buf, messageServerAddr)
	if err != nil {
		fmt.Printf("ERROR: SetSendingAllowed RPC failed: %v\n", err)
		return
	}

	var reply api.SetSendAllowedReply
	if err := gob.NewDecoder(resp).Decode(&reply); err != nil {
		fmt.Printf("ERROR: Failed to decode SetSendingAllowed reply: %v\n", err)
	}

	fmt.Printf("DEBUG: SetSendingAllowed completed successfully\n")
}
