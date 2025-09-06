package clientStub

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"local/auth/rpc/api"
	"local/lib/transport"
)

func Signup(id, pw, authAddr string) bool {

	fmt.Printf("Calling Signup for user: %s\n", id) // Debug line

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode("Signup")
	enc.Encode(api.SignupArgs{ID: id, Password: pw})

	resp, err := transport.Call(&buf, authAddr)
	if err != nil {
		return false
	}

	var reply api.SignupReply
	gob.NewDecoder(resp).Decode(&reply)
	return reply.Success
}

func Login(id, pw, authAddr string) int {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode("Login")
	enc.Encode(api.LoginArgs{ID: id, Password: pw})

	resp, err := transport.Call(&buf, authAddr)
	if err != nil {
		return 0
	}

	var reply api.LoginReply
	gob.NewDecoder(resp).Decode(&reply)
	return reply.Cap
}

func GetId(cap int, authAddr string) string {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode("GetId")
	enc.Encode(api.GetIdArgs{Cap: cap})

	resp, err := transport.Call(&buf, authAddr)
	if err != nil {
		return ""
	}

	var reply api.GetIdReply
	gob.NewDecoder(resp).Decode(&reply)
	return reply.Id
}
