package server

import (
	"OSS-Score/util"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func Pong(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "ping",
	})
}

func CalculateScore(c *gin.Context) {
	repoCatalog := c.Param("catalog")
	repoOwner := c.Param("owner")
	repoName := c.Param("name")
	// scoreType := c.Param("scoreType")

	go util.CalculateScore(repoCatalog, repoOwner, repoName, 24, 0)

	c.JSON(http.StatusOK, gin.H{
		"message": "Score request accepted",
	})
}

func GetScore(c *gin.Context) {
	uri := os.Getenv("MONGO_URI")
	// Create a new mongo_client and connect to the server
	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		log.Fatalln(err)
		panic(err)
	}
	defer func() {
		if err = mongoClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// Ping the primary
	if err := mongoClient.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected and pinged.")

	catalog := c.Param("catalog")
	owner := c.Param("owner")
	name := c.Param("name")
	scoreType := c.Param("scoreType")
	score, scoreStatus := util.GetScore(mongoClient, catalog, owner, name, scoreType, 12) // TEMP HARDCODED TO 12 MONTHS

	var message string
	if scoreStatus == 0 {
		message = "Score not yet calculated"
	} else if scoreStatus == 1 {
		message = "Score calculation in progress"
	} else {
		message = "Score ready"
	}
	// retrieve score from database
	//if score not in database send wait / error message
	//if score in database send score
	c.JSON(http.StatusOK, gin.H{
		"message": message,
		"score":   score,
	})
}
