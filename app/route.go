package app

import (
	"net/http"
)

// route is a handler with matchers
type route struct {
	matchers []matcher
	handler  http.Handler
}

func (r *route) match(req *http.Request, ctx *context) bool {
	for _, m := range r.matchers {
		if !m.match(req, ctx) {
			return false
		}
	}
	return true
}
