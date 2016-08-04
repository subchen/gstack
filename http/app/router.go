package app

import "strings"

type (
	router struct {
		root   *node
		routes []*route // a route list, unused currently
	}

	node struct {
		name         string // "{}", "*", "", ...
		children     map[string]*node
		methodRoutes methodRoutes
	}

	// routes is group of route, they have same path
	methodRoutes []*route

	route struct {
		path    string
		method  string
		handler HandlerFunc
	}
)

func newRouter() *router {
	return &router{
		root: &node{
			name:  "/",
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
				methodRoutes:   nil,
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

	n.methodRoutes = append(n.methodRoutes, route)
	r.routes = append(r.routes, route)
}

func (r *router) find(path string) methodRoutes {
	names := strings.Split(path, "/")
	if n := r.root.find(names[1], names[2:]); n != nil {
		return n.methodRoutes
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
			if child.methodRoutes != nil {
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
			if child.methodRoutes != nil {
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
func (r methodRoutes) find(method string) *route {
	for _, route := range r {
		if route.method == method {
			return route
		}
	}
	return nil
}

// methods returns all available methods (OPTIONS)
func (r methodRoutes) allows() string {
	var methods []string
	for _, route := range r {
		methods = append(methods, route.method)
	}
	return strings.Join(methods, ", ")
}
