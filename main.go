package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

var db = make(map[string]string)

func setupRouter() *gin.Engine {

	gin.DisableConsoleColor()
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Serve a link redirect
	r.GET("/:link", func(c *gin.Context) {
		link := c.Params.ByName("link")
		url, ok := db[link]
		if ok {
			c.Redirect(http.StatusFound, url)
		} else {
			c.String(http.StatusNotFound, "No URL associated with %s", link)
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

		db[shortUrl] = string(url)
		c.String(http.StatusOK, "Short URL is: %s", shortUrl)
	})

	return r
}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
