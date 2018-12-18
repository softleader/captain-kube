package main

import (
	"github.com/gin-gonic/gin"
	"net"
)

func main() {
	r := gin.Default()
	r.GET("/hosts/:host", func(c *gin.Context) {
		addrs, err := net.LookupHost(c.Param("host"))
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		c.JSON(200, addrs)
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
