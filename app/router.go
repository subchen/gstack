package app

import (
	"net/http"
	"strings"
)

// NewRouter create a Router with path prefix
func NewRouter(path string) *Router {
	if path == "" || path == "/" {
		path = ""
	} else {
		if !strings.HasPrefix(path, "/") {
			panic("path should be prefix with '/'")
		}
		if strings.HasSuffix(path, "/") {
			panic("path should NOT be suffix with '/'")
		}
	}

	return &Router{
		pathWithSlash:    path + "/",
		pathWithoutSlash: path,
	}
}

// Router is a group of middlewares and routes, implements http.Handler
type Router struct {
	parent           *Router
	pathWithSlash    string
	pathWithoutSlash string
	middleware       *middleware
	middlewares      []middleware
	routes           []route
	subrouters       []Router
	errorHandlerFunc ErrorHandlerFunc
}

// Use registers a Handler as middleware
func (r *Router) Use(handler Handler) {
	r.middlewares = append(r.middlewares, middleware{handler, nil})
}

// UseFunc registers a HandlerFunc as middleware
func (r *Router) UseFunc(handlerFunc func(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc)) {
	r.Use(HandlerFunc(handlerFunc))
}

// UseHandler registers a http.Handler as middleware
func (r *Router) UseHandler(handler http.Handler) {
	r.Use(wrap(handler))
}

// UseHandlerFunc registers http.HandlerFunc as middleware
func (r *Router) UseHandlerFunc(handlerFunc func(rw http.ResponseWriter, req *http.Request)) {
	r.UseHandler(http.HandlerFunc(handlerFunc))
}

// Handle registers the handler as route
func (r *Router) Handle(methods string, path string, handler http.Handler) {
	r.routes = append(r.routes, route{
		handler: handler,
		matchers: []matcher{
			newMethodsMatcher(methods),
			newPathMatcher(path),
		},
	})
}

// HandleFunc registers the handlerFunc as route
func (r *Router) HandleFunc(methods string, path string, handlerFunc func(w http.ResponseWriter, req *http.Request)) {
	r.Handle(methods, path, http.HandlerFunc(handlerFunc))
}

// GET registers the handlerFunc as route
func (r *Router) GET(path string, handlerFunc func(w http.ResponseWriter, req *http.Request)) {
	r.Handle("GET", path, http.HandlerFunc(handlerFunc))
}

// POST registers the handlerFunc as route
func (r *Router) POST(path string, handlerFunc func(w http.ResponseWriter, req *http.Request)) {
	r.Handle("POST", path, http.HandlerFunc(handlerFunc))
}

// PUT registers the handlerFunc as route
func (r *Router) PUT(path string, handlerFunc func(w http.ResponseWriter, req *http.Request)) {
	r.Handle("PUT", path, http.HandlerFunc(handlerFunc))
}

// PATCH registers the handlerFunc as route
func (r *Router) PATCH(path string, handlerFunc func(w http.ResponseWriter, req *http.Request)) {
	r.Handle("PATCH", path, http.HandlerFunc(handlerFunc))
}

// DELETE registers the handlerFunc as route
func (r *Router) DELETE(path string, handlerFunc func(w http.ResponseWriter, req *http.Request)) {
	r.Handle("DELETE", path, http.HandlerFunc(handlerFunc))
}

// Subrouters registers sub-routers
func (r *Router) Subrouters(routers ...*Router) {
	for _, router := range routers {
		router.parent = r
		r.subrouters = append(r.subrouters, *router)
	}
}

func (r *Router) HandleError(handlerFunc func(w http.ResponseWriter, req *http.Request, code int, err interface{})) {
	r.errorHandlerFunc = ErrorHandlerFunc(handlerFunc)
}

// ServeHTTP implements http.Handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// handle 500 panic
	defer func() {
		if err := recover(); err != nil {
			r.handleError(w, req, http.StatusInternalServerError, err)
		}
	}()

	// make middleware chain on first request
	if r.middleware == nil {
		// append r.route as last middleware
		r.middlewares = append(r.middlewares, middleware{HandlerFunc(r.route), nil})

		r.middleware = &(r.middlewares[0])
		for i := 1; i < len(r.middlewares); i++ {
			r.middlewares[i-1].next = &(r.middlewares[i])
		}
	}

	r.middleware.ServeHTTP(w, req)
}

// match router path prefix
func (r *Router) match(req *http.Request, ctx *context) bool {
	if r.pathWithSlash == "/" {
		return true
	}

	if ctx.path == r.pathWithoutSlash {
		ctx.path = "/"
		return true
	}
	if strings.HasPrefix(ctx.path, r.pathWithSlash) {
		ctx.path = strings.TrimPrefix(ctx.path, r.pathWithoutSlash)
		return true
	}

	return false
}

func (r *Router) route(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	// WARNING: Don't use param "next", here "next" is always nil.

	ctx := requestContext(req)

	for _, router := range r.subrouters {
		if router.match(req, ctx) {
			router.ServeHTTP(w, req)
			return
		}
	}

	for _, route := range r.routes {
		if route.match(req, ctx) {
			route.handler.ServeHTTP(w, req)
			return
		}
	}

	r.handleError(w, req, http.StatusNotFound, requestNotFound(req))
}

func (r *Router) handleError(w http.ResponseWriter, req *http.Request, code int, err interface{}) {
	if r.errorHandlerFunc != nil {
		r.errorHandlerFunc(w, req, code, err)
	} else if r.parent != nil {
		r.parent.handleError(w, req, code, err)
	} else {
		defaultErrorHandlerFunc(w, req, code, err)
	}
}
