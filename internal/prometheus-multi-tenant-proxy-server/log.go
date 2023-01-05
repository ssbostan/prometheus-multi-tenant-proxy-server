package server

import (
	"log"
	"net"
	"net/http"
)

func logRequest(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, _, _ := r.BasicAuth()
		project := r.Header.Get("X-Project-Name")
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}
		log.Printf("[INFO] Request from IP=%s User=%s Project=%s\n", ip, username, project)
		handler(w, r)
	}
}
