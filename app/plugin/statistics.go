package plugin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/subchen/gstack/iif"
)

type statistics struct {
	url       string        `json:"req_url"`
	method    string        `json:"req_method"`
	hitsAll   uint64        `json:"times"`
	hitsErr   uint64        `json:"times_failed"`
	timeTotal time.Duration `json:"time_total"`
	timeMin   time.Duration `json:"time_min"`
	timeMax   time.Duration `json:"time_max"`
}

// URLMap contains several statistics struct to log different data
type URLMap struct {
	mutex sync.RWMutex
	urls  map[string]map[string]*statistics // url->method->stats
}

var (
	urlmap = &URLMap{
		urls: make(map[string]map[string]*statistics, 16),
	}
)

type ResponseWriterWrapper struct {
	w          http.ResponseWriter
	statusCode int
}

func (w *ResponseWriterWrapper) Write(bytes []byte) (int, error) {
	return w.w.Write(bytes)
}

func (w *ResponseWriterWrapper) Header() http.Header {
	return w.w.Header()
}

func (w *ResponseWriterWrapper) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.w.WriteHeader(statusCode)
}

// Usage: router.UseFunc(StatisticsMiddleware)
func StatisticsMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	url := r.URL.Path
	method := r.Method
	start := time.Now()

	rw := &ResponseWriterWrapper{w, 0}
	next(rw, r) // to exec

	elapsed := time.Since(start)
	hitErr := iif.Uint64(rw.statusCode >= 400, 1, 0)

	// save stats
	urlmap.mutex.Lock()
	defer urlmap.mutex.Unlock()

	methodMap, ok := urlmap.urls[url]
	if !ok {
		methodMap = make(map[string]*statistics)
		urlmap.urls[url] = methodMap
	}

	stats, ok := methodMap[method]
	if !ok {
		stats = &statistics{
			url:       url,
			method:    method,
			hitsAll:   1,
			hitsErr:   uint64(hitErr),
			timeMin:   elapsed,
			timeMax:   elapsed,
			timeTotal: elapsed,
		}
		methodMap[method] = stats
	} else {
		if stats.timeMin > elapsed {
			stats.timeMin = elapsed
		} else if stats.timeMax < elapsed {
			stats.timeMax = elapsed
		}
		stats.timeTotal = stats.timeTotal + elapsed

		stats.hitsAll = stats.hitsAll + 1
		stats.hitsErr = stats.hitsErr + hitErr
	}
}

// Usage: router.HandleFunc("GET", "/stats", StatisticsHandler)
func StatisticsHandler(w http.ResponseWriter, r *http.Request) {
	resultLists := make([]map[string]interface{}, 0, len(urlmap.urls))

	for k, v := range urlmap.urls {
		fmt.Printf("%v.len = %v\n", k, len(v))
		for kk, vv := range v {
			fmt.Printf("%v = %v\n", kk, vv)
			result := map[string]interface{}{
				"req_url":    k,
				"req_method": kk,
				"req_total":  vv.hitsAll,
				"req_failed": vv.hitsErr,
				"time_min":   vv.timeMin,
				"time_max":   vv.timeMax,
				"time_avg":   uint64(vv.timeTotal) / vv.hitsAll,
			}
			resultLists = append(resultLists, result)
		}
	}

	if bytes, err := json.MarshalIndent(resultLists, "", "  "); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(bytes)
	}
}
