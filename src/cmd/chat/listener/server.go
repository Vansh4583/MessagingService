//the server stub which is listening for incoming message received calls

package listener

import (
	"bytes"
	"encoding/gob"
	"local/message/types"
)

func Dispatch(req *bytes.Buffer) []byte {
	dec := gob.NewDecoder(req)
	var method string
	dec.Decode(&method)

	var resp bytes.Buffer
	enc := gob.NewEncoder(&resp)

	switch method {
	case "MessageReceived":
		var msg types.Message
		dec.Decode(&msg)

		// Call the local MessageListener
		messageListener := MessageListener(0) // UserCap not needed for display
		messageListener.MessageReceived(msg)

		// Send empty response
		enc.Encode(struct{}{})
	}

	return resp.Bytes()
}
