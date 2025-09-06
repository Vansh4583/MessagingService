package clientStub

import (
	"bytes"
	"encoding/gob"
	"local/auth/types"
	"local/lib/transport"
	"local/message/rpc/api"
	messageTypes "local/message/types"
)

func Send(cap types.UserCap, to string, text string, messageAddr string) bool {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	enc.Encode("Send")
	enc.Encode(api.SendArgs{Cap: int(cap), ToID: to, Message: text})

	resp, err := transport.Call(&buf, messageAddr)
	if err != nil {
		return false
	}

	var reply api.SendReply
	if err := gob.NewDecoder(resp).Decode(&reply); err != nil {
		return false
	}

	return reply.Success
}

func Receive(cap types.UserCap, messageAddr string) *messageTypes.Message {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	enc.Encode("Receive")
	enc.Encode(api.ReceiveArgs{Cap: int(cap)})

	resp, err := transport.Call(&buf, messageAddr)
	if err != nil {
		return nil
	}

	var reply api.ReceiveReply
	if err := gob.NewDecoder(resp).Decode(&reply); err != nil {
		return nil
	}

	if !reply.Ok {
		return nil
	}

	return &messageTypes.Message{From: reply.From, Text: reply.Text}
}

func SetSendingAllowed(cap types.UserCap, from string, allowed bool, messageAddr string) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	enc.Encode("SetSendAllowed")
	enc.Encode(api.SetSendAllowedArgs{Cap: int(cap), Target: from, Allowed: allowed})

	transport.Call(&buf, messageAddr)
}

func SetReceiver(cap types.UserCap, receiverAddr string, receive bool, messageAddr string) bool {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode("SetReceiver")
	enc.Encode(api.SetReceiverArgs{Cap: int(cap), Receiver: api.RemoteReference{Address: receiverAddr}, Receive: receive})

	resp, err := transport.Call(&buf, messageAddr)
	if err != nil {
		return false
	}

	var reply api.SetReceiverReply
	if err := gob.NewDecoder(resp).Decode(&reply); err != nil {
		return false
	}

	return reply.Success
}
