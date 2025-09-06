package message

import (
	"testing"
	"time"

	"local/auth"
	"local/message/types"
)

func TestReceiver(t *testing.T) {
	auth.Signup("u1", "p1")
	auth.Signup("u2", "p2")
	c1 := auth.Login("u1", "p1")
	c2 := auth.Login("u2", "p2")
	SetSendingAllowed(c1, "u2", true)
	ch := make(receiver, 1)
	SetReceiver(c1, ch, true)
	Send(c2, "u1", "hello")
	go func() {
		time.Sleep(time.Millisecond)
		ch <- "not notified"
	}()
	note := <-ch
	if note != "received u2 hello" {
		t.Fatal("notification failure")
	}
	defer close(ch)
}

type receiver chan string

func (ch receiver) MessageReceived(msg types.Message) {
	ch <- "received " + msg.From + " " + msg.Text
}
