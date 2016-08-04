package app

import "strings"

type (
	router struct {
		root       *node
		routesList []*routes // all path routes
	}

	node struct {
		name     string // may be "{}", "*", "", ...
		children map[string]*node
		routes   *routes
	}

	// routes is group of route, they have same path
	routes struct {
		path   string
		routes []*route
	}

	// route is a defined handler
	route struct {
		path    string // origin path
		method  string
		handler HandlerFunc
	}
)

func newRouter() *router {
	return &router{
		root: &node{
			name: "/",
		},
	}
}

// add registers a handler into router tree
func (r *router) add(method string, path string, handler HandlerFunc) {
	names := strings.Split(path, "/")

	n := r.root
	for i := 1; i < len(names); i++ {
		name := names[i]
		if strings.Contains(name, "*") {
			name = "*" // any
		} else if strings.Contains(name, "{") {
			name = "{}" // param
		}

		nn, ok := n.children[name]
		if !ok {
			nn = &node{
				name:     name,
				children: nil,
				routes:   nil,
			}

			if n.children == nil {
				n.children = make(map[string]*node, 4)
			}
			n.children[name] = nn
		}

		n = nn
	}

	route := &route{
		path:    path,
		method:  method,
		handler: handler,
	}

	if n.routes == nil {
		n.routes = &routes{
			path:   path,
			routes: nil,
		}
		r.routesList = append(r.routesList, n.routes)
	}

	n.routes.routes = append(n.routes.routes, route)
}

func (r *router) find(path string) *routes {
	names := strings.Split(path, "/")
	if n := r.root.find(names[1], names[2:]); n != nil {
		return n.routes
	}
	return nil
}

func (n *node) find(name string, path []string) *node {
	if len(n.children) == 0 {
		return nil
	}

	// static
	if child, ok := n.children[name]; ok {
		if len(path) == 0 {
			if child.routes != nil {
				return child // match
			}
		} else {
			nn := child.find(path[0], path[1:])
			if nn != nil {
				return nn
			}
		}
	}

	// param
	if child, ok := n.children["{}"]; ok {
		if len(path) == 0 {
			if child.routes != nil {
				return child // match
			}
		} else {
			nn := child.find(path[0], path[1:])
			if nn != nil {
				return nn
			}
		}
	}

	// any *
	if child, ok := n.children["*"]; ok {
		return child
	}

	return nil
}

// find returns the matched route
func (r *routes) find(method string) *route {
	for _, route := range r.routes {
		if route.method == method {
			return route
		}
	}
	return nil
}

// methods returns all available methods (OPTIONS)
func (r *routes) allows() string {
	var methods []string
	for _, route := range r.routes {
		methods = append(methods, route.method)
	}
	return strings.Join(methods, ", ")
}
