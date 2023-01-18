package server

import (
	"log"
	"net"
	"net/http"
)

func processRequest(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !authUser(w, r) {
			unauthorizedRequest(w)
			return
		}
		if !injectAccessLabel(w, r) {
			badRequest(w)
			return
		}
		logRequest(w, r)
		updateRequest(w, r)
		handler(w, r)
	}
}

func logRequest(w http.ResponseWriter, r *http.Request) {
	username, _, _ := r.BasicAuth()
	access := r.Header.Get(config.Global.AccessRequestHeader)
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		ip = r.RemoteAddr
	}
	log.Printf("[INFO] Accept request ip=%s user=%s access=%s\n", ip, username, access)
}

func updateRequest(w http.ResponseWriter, r *http.Request) {
	r.Host = prometheusAddress.Host
	r.URL.Host = prometheusAddress.Host
	r.URL.Scheme = prometheusAddress.Scheme
	r.Header.Del("Authorization")
}

func unauthorizedRequest(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Basic realm="Prometheus Multi-tenant Proxy Server"`)
	w.WriteHeader(401)
	w.Write([]byte("401 Unauthorized Request\n"))
}

func badRequest(w http.ResponseWriter) {
	w.WriteHeader(400)
	w.Write([]byte("400 Bad Request\n"))
}
