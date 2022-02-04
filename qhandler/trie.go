package qhandler

import (
	"fmt"
	"sort"
	"strings"
)

// Q: 为什么不用 map 来存储路由信息？
// A: 因为 map 提供不了动态性，考虑用前缀树Trie（支持参数匹配':'和通配'*'）
type node struct {
	children []*node // 子节点，例如: [docs, tutorial]
	pattern  string  // 待匹配路由，例如 /p/:lang
	part     string  // 路由中的一部分，例如 :lang
	isWild   bool    // 是否精确匹配，part 含有 : 或 * 时为true
}

// 寻找第一个模糊匹配的节点，用于[插入]
func (n *node) findWildChild(part string) *node {
	for _, child := range n.children {
		if child.isWild {
			return child
		}
	}
	return nil
}

// 寻找第一个精确匹配的节点，用于[插入]
func (n *node) findSpecificChild(part string) *node {
	for _, child := range n.children {
		if child.part == part {
			return child
		}
	}
	return nil
}

// 所有匹配成功的节点，用于[查找]
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// 插入
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	part := parts[height]
	// 解决模糊匹配节点的冲突问题
	child := n.findWildChild(part)
	if child != nil && child.isWild && (part[0] == ':' || part[0] == '*') {
		panic(fmt.Sprintf("now %s(in %s) is conflict with %s", part, pattern, child.part))
	}
	// 再找一遍有没有其他节点，没有的话就插入
	child = n.findSpecificChild(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

// 查找
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}

// 优先级保证 eg. /18 > /:age
func (n *node) sort() {
	if n == nil {
		return
	}
	list := n.children
	sort.Slice(n.children, func(i, j int) bool {
		if !n.children[i].isWild &&  n.children[j].isWild {
			return true
		} else if n.children[i].isWild && !n.children[j].isWild {
			return false
		} else {
			return len(n.children[i].pattern) < len(n.children[j].pattern)
		}
	})
	if list != nil && len(list) > 0 {
		for i := range list {
			list[i].sort()
		}
	}
}