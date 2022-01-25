package main

import (
	"OSS-Score/server"
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {

	//uri := os.Getenv("MONGO_URI")
	uri := "mongodb+srv://local-user:oss-score@repos.76o1e.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"

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

	router := gin.Default()
	router.GET("/", server.Pong)
	router.Run("localhost:8080")

	// repoCatalog := "github"
	// repoOwner := "swagger-api"
	// repoName := "swagger-ui"
	// // repoOwner := "facebook"
	// // repoName := "react"
	// // repoOwner := "jasonlong"
	// // repoName := "isometric-contributions"
	// // repoOwner := "fonarevvichka"
	// // repoName := "OSS-Score"

	// score, ready := util.GetScore(mongoClient, repoCatalog, repoOwner, repoName, 12, 1)
	// if ready {
	// 	fmt.Println(score)
	// }
}
