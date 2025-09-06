package transport

import (
	"bytes"
	"context"
	"fmt"
	"net"
)

type ctxKey string

const udpConnKey ctxKey = "conn"

func WithUDPListenerContext(ctx context.Context, conn *net.UDPConn) context.Context {
	return context.WithValue(ctx, udpConnKey, conn)
}

func getConnFromContext(ctx context.Context) *net.UDPConn {
	val := ctx.Value(udpConnKey)
	if conn, ok := val.(*net.UDPConn); ok {
		return conn
	}
	return nil
}

// Call sends the payload to a server and waits for a response.
func Call(payload *bytes.Buffer, to string) (*bytes.Buffer, error) {
	raddr, err := net.ResolveUDPAddr("udp", to)
	if err != nil {
		return nil, fmt.Errorf("resolve addr: %w", err)
	}

	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return nil, fmt.Errorf("dial UDP: %w", err)
	}
	defer conn.Close()

	// Send payload
	_, err = conn.Write(payload.Bytes())
	if err != nil {
		return nil, fmt.Errorf("write to UDP: %w", err)
	}

	// Wait for response
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("read from UDP: %w", err)
	}

	return bytes.NewBuffer(buf[:n]), nil
}

// Listen starts handling incoming requests using the provided handler
func Listen(ctx context.Context, handler func(msg *bytes.Buffer) []byte) {
	conn := getConnFromContext(ctx)
	if conn == nil {
		panic("Listen: missing UDP connection in context")
	}

	buf := make([]byte, 4096)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			n, addr, err := conn.ReadFromUDP(buf)
			if err != nil {
				continue
			}

			// Dispatch each request in a goroutine
			go func(data []byte, addr *net.UDPAddr) {
				req := bytes.NewBuffer(data)
				resp := handler(req)
				conn.WriteToUDP(resp, addr)
			}(append([]byte(nil), buf[:n]...), addr)
		}
	}
}
