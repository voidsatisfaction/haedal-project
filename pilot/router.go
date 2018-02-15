package haedal

import (
	"net/http"
	"strings"
)

type HandlerFunc func(*Context)

type router struct {
	// method / pattern / handler
	// e.g POST / "/users" / createUser
	handlers map[string]map[string]HandlerFunc
}

func NewRouter() *router {
	return &router{
		make(map[string]map[string]HandlerFunc),
	}
}

func (r *router) HandleFunc(method, pattern string, h HandlerFunc) {
	m, ok := r.handlers[method]
	if !ok {
		m = make(map[string]HandlerFunc)
		r.handlers[method] = m
	}

	m[pattern] = h
}

func (r *router) handler() HandlerFunc {
	return func(c *Context) {
		for pattern, handler := range r.handlers[c.Request.Method] {
			if params, ok := match(pattern, c.Request.URL.Path); ok {
				for k, v := range params {
					c.Params[k] = v
				}
				handler(c)
				return
			}
		}
		http.NotFound(c.ResponseWriter, c.Request)
		return
	}
}

func match(pattern, path string) (map[string]string, bool) {
	if pattern == path {
		return nil, true
	}

	patterns := strings.Split(pattern, "/")
	paths := strings.Split(path, "/")

	if len(patterns) != len(paths) {
		return nil, false
	}

	params := make(map[string]string)

	for i := 0; i < len(patterns); i++ {
		switch {
		case patterns[i] == paths[i]:

		case len(patterns[i]) > 0 && patterns[i][0] == ':':
			params[patterns[i][1:]] = paths[i]
		default:
			return nil, false
		}
	}

	return params, true
}
