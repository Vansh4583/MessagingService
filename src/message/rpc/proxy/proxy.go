//nothing

package proxy

import (
	"bytes"
	"encoding/gob"
	"local/lib/transport"
	"local/message/types"
)

type ReceiverProxy struct {
	Address string
}

func (p *ReceiverProxy) MessageReceived(msg types.Message) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode("MessageReceived")
	enc.Encode(msg)

	// Send RPC to chat client address
	transport.Call(&buf, p.Address)
}
