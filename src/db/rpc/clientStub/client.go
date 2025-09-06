// db/rpc/clientStub/client.go
package clientStub

import (
    "bytes"
    "encoding/gob"

    "local/db/rpc/api"
    "local/lib/transport"
)

func Put(key, value string, dbAddr string) error {
    var buf bytes.Buffer
    enc := gob.NewEncoder(&buf)
    enc.Encode("Put")
    enc.Encode(api.PutArgs{Key: key, Value: value})

    _, err := transport.Call(&buf, dbAddr)
    return err
}
