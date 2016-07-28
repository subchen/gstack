// +build !go1.7

package app

import (
	"net/http"
	"sync"
)

var (
	mutex sync.RWMutex
	data  = make(map[*http.Request]*context, 32)
)

type context struct {
	path string
	vars map[interface{}]interface{}
}

func requestWithContext(r *http.Request) (*http.Request, *context) {
	ctx := &context{
		path: r.URL.Path,
		vars: make(map[interface{}]interface{}, 4),
	}

	mutex.Lock()
	data[r] = ctx
	mutex.Unlock()

	return r, ctx
}

func requestContext(r *http.Request) *context {
	mutex.RLock()
	ctx := data[r]
	mutex.RUnlock()

	return ctx
}

func clearContext(r *http.Request) {
	mutex.Lock()
	delete(data, r)
	mutex.Unlock()
}
