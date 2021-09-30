package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	routes "go_exploring/server"
)

func main() {
	router := gin.Default()

	router.GET("./pong", routes.Pong)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	
	router.NoRoute(func(c *gin.Context) {
              c.AbortWithStatus(http.StatusNotFound)
       })
	router.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
