package app

import (
	"net/http"
)

// Handler is an interface that objects can implement to be registered to serve as middleware
type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

// HandlerFunc is an adapter to allow the use of ordinary functions as handlers.
type HandlerFunc func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)

func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	f(w, r, next)
}

type middleware struct {
	handler Handler
	next    *middleware
}

// ServeHTTP implements http.HandlerFunc
func (m *middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.handler.ServeHTTP(w, r, m.next.ServeHTTP)
}

// wrap converts a http.Handler into a mux.Handler so it can be used as a middleware.
// The next http.HandlerFunc is automatically called after the Handler is executed.
func wrap(handler http.Handler) Handler {
	return HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		handler.ServeHTTP(w, r)
		next(w, r)
	})
}
