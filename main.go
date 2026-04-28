//go:build localserver

package main

import (
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	Init()
	r := gin.Default()
	SetRouter(r)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	_ = r.Run(":" + port)
}
