package app

import (
	"errors"
	"fmt"
	"github.com/subchen/gstack/log"
	"net/http"
)

func requestNotFound(req *http.Request) error {
	return errors.New("The requested URL " + req.URL.Path + " was not found on this server.")
}

type ErrorHandlerFunc func(w http.ResponseWriter, req *http.Request, code int, err interface{})

func defaultErrorHandlerFunc(w http.ResponseWriter, req *http.Request, code int, err interface{}) {
	log.Fatal("test")

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	w.WriteHeader(code)
	fmt.Fprintf(w, "%v\n", err)
}
