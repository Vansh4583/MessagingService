package message

import (
	"fmt"
	"testing"

	"local/auth"
)

func TestMain(t *testing.T) {
	TestSimple(t)
	TestMultiSend(t)
}

func TestSimple(t *testing.T) {
	auth.Signup("u1", "p1")
	auth.Signup("u2", "p2")
	cap1 := auth.Login("u1", "p1")
	cap2 := auth.Login("u2", "p2")
	if Send(cap2, "u1", "b") {
		t.Fatal("Disallowed send permitted")
	}
	SetSendingAllowed(cap1, "u2", true)
	if !Send(cap2, "u1", "b") {
		t.Fatal("Allowed send not permitted")
	}
	m := Receive(cap1)
	if m == nil {
		t.Fatal(("no receive"))
	}
	if m.From != "u2" || m.Text != "b" {
		t.Fatalf("Wrong receive %s %s\n", m.From, m.Text)
	}
	Send(cap2, "u1", "b")
	Receive(cap1)
}

func TestMultiSend(t *testing.T) {
	auth.Signup("u1", "p1")
	auth.Signup("u2", "p2")
	cap1 := auth.Login("u1", "p1")
	cap2 := auth.Login("u2", "p2")
	SetSendingAllowed(cap1, "u2", true)
	for i := range 10 {
		Send(cap2, "u1", fmt.Sprintf("b%d", i))
	}
	for i := range 10 {
		m := Receive(cap1)
		if m == nil {
			t.Fatalf("recieve fails at %d\n", i)
		}
		if m.Text != fmt.Sprintf("b%d", i) {
			t.Fatalf("wrong receive %s at %d\n", m.Text, i)
		}
	}
	if Receive(cap1) != nil {
		t.Fatalf("extra receive")
	}
	if inboxes["u1"].head != nil {
		t.Fatal("inbox memory leak")
	}
}
