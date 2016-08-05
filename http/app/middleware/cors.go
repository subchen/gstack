package middleware

import (
	"github.com/subchen/gstack/http/app"
)

const (
	HeaderAccessControlRequestMethod    = "Access-Control-Request-Method"
	HeaderAccessControlRequestHeaders   = "Access-Control-Request-Headers"
	HeaderAccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	HeaderAccessControlAllowMethods     = "Access-Control-Allow-Methods"
	HeaderAccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	HeaderAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	HeaderAccessControlExposeHeaders    = "Access-Control-Expose-Headers"
	HeaderAccessControlMaxAge           = "Access-Control-Max-Age"
)

func CORS() app.MiddlewareFunc {
	return func(next app.HandlerFunc) app.HandlerFunc {
		return func(ctx *app.Context) {
			header := ctx.ResponseWriter.Header()

			header.Set(HeaderAccessControlAllowOrigin, "*")
			header.Set(HeaderAccessControlAllowMethods, app.DEFAULT_ALL_METHODS)
			header.Set(HeaderAccessControlAllowCredentials, "true")

			h := ctx.Header().Get(HeaderAccessControlRequestHeaders)
			header.Set(HeaderAccessControlAllowHeaders, h)

			//header.Set(HeaderAccessControlMaxAge, "0")

			next(ctx)
		}
	}
}
