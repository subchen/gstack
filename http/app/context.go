package app

import (
	"fmt"
	"net/http"
)

type Context struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	vars           map[string]string
	app            *App
}

func newContext(w http.ResponseWriter, r *http.Request, app *App) *Context {
	return &Context{
		ResponseWriter: w,
		Request:        r,
		vars:           nil,
		app:            app,
	}
}

func (ctx *Context) Path() string {
	return ctx.Request.URL.Path
}

func (ctx *Context) Method() string {
	return ctx.Request.Method
}

func (ctx *Context) Header() http.Header {
	return ctx.Request.Header
}

func (ctx *Context) Vars(name string) string {
	return ctx.vars[name]
}

func (ctx *Context) Redirect(url string) {
	code := http.StatusMovedPermanently // 301
	if ctx.Request.Method != "GET" {
		code = http.StatusTemporaryRedirect // 307
	}
	http.Redirect(ctx.ResponseWriter, ctx.Request, url, code)
}

func (ctx *Context) Error(code int, err error) {
	w := ctx.ResponseWriter

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	fmt.Fprintf(w, "%+v\n", err)
}
