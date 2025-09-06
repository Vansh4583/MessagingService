package message

import (
	"local/auth"
	authTypes "local/auth/types"
	"local/message/types"
)

// Register a notification receiver for user identified by userCap
func SetReceiver(user authTypes.UserCap, receiver types.Receiver, receive bool) {
	id := auth.GetId(user)
	if id != "" {
		receivers := registry[id]
		if receivers == nil {
			receivers = make(map[types.Receiver]bool)
			registry[id] = receivers
		}
		if receive {
			receivers[receiver] = true
		} else {
			delete(receivers, receiver)
		}
	}
}

// List of registered registry
var registry = make(map[string]map[types.Receiver]bool)

// Notify registered receiver (if there is one) for user identified by userId
func notifyReceiver(userId string, msg types.Message) {
	receivers := registry[userId]
	for receiver := range receivers {
		receiver.MessageReceived(msg)
	}
}
