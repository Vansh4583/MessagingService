package serverStub

import (
	"bytes"
	"encoding/gob"
)

var db = map[string]string{}

func Dispatch(req *bytes.Buffer) []byte {
	dec := gob.NewDecoder(req)
	var method string
	dec.Decode(&method)

	var res bytes.Buffer
	enc := gob.NewEncoder(&res)

	switch method {
	case "Put":
		var key, val string
		dec.Decode(&key)
		dec.Decode(&val)
		db[key] = val
		enc.Encode(true)

	case "Get":
		var key string
		dec.Decode(&key)
		val, ok := db[key]
		if !ok {
			enc.Encode("")
		} else {
			enc.Encode(val)
		}
	}
	return res.Bytes()
}
