package qhandler

import (
	"fmt"
	"reflect"
	"testing"
)

func newTestRouter() *router {
	r := newRouter()
	r.addRoute("GET", "/", nil)
	r.addRoute("GET", "/hello/:name", nil)
	r.addRoute("GET", "/hello/b", nil)
	r.addRoute("GET", "/hello/b/c", nil)
	r.addRoute("GET", "/assets/*filepath", nil)
	r.addRoute("GET", "/*assets/*filepath", nil)
	return r
}

func TestParsePattern(t *testing.T) {
	ok := reflect.DeepEqual(parsePattern("/p/:name"), []string{"p", ":name"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*"), []string{"p", "*"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*name/*"), []string{"p", "*name"})
	if !ok {
		t.Fatal("test parsePattern failed")
	}
}

func TestGetRoute(t *testing.T) {
	r := newTestRouter()
	n, ps := r.getRoute("GET", "/hello/wq")

	if n == nil {
		t.Fatal("<nil> shouldn't be returned")
	}

	if n.pattern != "/hello/:name" {
		t.Fatal("should match /hello/:name")
	}

	if ps["name"] != "wq" {
		t.Fatal("name should be equal to 'wq'")
	}

	fmt.Printf("matched path: %s, params['name']: %s\n", n.pattern, ps["name"])

}

func printNode(n *node, h int) {
	if n == nil {
		return
	}
	fmt.Printf("%p ", n)
	fmt.Println(h, n)
	list := n.children
	if list != nil && len(list) > 0 {
		for i := range list {
			printNode(list[i], h+1)
		}
	}
}

func TestTrie(t *testing.T) {
	r := newTestRouter()
	for _, root := range r.roots {
		fmt.Println("=================================")
		printNode(root, 0)
		fmt.Println("=================================")
	}
	n1, m1 := r.getRoute("GET", "/hello/b")
	fmt.Println(n1, m1)
	n2, m2 := r.getRoute("GET", "/hello/wq")
	fmt.Println(n2, m2)
	n3, m3 := r.getRoute("GET", "/hello/b/c")
	fmt.Println(n3, m3)
}
