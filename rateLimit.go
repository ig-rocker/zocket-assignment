package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"

	"github.com/gin-gonic/gin"
)

type Client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	mu      sync.Mutex
	clients = make(map[string]*Client)
)

func ClearRateLimiter() {
	for {
		time.Sleep(time.Second * 5)
		mu.Lock()
		for ip, client := range clients {
			if time.Since(client.lastSeen) > 5*time.Second {
				delete(clients, ip)
			}
		}
		mu.Unlock()
	}
}

func RateLimit() gin.HandlerFunc {
	go ClearRateLimiter()
	return func(c *gin.Context) {
		mu.Lock()
		ip := c.Request.Host
		if _, found := clients[ip]; !found {
			clients[ip] = &Client{limiter: rate.NewLimiter(1, 2)}
		}

		clients[ip].lastSeen = time.Now()

		if !clients[ip].limiter.Allow() {
			mu.Unlock()
			message := Message{
				Status: "Request failed",
				Body:   " API rate limit exceeded, Please try after some time",
			}
			mes, err := json.Marshal(message)
			if err != nil {
				log.Fatal(err)
			}
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    http.StatusTooManyRequests,
				"message": string(mes),
			})
			return
		}
		mu.Unlock()
		c.Next()
	}

}
