package serverStub

import (
	"bytes"
	"encoding/gob"
	"local/auth/types"
	"local/message"
	"local/message/rpc/api"
)

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
		success := message.Send(types.UserCap(args.Cap), args.ToID, args.Message)
		enc.Encode(api.SendReply{Success: success})

	case "Receive":
		var args api.ReceiveArgs
		dec.Decode(&args)
		msg := message.Receive(types.UserCap(args.Cap))
		if msg != nil {
			enc.Encode(api.ReceiveReply{From: msg.From, Text: msg.Text, Ok: true})
		} else {
			enc.Encode(api.ReceiveReply{Ok: false})
		}

	case "SetSendAllowed":
		var args api.SetSendAllowedArgs
		dec.Decode(&args)
		message.SetSendingAllowed(types.UserCap(args.Cap), args.Target, args.Allowed)
		enc.Encode(api.SetSendAllowedReply{})

	case "SetReceiver":
		var args api.SetReceiverArgs
		dec.Decode(&args)
		success := message.SetReceiver(types.UserCap(args.Cap), args.Receiver.Address, args.Receive)
		enc.Encode(api.SetReceiverReply{Success: success})
	}

	return resp.Bytes()
}
