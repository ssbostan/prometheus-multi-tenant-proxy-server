package server

import (
	"crypto/sha1"
	"crypto/subtle"
	"encoding/hex"
	"log"
	"net"
	"net/http"
)

func authUser(w http.ResponseWriter, r *http.Request) bool {
	username, password, ok := r.BasicAuth()
	user, exists := getUser(username)
	project := r.Header.Get("X-Project-Name")
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		ip = r.RemoteAddr
	}
	if !ok || exists == false {
		return false
	}
	if !authenticate(user, password) || !authorize(user, project) {
		log.Printf("[WARNING] Unauthorized request ip=%s user=%s\n", ip, username)
		return false
	}
	return true
}

func getUser(username string) (UserConfiguration, bool) {
	for _, user := range config.Users {
		if subtle.ConstantTimeCompare([]byte(username), []byte(user.Username)) == 1 {
			return user, true
		}
	}
	return UserConfiguration{}, false
}

func authenticate(user UserConfiguration, password string) bool {
	sha1Hash := sha1.New()
	sha1Hash.Write([]byte(password))
	passwordHash := hex.EncodeToString(sha1Hash.Sum(nil))
	if subtle.ConstantTimeCompare([]byte(passwordHash), []byte(user.Password)) == 1 {
		return true
	}
	return false
}

func authorize(user UserConfiguration, project string) bool {
	for _, p := range user.Projects {
		if subtle.ConstantTimeCompare([]byte(p), []byte(project)) == 1 {
			return true
		}
	}
	return false
}
