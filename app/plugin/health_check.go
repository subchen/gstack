package plugin

import (
	"fmt"
	"net/http"
	"strings"
)

// Usage: router.HandleFunc("GET", "/health", HealthCheckHandler)
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	accept := r.Header.Get("Accept")

	w.WriteHeader(http.StatusOK)

	if strings.Contains(accept, "application/json") {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		fmt.Fprintf(w, "{\"health\": true}\n")
		return
	}

	if strings.Contains(accept, "application/xml") {
		w.Header().Set("Content-Type", "application/xml; charset=utf-8")
		fmt.Fprintf(w, "<health>true</health>\n")
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "OK\n")
}
