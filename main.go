package main

import (
	"QHandler/qhandler"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

type student struct {
	Name string
	Age  int8
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	r := qhandler.Default()
	r.Use(qhandler.Logger())
	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./static")

	stu1 := &student{Name: "Sam", Age: 20}
	stu2 := &student{Name: "Jack", Age: 22}
	r.GET("/", func(c *qhandler.Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})
	r.GET("/students", func(c *qhandler.Context) {
		c.HTML(http.StatusOK, "arr.tmpl", qhandler.H{
			"title":  "gee",
			"stuArr": [2]*student{stu1, stu2},
		})
	})

	r.GET("/date", func(c *qhandler.Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", qhandler.H{
			"title": "gee",
			"now":   time.Date(2022, 2, 6, 0, 0, 0, 0, time.UTC),
		})
	})

	r.GET("/panic", func(c *qhandler.Context) {
		names := []string{"wq"}
		c.String(http.StatusOK, names[100])
	})

	r.Run(":9999")
}
