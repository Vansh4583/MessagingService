package auth

import (
	"math/rand"
	"time"

	"local/auth/types"
)

// Add new user
func Signup(id string, password string) bool {
	if users[id] != nil {
		return false
	} else {
		users[id] = &user{id, password, 0}
		return true
	}
}

// Check credentials and conditionally issue a login cap
func Login(id string, password string) types.UserCap {
	user := users[id]
	if user == nil || user.password != password {
		return 0
	} else {
		if user.cap == 0 {
			user.cap = types.UserCap(rand.Uint64())
		}
		logins[user.cap] = &login{user, now().Add(loginMaxLifetime)}
		return user.cap
	}
}

// Return id associated with cap (if and only if cap is valid and not expired)
func GetId(cap types.UserCap) string {
	user := getUser(cap)
	if user != nil {
		return user.id
	} else {
		return ""
	}
}

// List of valid users
type user struct {
	id       string
	password string
	cap      types.UserCap
}

var users = make(map[string]*user)

// List of logged in users
type login struct {
	user   *user
	expiry time.Time
}

var logins = make(map[types.UserCap]*login)

// Maximum lifetime of a login session
const loginMaxLifetime = time.Hour

// Determine whether cap is still valid (i.e., exists and has not expired)
func isValid(cap types.UserCap) bool {
	login := logins[cap]
	if login == nil {
		return false
	} else if now().After(login.expiry) {
		login.user.cap = 0
		delete(logins, cap)
		return false
	} else {
		return true
	}
}

// Check cap and get user struct if it is valid
func getUser(cap types.UserCap) *user {
	if isValid(cap) {
		login := logins[cap]
		if login != nil {
			return login.user
		}
	}
	return nil
}

// Stubbable time.Now for testing cap expiry
var now = time.Now
