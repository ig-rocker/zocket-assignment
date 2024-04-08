package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Message struct {
	Status string `json:"status"`
	Body   string `json:"body"`
}

func RouteHandler(ctx *gin.Context) {
	message := Message{
		Status: "Successful",
		Body:   "Hello from inside the route",
	}

	mes, err := json.Marshal(message)
	if err != nil {
		log.Fatal(err)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": string(mes),
	})
}




func main() {
	r := gin.Default()

	r.Use(RateLimit())

	r.GET("/route1", RouteHandler)
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Pong",
		})
	})

	r.Run()
}
