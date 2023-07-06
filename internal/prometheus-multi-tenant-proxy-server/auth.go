package server

import (
	"crypto/sha1"
	"crypto/subtle"
	"encoding/hex"
	"log"
	"net"
	"net/http"
)

func authUser(_ http.ResponseWriter, r *http.Request) bool {
	username, password, ok := r.BasicAuth()
	user, exists := getUser(username)
	access := r.Header.Get(config.Global.AccessRequestHeader)
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		ip = r.RemoteAddr
	}
	if !ok || exists == false {
		return false
	}
	if !authenticate(user, password) || !authorize(user, access) {
		log.Printf("[WARNING] Unauthorized request ip=%s user=%s\n", ip, username)
		return false
	}
	return true
}

func getUser(username string) (UserConfig, bool) {
	for _, user := range config.Users {
		if subtle.ConstantTimeCompare([]byte(username), []byte(user.Username)) == 1 {
			return user, true
		}
	}
	return UserConfig{}, false
}

func authenticate(user UserConfig, password string) bool {
	sha1Hash := sha1.New()
	sha1Hash.Write([]byte(password))
	passwordHash := hex.EncodeToString(sha1Hash.Sum(nil))
	if subtle.ConstantTimeCompare([]byte(passwordHash), []byte(user.Password)) == 1 {
		return true
	}
	return false
}

func authorize(user UserConfig, access string) bool {
	for _, a := range user.Accesses {
		if subtle.ConstantTimeCompare([]byte(a), []byte(access)) == 1 {
			return true
		}
	}
	return false
}
