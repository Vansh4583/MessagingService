package message

import (
	"local/auth/rpc/clientStub"
	authTypes "local/auth/types"
	"local/message/rpc/proxy"
	"local/message/types"
)

// Registry stores remote references for each user
var registry = make(map[string]map[string]bool) // map[userId]map[address]bool

// SetReceiver registers or unregisters a receiver for push notifications
func SetReceiver(user authTypes.UserCap, receiverAddr string, receive bool) bool {
	id := clientStub.GetId(int(user), authServerAddr)
	if id == "" {
		return false
	}

	receivers := registry[id]
	if receivers == nil {
		receivers = make(map[string]bool)
		registry[id] = receivers
	}

	if receive {
		receivers[receiverAddr] = true
	} else {
		delete(receivers, receiverAddr)
	}
	return true
}

// notifyReceiver sends message to all registered receivers for userId
func notifyReceiver(userId string, msg types.Message) {
	receivers := registry[userId]
	for addr := range receivers {
		proxy := &proxy.ReceiverProxy{Address: addr}
		proxy.MessageReceived(msg)
	}
}
