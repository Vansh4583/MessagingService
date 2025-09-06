package api

type SendArgs struct {
	Cap     int
	ToID    string
	Message string
}

type SendReply struct {
	Success bool
}

type ReceiveArgs struct {
	Cap int
}

type ReceiveReply struct {
	From string
	Text string
	Ok   bool
}

type SetSendAllowedArgs struct {
	Cap     int
	Target  string
	Allowed bool
}

type SetSendAllowedReply struct{}

type RemoteReference struct {
	Address string
}

type SetReceiverArgs struct {
	Cap      int
	Receiver RemoteReference
	Receive  bool
}

type SetReceiverReply struct {
	Success bool
}
