package app

import (
	"fmt"
	"net/http"
	"net/url"
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

func (ctx *Context) URL() *url.URL {
	return ctx.Request.URL
}

func (ctx *Context) Path() string {
	return ctx.Request.URL.Path
}

func (ctx *Context) Method() string {
	return ctx.Request.Method
}

func (ctx *Context) GetHeader(name string) string {
	return ctx.Request.Header.Get(name)
}

func (ctx *Context) SetHeader(name string, value string) {
	ctx.ResponseWriter.Header().Set(name, value)
}

func (ctx *Context) Vars(name string) string {
	return ctx.vars[name]
}

func (ctx *Context) Form(name string) string {
	return ctx.Request.Form.Get(name)
}

func (ctx *Context) FormValues(name string) []string {
	return ctx.Request.Form[name]
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
