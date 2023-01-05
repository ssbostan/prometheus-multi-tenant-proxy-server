package server

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/urfave/cli/v2"
)

var config Config

func Run(c *cli.Context) error {
	config = parseConfigFile(c.String("config"))
	prometheusAddress, err := url.Parse(config.Global.PrometheusAddress)
	if err != nil {
		log.Fatal(err)
	}
	reverseProxy := httputil.NewSingleHostReverseProxy(prometheusAddress)
	originalDirector := reverseProxy.Director
	reverseProxy.Director = func(r *http.Request) {
		originalDirector(r)
		injectProjectLabel(r)
		func(r *http.Request) {
			r.Host = prometheusAddress.Host
			r.URL.Host = prometheusAddress.Host
			r.URL.Scheme = prometheusAddress.Scheme
			r.Header.Del("Authorization")
		}(r)
	}
	http.HandleFunc("/", logRequest(checkAccess(reverseProxy.ServeHTTP)))
	if err = http.ListenAndServe(":9999", nil); err != nil {
		log.Fatal(err)
	}
	return nil
}
