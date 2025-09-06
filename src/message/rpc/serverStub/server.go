package serverStub

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"

	authAPI "local/auth/rpc/api"
	"local/lib/transport"
	"local/message/rpc/api"
	messageTypes "local/message/types"
)

var dbAddr string // Store DB address for auth server lookup

// Define inbox types locally to avoid import issues
type inboxElement struct {
	message *messageTypes.Message
	next    *inboxElement
}

type inboxDesc struct {
	validSenders map[string]bool
	head         *inboxElement
	tail         *inboxElement
}

var inboxes = make(map[string]*inboxDesc)

func getInbox(id string) *inboxDesc {
	inbox := inboxes[id]
	if inbox == nil {
		inbox = &inboxDesc{validSenders: make(map[string]bool)}
		inboxes[id] = inbox
	}
	return inbox
}

// Initialize the server stub with DB address
func Initialize(dbAddress string) {
	dbAddr = dbAddress
}

// Get auth server address from DB
func getAuthAddr() string {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode("Get")
	enc.Encode("auth")

	resp, err := transport.Call(&buf, dbAddr)
	if err != nil {
		log.Printf("Failed to get auth address: %v", err)
		return ""
	}

	var authAddr string
	if err := gob.NewDecoder(resp).Decode(&authAddr); err != nil {
		log.Printf("Failed to decode auth address: %v", err)
		return ""
	}

	return authAddr
}

// Get user ID via RPC call to auth server
func getUserID(cap int) string {
	authAddr := getAuthAddr()
	if authAddr == "" {
		return ""
	}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode("GetID")
	enc.Encode(authAPI.GetIDArgs{Cap: cap})

	resp, err := transport.Call(&buf, authAddr)
	if err != nil {
		log.Printf("Failed to get user ID: %v", err)
		return ""
	}

	var reply authAPI.GetIDReply
	if err := gob.NewDecoder(resp).Decode(&reply); err != nil {
		log.Printf("Failed to decode GetID reply: %v", err)
		return ""
	}

	fmt.Printf("DEBUG: Got user ID '%s' for cap %d\n", reply.ID, cap)
	return reply.ID
}

// Handle sending messages directly without calling message.Send
func sendMessageDirect(fromID, toID, text string) bool {
	if fromID == "" {
		return false
	}

	fmt.Printf("DEBUG: SERVER - Checking send permission from '%s' to '%s'\n", fromID, toID)

	msg := &messageTypes.Message{From: fromID, Text: text}

	inbox := getInbox(toID)
	fmt.Printf("DEBUG: SERVER - ValidSenders for '%s': %v\n", toID, inbox.validSenders)

	if inbox.validSenders[fromID] {
		el := &inboxElement{message: msg}
		if inbox.tail == nil {
			inbox.head = el
			inbox.tail = el
		} else {
			inbox.tail.next = el
			inbox.tail = el
		}
		fmt.Printf("DEBUG: Message sent from '%s' to '%s'\n", fromID, toID)
		return true
	}

	fmt.Printf("DEBUG: Send failed - '%s' not allowed to send to '%s'\n", fromID, toID)
	return false
}

func Dispatch(req *bytes.Buffer) []byte {
	dec := gob.NewDecoder(req)
	var method string
	dec.Decode(&method)

	var resp bytes.Buffer
	enc := gob.NewEncoder(&resp)

	switch method {
	case "Send":
		var args api.SendArgs
		dec.Decode(&args)

		senderID := getUserID(args.Cap)
		if senderID == "" {
			fmt.Printf("DEBUG: Failed to get sender ID for cap %d\n", args.Cap)
			enc.Encode(api.SendReply{Success: false})
			return resp.Bytes()
		}

		success := sendMessageDirect(senderID, args.ToID, args.Message)
		enc.Encode(api.SendReply{Success: success})

	case "Receive":
		var args api.ReceiveArgs
		dec.Decode(&args)

		userID := getUserID(args.Cap)
		if userID == "" {
			enc.Encode(api.ReceiveReply{Ok: false})
			return resp.Bytes()
		}

		inbox := getInbox(userID)
		if inbox.head == nil {
			enc.Encode(api.ReceiveReply{Ok: false})
			return resp.Bytes()
		}

		el := inbox.head
		inbox.head = el.next
		if inbox.head == nil {
			inbox.tail = nil
		}

		enc.Encode(api.ReceiveReply{From: el.message.From, Text: el.message.Text, Ok: true})

	case "SetSendAllowed":
		var args api.SetSendAllowedArgs
		dec.Decode(&args)

		receiverID := getUserID(args.Cap)
		if receiverID == "" {
			enc.Encode(api.SetSendAllowedReply{})
			return resp.Bytes()
		}

		fmt.Printf("DEBUG: SERVER - SetSendAllowed for receiver '%s', allowing '%s'\n", receiverID, args.Target)

		inbox := getInbox(receiverID)
		if args.Allowed {
			inbox.validSenders[args.Target] = true
			fmt.Printf("DEBUG: Allowed '%s' to send to '%s'\n", args.Target, receiverID)
			fmt.Printf("DEBUG: SERVER - Permission stored. ValidSenders for '%s': %v\n", receiverID, inbox.validSenders)
		} else {
			delete(inbox.validSenders, args.Target)
			fmt.Printf("DEBUG: Blocked '%s' from sending to '%s'\n", args.Target, receiverID)
		}

		enc.Encode(api.SetSendAllowedReply{})
	}

	return resp.Bytes()
}
