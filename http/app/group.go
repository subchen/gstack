package app

import (
	"strings"
)

type (
	Group struct {
		prefix     string
		middleware []MiddlewareFunc
		router     *router
	}
)

func newGroup(prefix string, middleware []MiddlewareFunc, router *router) *Group {
	return &Group{
		prefix:     prefix,
		middleware: append([]MiddlewareFunc{}, middleware...), // copied
		router:     router,
	}
}

func (g *Group) Configure(fn func(*Group)) {
	fn(g)
}

func (g *Group) Use(middleware ...MiddlewareFunc) {
	g.middleware = append(g.middleware, middleware...)
}

func (g *Group) GET(path string, handler HandlerFunc) {
	g.route("GET", path, handler)
}

func (g *Group) POST(path string, handler HandlerFunc) {
	g.route("POST", path, handler)
}

func (g *Group) PUT(path string, handler HandlerFunc) {
	g.route("PUT", path, handler)
}

func (g *Group) PATCH(path string, handler HandlerFunc) {
	g.route("PATCH", path, handler)
}

func (g *Group) DELETE(path string, handler HandlerFunc) {
	g.route("DELETE", path, handler)
}

func (g *Group) OPTIONS(path string, handler HandlerFunc) {
	g.route("OPTIONS", path, handler)
}

func (g *Group) HEAD(path string, handler HandlerFunc) {
	g.route("HEAD", path, handler)
}

func (g *Group) CONNECT(path string, handler HandlerFunc) {
	g.route("CONNECT", path, handler)
}

func (g *Group) TRACE(path string, handler HandlerFunc) {
	g.route("TRACE", path, handler)
}

func (g *Group) Handle(methods string, path string, handler HandlerFunc) {
	if methods == "*" {
		methods = ALL_METHODS
	}
	for _, method := range strings.Split(methods, ",") {
		g.route(strings.TrimSpace(method), path, handler)
	}
}

func (g *Group) Group(path string) *Group {
	return newGroup(g.prefix+path, g.middleware, g.router)
}

func (g *Group) route(method string, path string, handler HandlerFunc) {
	// make handler chain
	for i := len(g.middleware) - 1; i >= 0; i-- {
		handler = g.middleware[i](handler)
	}
	g.router.add(method, g.prefix+path, handler)
}
