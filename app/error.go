package app

import (
	"fmt"
	"net/http"

    "github.com/subchen/gstack/errors"
)

func requestNotFound(req *http.Request) error {
	return errors.New("The requested URL " + req.URL.Path + " was not found on this server.")
}

type ErrorHandlerFunc func(w http.ResponseWriter, req *http.Request, code int, err interface{})

func defaultErrorHandlerFunc(w http.ResponseWriter, req *http.Request, code int, err interface{}) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	w.WriteHeader(code)
	fmt.Fprintf(w, "%+v\n", err)
}
