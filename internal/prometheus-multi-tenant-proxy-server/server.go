package server

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/urfave/cli/v2"
)

var config Config
var prometheusAddress *url.URL

func Run(c *cli.Context) error {
	config = parseConfigFile(c.String("config"))
	promURL, err := url.Parse(config.Global.PrometheusAddress)
	if err != nil {
		log.Fatal(err)
	}
	prometheusAddress = promURL
	reverseProxy := httputil.NewSingleHostReverseProxy(prometheusAddress)
	http.HandleFunc("/", processRequest(reverseProxy.ServeHTTP))
	if err = http.ListenAndServe(":9999", nil); err != nil {
		log.Fatal(err)
	}
	return nil
}
