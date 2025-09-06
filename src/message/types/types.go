package types

// Message
type Message struct {
	From string
	Text string
}

// Notication receiver
type Receiver interface {
	MessageReceived(Message)
}
