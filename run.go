package main

import (
	"github.com/gin-gonic/gin"
	server "go_exploring/server"
)

func main() {
	router := gin.Default()

	router.GET("./pong", server.Pong)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
