package app

import (
	"net/http"
)

func withContext(router *Router) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		req, ctx := requestWithContext(req)
		if !router.match(req, ctx) {
			router.handleError(w, req, http.StatusNotFound, requestNotFound(req))
			return
		}
		router.ServeHTTP(w, req)
	})
}

func ListenAndServe(addr string, router *Router) error {
	server := &http.Server{
		Addr:    addr,
		Handler: withContext(router),
	}
	return server.ListenAndServe()
}

func ListenAndServeTLS(addr, certFile, keyFile string, router *Router) error {
	server := &http.Server{
		Addr:    addr,
		Handler: withContext(router),
	}
	return server.ListenAndServeTLS(certFile, keyFile)
}
