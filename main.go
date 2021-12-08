package main

import (
	"OSS-Score/util"
	"fmt"
	"time"
)

func main() {
	// router := gin.Default()
	// router.GET("/Pong", server.Pong)
	// router.POST("/catalog/:catalog/owner/:owner/name/:name", server.QueryProject)
	// router.GET("/catalog/:catalog/owner/:owner/name/:name/scoreType/:scoreType", server.GetCachedScore)
	// router.Run()

	repoInfo := util.QueryGithub("github", "facebook", "react", time.Now().AddDate(-1, 0, 0))
	fmt.Println(repoInfo)
}

// func async_test(c1 chan int, c2 chan int) {
// 	time.Sleep(time.Second * 1)

// 	c1 <- 2
// 	c2 <- 4
// }

// func main() {
// 	c1 := make(chan int)
// 	c2 := make(chan int)
// 	go async_test(c1, c2)

// 	// var r1 int
// 	// var r2 int
// 	// for i := 0; i < 2; i++ {
// 	// 	// Await both of these values
// 	// 	// simultaneously, printing each one as it arrives.
// 	// 	select {
// 	// 	case r1 = <-c1:
// 	// 		r1 = r1
// 	// 	case r2 = <-c2:
// 	// 		r2 = r2
// 	// 	}
// 	// }

// 	r1 := <-c1
// 	r2 := <-c2
// 	fmt.Println(r1 + r2)
// }
