// auth/rpc/api/api.go
package api

type SignupArgs struct {
    ID       string
    Password string
}

type SignupReply struct {
    Success bool
}

type LoginArgs struct {
    ID       string
    Password string
}

type LoginReply struct {
    Cap int // Placeholder for a real capability type
}

type GetIDArgs struct {
    Cap int
}

type GetIDReply struct {
    ID string
}
