package auth

import (
	"testing"
	"time"
)

func TestSimple(t *testing.T) {
	Signup("u1", "p1")
	if Login("u1", "p2") != 0 {
		t.Fatal("password check")
	}
	cap := Login("u1", "p1")
	if cap == 0 {
		t.Fatal("login")
	}
	if GetId(cap) == "" {
		t.Fatal("says valid cap is invalid")
	}
	if GetId(0) != "" {
		t.Fatal("says invalid cap is valid")
	}
	Signup("u2", "p2")
	now = func() time.Time { return time.Now().Add(loginMaxLifetime + time.Millisecond) }
	if GetId(cap) != "" {
		t.Fatal("expired cap still valid")
	}
	now = time.Now
}
