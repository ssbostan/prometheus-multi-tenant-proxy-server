package server

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/prometheus-community/prom-label-proxy/injectproxy"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/parser"
)

func injectProjectLabel(w http.ResponseWriter, r *http.Request) bool {
	labelEnforcer := injectproxy.NewEnforcer(
		false,
		&labels.Matcher{
			Type:  labels.MatchEqual,
			Name:  "project",
			Value: r.Header.Get("X-Project-Name"),
		},
	)
	if err := r.ParseForm(); err != nil {
		log.Printf("[ERROR] %s\n", err)
		return false
	}
	formData := r.Form
	for key, values := range formData {
		if key == "query" || key == "match[]" {
			formData.Del(key)
			for _, value := range values {
				expr, err := parser.ParseExpr(value)
				if err != nil {
					log.Printf("[ERROR] %s\n", err)
					return false
				}
				if err := labelEnforcer.EnforceNode(expr); err != nil {
					log.Printf("[ERROR] %s\n", err)
					return false
				}
				formData.Add(key, expr.String())
			}
		}
	}
	newFormData := formData.Encode()
	if r.Method == "GET" {
		r.URL.RawQuery = newFormData
	} else if r.Method == "POST" {
		r.Body = ioutil.NopCloser(strings.NewReader(newFormData))
		r.ContentLength = int64(len(newFormData))
	}
	return true
}
