// +build go1.7

package app

import (
	"fmt"
	"net/http"
	"time"
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
	return r.WithContext(ctx), ctx
}

func requestContext(r *http.Request) *context {
	ctx, ok := r.Context().(*context)
	if !ok {
		panic("Missing context for request")
	}
	return ctx
}

func (*context) Deadline() (deadline time.Time, ok bool) {
	return
}

func (*context) Done() <-chan struct{} {
	return nil
}

func (*context) Err() error {
	return nil
}

func (c *context) Value(key interface{}) interface{} {
	return c.vars[key]
}

func (c *context) String() string {
	return fmt.Sprintf("%v", c.vars)
}
