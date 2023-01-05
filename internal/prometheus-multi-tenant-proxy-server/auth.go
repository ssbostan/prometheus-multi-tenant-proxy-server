package server

import (
	"crypto/sha1"
	"crypto/subtle"
	"encoding/hex"
	"net/http"
	"reflect"
)

func checkAccess(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		user := getUser(username)
		if !ok || reflect.ValueOf(user).IsZero() == true {
			unauthorizedAccess(w)
			return
		}
		if !authenticate(user, password) || !authorize(user, r.Header.Get("X-Project-Name")) {
			unauthorizedAccess(w)
			return
		}
		handler(w, r)
	}
}

func unauthorizedAccess(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Basic realm="Prometheus Multi-tenant Proxy Server"`)
	w.WriteHeader(401)
	w.Write([]byte("401 Unauthorized\n"))
}

func getUser(username string) UserConfiguration {
	for _, user := range config.Users {
		if subtle.ConstantTimeCompare([]byte(username), []byte(user.Username)) == 1 {
			return user
		}
	}
	return UserConfiguration{}
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
