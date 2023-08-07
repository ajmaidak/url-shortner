package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var db = make(map[string]string)

func setupRouter() *gin.Engine {

	gin.DisableConsoleColor()
	r := gin.Default()

	var ctx = context.Background()

	redisAddr := os.Getenv("REDIS_ADDRESS")

	redisDB := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Serve a link redirect
	r.GET("/:link", func(c *gin.Context) {
		link := c.Params.ByName("link")
		url, err := redisDB.Get(ctx, link).Result()
		if err == nil {
			c.Redirect(http.StatusFound, url)
		} else if err == redis.Nil {
			c.String(http.StatusNotFound, "No URL associated with %s", link)
		} else {
			c.String(http.StatusInternalServerError, "database connection error")
		}
	})

	r.POST("/link", func(c *gin.Context) {
		url, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.String(http.StatusInternalServerError, "Interal Server Error")
		}
		// Calculate the SHA256 hash
		hash := sha256.Sum256(url)
		// Get the first 5 characters of the hash
		shortUrl := fmt.Sprintf("%x", hash)[:5]

		err = redisDB.Set(ctx, shortUrl, string(url), 0).Err()
		if err != nil {
			c.String(http.StatusInternalServerError, "Database Connection Error")
		}
		c.String(http.StatusOK, "Short URL is: %s", shortUrl)
	})

	return r
}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
