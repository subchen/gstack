package main

import (
	"fmt"
	"github.com/subchen/gstack/app"
	"net/http"
	"os"
)

func mainRouter() *app.Router {
	router := app.NewRouter("/v2")

	router.UseFunc(StatisticsMiddleware)

	//router.Subrouter("/admin", adminRouter())
	//router.Subrouter("/users", userRouter())

	router.Subrouters(
		adminRouter("/admin"),
		userRouter("/users"),
	)

	router.HandleFunc("GET", "/health", HealthCheckHandler)
	router.HandleFunc("GET", "/stats", StatisticsHandler)

	return router
}

func adminRouter(path string) *app.Router {
	router := app.NewRouter(path)

	//router.HandleFunc("GET", "/", handler)
	router.HandleFunc("GET", "/env", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%v\n", os.Environ())
	})
	router.HandleFunc("GET", "/stats", handler)

	return router
}

func userRouter(path string) *app.Router {
	router := app.NewRouter(path)

	router.HandleFunc("GET", "/", handler)
	router.HandleFunc("GET", "/{id}", handler)
	router.HandleFunc("POST", "/{id}", handler)
	router.HandleFunc("PUT", "/{id}", handler)
	router.HandleFunc("PATCH", "/{id}", handler)
	router.HandleFunc("DELETE", "/{id}", handler)

	router.HandleFunc("GET", "/{id}/profiles", handler)

	return router
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}

func main() {
	fmt.Println("Listening http://127.0.0.1:8080/v2/")
	app.ListenAndServe(":8080", mainRouter())
}
