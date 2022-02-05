package qhandler

import (
	"log"
	"time"
)

// middlewares 就是非业务的技术类组件
// Web 框架本身不可能去理解所有的业务，因而不可能实现所有的功能
// 因此，框架需要有一个插口，允许用户自己定义功能，嵌入到框架中，仿佛这个功能是框架原生支持的一样

// 关于插入点：
// 如果插入点太底层，中间件逻辑就会非常复杂
// 如果插入点离用户太近，那和用户直接定义一组函数，每次在 Handler 中手工调用没有多大的优势了
// 考虑放在 Group 上，如果作用于具体规则，还不如用户直接在 Handler 中调用直观

// 关于中间件的输入是什么？
// 中间件的输入，决定了扩展能力
// 暴露的参数如果太少，用户发挥空间有限

func Logger() HandlerFunc {
	return func(c *Context) {
		start := time.Now()
		c.Next()
		log.Println("[QHandler]", c.StatusCode, c.Path, "=> cost time:", time.Since(start))
	}
}
