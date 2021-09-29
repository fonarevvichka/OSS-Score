package server

import "github.com/gin-gonic/gin"

func Pong(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "ping",
	})
}
