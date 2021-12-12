package main

import (
	"OSS-Score/server"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.POST("/catalog/:catalog/owner/:owner/name/:name", server.QueryProject)
	router.GET("/catalog/:catalog/owner/:owner/name/:name/scoreType/:scoreType", server.GetCachedScore)
	router.Run()
}
