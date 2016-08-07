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
			ctx.SetHeader(HeaderAccessControlAllowOrigin, "*")
			ctx.SetHeader(HeaderAccessControlAllowMethods, app.DEFAULT_ALL_METHODS)
			ctx.SetHeader(HeaderAccessControlAllowCredentials, "true")

			ctx.SetHeader(HeaderAccessControlAllowHeaders, ctx.GetHeader(HeaderAccessControlRequestHeaders))

			//ctx.SetHeader(HeaderAccessControlMaxAge, "0")

			next(ctx)
		}
	}
}
