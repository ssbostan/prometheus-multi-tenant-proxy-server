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
	log.Printf("[INFO] Config loaded from: %s\n", c.String("config"))
	log.Printf("[INFO] Access request header: %s\n", config.Global.AccessRequestHeader)
	log.Printf("[INFO] Access target label: %s\n", config.Global.AccessTargetLabel)
	promAddress, err := url.Parse(config.Global.PrometheusAddress)
	if err != nil {
		log.Fatal(err)
	}
	prometheusAddress = promAddress
	log.Printf("[INFO] Prometheus server address: %s\n", prometheusAddress.String())
	log.Printf("[INFO] Proxy server address: %s\n", config.Global.ListenAddress)
	reverseProxy := httputil.NewSingleHostReverseProxy(prometheusAddress)
	http.HandleFunc("/", processRequest(reverseProxy.ServeHTTP))
	if err = http.ListenAndServe(config.Global.ListenAddress, nil); err != nil {
		log.Fatal(err)
	}
	return nil
}
