package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()

	r.Static("/static", "static")
	r.LoadHTMLGlob("templates/index.html")
	//r.LoadHTMLGlob("templates/*.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
		//c.HTML(http.StatusOK, "index.html", gin.H{
		//})
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
