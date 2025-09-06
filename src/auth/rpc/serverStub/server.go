// auth/rpc/serverStub/server.go
package serverStub

import (
	"bytes"
	"encoding/gob"

	"local/auth"
	"local/auth/rpc/api"
	authTypes "local/auth/types"
)

func Dispatch(req *bytes.Buffer) []byte {
	dec := gob.NewDecoder(req)
	var method string
	dec.Decode(&method)

	var resp bytes.Buffer
	enc := gob.NewEncoder(&resp)

	switch method {
	case "Signup":
		var args api.SignupArgs
		dec.Decode(&args)
		result := auth.Signup(args.ID, args.Password)
		enc.Encode(api.SignupReply{Success: result})

	case "Login":
		var args api.LoginArgs
		dec.Decode(&args)
		cap := auth.Login(args.ID, args.Password)
		enc.Encode(api.LoginReply{Cap: int(cap)})

	case "GetId": // Changed from "GetID" to "GetId"
		var args api.GetIdArgs // Changed from GetIDArgs to GetIdArgs
		dec.Decode(&args)
		id := auth.GetId(authTypes.UserCap(args.Cap))
		enc.Encode(api.GetIdReply{Id: id}) // Changed from GetIDReply to GetIdReply, ID to Id
	}

	return resp.Bytes()
}
