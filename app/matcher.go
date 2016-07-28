package app

import (
	"net/http"
	"regexp"
	"strings"
)

type matcher interface {
	match(req *http.Request, ctx *context) bool
}

//
var (
	always = &alwaysMatcher{}
)

type alwaysMatcher struct{}

func (m *alwaysMatcher) match(req *http.Request, ctx *context) bool {
	return true
}

// match request method
func newMethodsMatcher(methods string) matcher {
	if methods == "*" {
		return always
	}
	if strings.ContainsRune(methods, ',') {
		return &methodsMatcher{strings.Split(methods, ",")}
	} else {
		return &methodMatcher{methods}
	}
}

type methodMatcher struct {
	method string
}

func (m *methodMatcher) match(req *http.Request, ctx *context) bool {
	return req.Method == m.method
}

type methodsMatcher struct {
	methods []string
}

func (m *methodsMatcher) match(req *http.Request, ctx *context) bool {
	for _, method := range m.methods {
		if req.Method == method {
			return true
		}
	}
	return false
}

// match url path
func newPathMatcher(path string) matcher {
	if strings.ContainsRune(path, '{') || strings.Contains(path, "**") {
		return &pathRegexpMatcher{makePathRegexp(path)}
	} else {
		return &pathExactMatcher{path}
	}
}

func makePathRegexp(path string) *regexp.Regexp {
	pairs := strings.Split(path, "/")
	for i, name := range pairs {
		if strings.HasPrefix(name, "{") && strings.HasSuffix(name, "}") {
			pairs[i] = "(?P<" + name[1:len(name)-1] + ">[^/]+)"
		} else if name == "**" {
			pairs[i] = "(?P<suffix>.*)"
		} else {
			pairs[i] = regexp.QuoteMeta(name)
		}
	}

	re := strings.Join(pairs, "/")
	if strings.HasSuffix(re, "/") {
		re = re + "?"
	}
	re = "^" + re + "$"

	return regexp.MustCompile(re)
}

type pathExactMatcher struct {
	path string
}

func (m *pathExactMatcher) match(req *http.Request, ctx *context) bool {
	return ctx.path == m.path || ctx.path == m.path+"/"
}

type pathRegexpMatcher struct {
	pattern *regexp.Regexp
}

func (m *pathRegexpMatcher) match(req *http.Request, ctx *context) bool {
	matches := m.pattern.FindStringSubmatch(ctx.path)
	if len(matches) == 0 {
		return false
	}

	for i, name := range m.pattern.SubexpNames() {
		ctx.vars[name] = matches[i+1]
	}
	return true
}
