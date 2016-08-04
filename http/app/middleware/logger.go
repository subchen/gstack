package middleware

import (
	"github.com/subchen/gstack/http/app"
	"fmt"
	"time"
)

func Logger() app.MiddlewareFunc {
	return func(next app.HandlerFunc) app.HandlerFunc {
		return func(ctx *app.Context) {
			fmt.Printf("%s %s\n", ctx.Method(), ctx.Request.URL.String())

			start := time.Now()
			next(ctx)

			fmt.Printf("time: %d\n", time.Since(start))
		}
	}
}
