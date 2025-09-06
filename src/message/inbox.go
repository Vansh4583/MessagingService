package message

import (
	//"local/auth"
	clientStub "local/auth/rpc/clientStub"
	authTypes "local/auth/types"
	"local/message/types"
)

// need to store the auth server address - get it from DB during messaged startup
var authServerAddr string

func SetAuthServerAddr(addr string) {
	authServerAddr = addr
}

// Sends a message to the user's inbox if the sender is permitted, returning true if successful
func Send(user authTypes.UserCap, to string, text string) bool {

	// Use auth client stub to validate capability
	id := clientStub.GetId(int(user), authServerAddr)
	if id != "" {
		msg := types.Message{From: id, Text: text}
		el := &inboxElement{&msg, nil}
		inbox := getInbox(to)
		if inbox.validSenders[id] {
			if inbox.tail == nil {
				inbox.head = el
				inbox.tail = el
			} else {
				inbox.tail.next = el
				inbox.tail = el
			}
			notifyReceiver(to, msg)
			return true
		}
	}
	return false
}

// Dequeue and return next message in user's inbox fifo
func Receive(user authTypes.UserCap) *types.Message {

	// Use auth client stub to validate capability
	userId := clientStub.GetId(int(user), authServerAddr)
	if userId != "" {
		inbox := inboxes[userId]
		if inbox == nil || inbox.head == nil {
			return nil
		}
		el := inbox.head
		inbox.head = inbox.head.next
		if inbox.head == nil {
			inbox.tail = nil
		}
		return el.message
	}
	return nil
}

// Update whether the user accepts messages from the specified sender
func SetSendingAllowed(user authTypes.UserCap, from string, allowed bool) {

	// Use auth client stub to validate capability
	id := clientStub.GetId(int(user), authServerAddr)
	if id != "" {
		inbox := getInbox(id)
		if allowed {
			inbox.validSenders[from] = true
		} else {
			delete(inbox.validSenders, from)
		}
	}
}

// Element of inbox linked list
type inboxElement struct {
	message *types.Message
	next    *inboxElement
}

// Descriptor for inbox linked list
type inboxDesc struct {
	validSenders map[string]bool
	head         *inboxElement
	tail         *inboxElement
}

// Map of linked lists
var inboxes = make(map[string]*inboxDesc)

func getInbox(id string) *inboxDesc {
	inbox := inboxes[id]
	if inbox == nil {
		inbox = &inboxDesc{make(map[string]bool), nil, nil}
		inboxes[id] = inbox
	}
	return inbox
}
