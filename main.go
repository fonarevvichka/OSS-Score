package main

import (
	"OSS-Score/server"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/Pong", server.Pong)
	router.POST("/catalog/:catalog/owner/:owner/name/:name/scoreType/:scoreType", server.CalculateScore)
	router.GET("/catalog/:catalog/owner/:owner/name/:name/scoreType/:scoreType", server.GetCachedScore)
	router.Run()
}
