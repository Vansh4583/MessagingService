// auth/rpc/serverStub/server.go
package serverStub

import (
    "bytes"
    "encoding/gob"

    "local/auth"
    authTypes "local/auth/types"
    "local/auth/rpc/api"
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
    case "GetID":
        var args api.GetIDArgs
        dec.Decode(&args)
        id := auth.GetId(authTypes.UserCap(args.Cap))
        enc.Encode(api.GetIDReply{ID: id})
    }

    return resp.Bytes()
}
