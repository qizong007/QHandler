package main

import (
	"QHandler/qhandler"
	"net/http"
)

func main() {
	r := qhandler.New()
	r.GET("/", func(c *qhandler.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})
	r.POST("/login", func(c *qhandler.Context) {
		c.JSON(http.StatusOK, qhandler.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})
	r.GET("/hello/:name", func(c *qhandler.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})

	r.GET("/assets/*filepath", func(c *qhandler.Context) {
		c.JSON(http.StatusOK, qhandler.H{"filepath": c.Param("filepath")})
	})

	r.Run(":9999")

}
