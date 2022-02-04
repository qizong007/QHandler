package qhandler

import (
	"log"
	"net/http"
	"strings"
	"time"
)

type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

// roots key be like: roots['GET'], roots['POST']
// handlers key eg, handlers['GET-/p/:lang/doc'], handlers['POST-/p/book']

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// only one '*' is allowed
func parsePattern(pattern string) []string {
	list := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range list {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

// 路由注册
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("[Register] %4s - %s\n", method, pattern)
	parts := parsePattern(pattern)
	key := method + "-" + pattern
	// add root
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.roots[method].sort()
	r.handlers[key] = handler
}

// 路由查找
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	n := root.search(searchParts, 0)
	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

// 处理函数（真正执行）
func (r *router) handle(c *Context) {
	start := time.Now()
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
	// 计时
	log.Println("[QHandler]", c.Path, "=> cost time:", time.Since(start))
}
