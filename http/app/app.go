package app

import (
	"github.com/subchen/gstack/errors"
	"net/http"
	"strings"
)

type (
	HandlerFunc func(ctx *Context)

	MiddlewareFunc func(next HandlerFunc) HandlerFunc

	App struct {
		prefix        string
		middlewarePre []MiddlewareFunc
		middleware    []MiddlewareFunc
		router        *router
		handler       HandlerFunc

		// Enables automatic redirection if the current route can't be matched but a
		// handler for the path with (without) the trailing slash exists.
		// For example if /foo/ is requested but a route only exists for /foo, the
		// client is redirected to /foo with http status code 301 for GET requests
		// and 307 for all other request methods.
		RedirectTrailingSlash bool
	}
)

const (
	ALL_METHODS = "GET,POST,PUT,PATCH,DELETE,OPTIONS,HEAD,CONNECT,TRACE"
)

func New(prefix string) *App {
	if prefix == "" || prefix == "/" {
		prefix = ""
	} else {
		if !strings.HasPrefix(prefix, "/") {
			panic("prefix should be start with '/'")
		}
		if strings.HasSuffix(prefix, "/") {
			panic("prefix should NOT be end with '/'")
		}
	}

	app := &App{}
	app.prefix = prefix
	app.router = newRouter()
	app.RedirectTrailingSlash = true
	return app
}

func (app *App) UsePre(middleware ...MiddlewareFunc) {
	app.middlewarePre = append(app.middlewarePre, middleware...)
}

func (app *App) Use(middleware ...MiddlewareFunc) {
	app.middleware = append(app.middleware, middleware...)
}

func (app *App) GET(path string, handler HandlerFunc) {
	app.add("GET", path, handler)
}

func (app *App) POST(path string, handler HandlerFunc) {
	app.add("POST", path, handler)
}

func (app *App) PUT(path string, handler HandlerFunc) {
	app.add("PUT", path, handler)
}

func (app *App) PATCH(path string, handler HandlerFunc) {
	app.add("PATCH", path, handler)
}

func (app *App) DELETE(path string, handler HandlerFunc) {
	app.add("DELETE", path, handler)
}

func (app *App) OPTIONS(path string, handler HandlerFunc) {
	app.add("OPTIONS", path, handler)
}

func (app *App) HEAD(path string, handler HandlerFunc) {
	app.add("HEAD", path, handler)
}

func (app *App) CONNECT(path string, handler HandlerFunc) {
	app.add("CONNECT", path, handler)
}

func (app *App) TRACE(path string, handler HandlerFunc) {
	app.add("TRACE", path, handler)
}

func (app *App) Handle(methods string, path string, handler HandlerFunc) {
	if methods == "*" {
		methods = ALL_METHODS
	}
	for _, method := range strings.Split(methods, ",") {
		app.add(strings.TrimSpace(method), path, handler)
	}
}

func (app *App) Group(path string) *Group {
	return newGroup(app.prefix+path, app.middleware, app.router)
}

func (app *App) add(method string, path string, handler HandlerFunc) {
	// make handler chain
	for i := len(app.middleware) - 1; i >= 0; i-- {
		handler = app.middleware[i](handler)
	}
	app.router.add(method, app.prefix+path, handler)
}

// Routes returns all register route path
func (app *App) Routes() []string {
	var paths []string
	for _, routes := range app.router.routesList {
		paths = append(paths, routes.path)
	}
	return paths
}

// ServeHTTP implements http.Handler
func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := newContext(w, r, app)

	// handle 500 panic
	defer func() {
		if err := recover(); err != nil {
			ctx.Error(http.StatusInternalServerError, errors.Newf("panic: %v", err))
		}
	}()

	if app.handler == nil {
		// it make middleware and handler as a single chain handler
		// and cached in app.handler
		// it will be built on first request
		handler := func(ctx *Context) {
			path := ctx.Path()
			routes := app.router.find(path)

			if routes == nil {
				// try to fix url and redirect
				if app.RedirectTrailingSlash {
					last := len(path) - 1
					if path[last] == '/' {
						path = path[0:last]
					} else {
						path = path + "/"
					}
					routes = app.router.find(path)
					if routes != nil {
						ctx.Redirect(path) // trim slash redirect
						return
					}
				}

				ctx.Error(http.StatusNotFound, errors.Newf("Request not found: %s", ctx.Path()))
				return
			}

			route := routes.find(ctx.Method())
			if route == nil {
				ctx.ResponseWriter.Header().Set("Allow", routes.allows())
				if ctx.Method() == "OPTIONS" {
					ctx.ResponseWriter.WriteHeader(http.StatusNoContent)
				} else {
					ctx.ResponseWriter.WriteHeader(http.StatusMethodNotAllowed)
				}
				return
			}

			// execute middleware and handler
			route.handler(ctx)
		}

		// chain pre middleware
		for i := len(app.middlewarePre) - 1; i >= 0; i-- {
			handler = app.middlewarePre[i](handler)
		}

		app.handler = handler
	}

	// execute middleware and handler
	app.handler(ctx)
}

func (app *App) Run(addr string) error {
	return http.ListenAndServe(addr, app)
}

func (app *App) RunTLS(addr, certFile, keyFile string) error {
	return http.ListenAndServeTLS(addr, certFile, keyFile, app)
}
