package middleware

import (
	"fmt"
	"github.com/subchen/gstack/http/app"
	"strings"
)

func HealthCheckHandler(ctx *app.Context) {
	accept := ctx.GetHeader("Accept")

	if strings.Contains(accept, "application/json") {
		ctx.SetHeader("Content-Type", "application/json; charset=utf-8")
		fmt.Fprintf(ctx.ResponseWriter, "{\"health\": true}\n")
		return
	}

	if strings.Contains(accept, "application/xml") {
		ctx.SetHeader("Content-Type", "application/xml; charset=utf-8")
		fmt.Fprintf(ctx.ResponseWriter, "<health>true</health>\n")
		return
	}

	ctx.SetHeader("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(ctx.ResponseWriter, "health OK\n")
}
