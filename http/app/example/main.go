package main

import (
	"fmt"
	"github.com/subchen/gstack/http/app"
	"github.com/subchen/gstack/http/app/middleware"
	"github.com/subchen/gstack/log"
)

/*
func userGroupHandles(g *app.Group) {
	g.Handle("GET", "/", handler)
	g.Handle("GET", "/{id}", handler)
	g.Handle("POST", "/{id}", handler)
	g.Handle("PUT", "/{id}", handler)
	g.Handle("PATCH", "/{id}", handler)
	g.Handle("DELETE", "/{id}", handler)
	g.Handle("GET", "/{id}/profiles", handler)
}
*/

func main() {

	a := app.New("/v2")
	a.RedirectTrailingSlash = true

	a.UsePre(middleware.Logger())

	a.Use(middleware.CORS())

	a.Use(WrapMiddleware(func(ctx *app.Context) {
		log.Fatal("Middleware")
	}))
	/*
		app.GET("/health", handler)
		app.POST("/stats", handler)

		app.Handle("*", "/stats", handler)
		app.Handle("POST,PUT", "/stats", handler)

		g := app.Group("/admin")
		g.Use(...)
		g.GET(...)

		app.Group("/users").Apply(userGroupHandles)
	*/

	a.GET("/health/", middleware.HealthCheckHandler)
	a.POST("/health/", middleware.HealthCheckHandler)

	a.Routes()

	fmt.Println("Listening http://127.0.0.1:8080/v2/")
	a.Run(":8080")
}

func WrapMiddleware(handler app.HandlerFunc) app.MiddlewareFunc {
	return func(next app.HandlerFunc) app.HandlerFunc {
		return func(ctx *app.Context) {
			handler(ctx)
			next(ctx)
		}
	}
}
