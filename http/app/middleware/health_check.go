package middleware

import (
	"fmt"
	"github.com/subchen/gstack/http/app"
	"strings"
)

func HealthCheckHandler(ctx *app.Context) {
	header := ctx.Header()
	accept := header.Get("Accept")

	if strings.Contains(accept, "application/json") {
		header.Set("Content-Type", "application/json; charset=utf-8")
		fmt.Fprintf(ctx.ResponseWriter, "{\"health\": true}\n")
		return
	}

	if strings.Contains(accept, "application/xml") {
		header.Set("Content-Type", "application/xml; charset=utf-8")
		fmt.Fprintf(ctx.ResponseWriter, "<health>true</health>\n")
		return
	}

	header.Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(ctx.ResponseWriter, "OK\n")
}
