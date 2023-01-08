package server

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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
	rawQueryForm, postForm := splitForm(r.Form, r.PostForm)
	newRawQueryForm, newRawQueryFormOK := enforceProjectLabel(rawQueryForm, labelEnforcer)
	newPostForm, newPostFormOK := enforceProjectLabel(postForm, labelEnforcer)
	if newRawQueryFormOK == false || newPostFormOK == false {
		return false
	}
	r.URL.RawQuery = newRawQueryForm.Encode()
	r.Body = ioutil.NopCloser(strings.NewReader(newPostForm.Encode()))
	r.ContentLength = int64(len(newPostForm.Encode()))
	return true
}

func enforceProjectLabel(form url.Values, labelEnforcer *injectproxy.Enforcer) (url.Values, bool) {
	for key, values := range form {
		if key == "query" || key == "match[]" {
			form.Del(key)
			for _, value := range values {
				expr, err := parser.ParseExpr(value)
				if err != nil {
					log.Printf("[ERROR] %s\n", err)
					return url.Values{}, false
				}
				if err := labelEnforcer.EnforceNode(expr); err != nil {
					log.Printf("[ERROR] %s\n", err)
					return url.Values{}, false
				}
				form.Add(key, expr.String())
			}
		}
	}
	return form, true
}
