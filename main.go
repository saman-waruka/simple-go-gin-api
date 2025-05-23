package main

import (
	"fmt"
	"go-gin-api/handler"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// โหลด .env
	err := godotenv.Load()
	if err != nil {
		fmt.Println("⚠️  No .env file found, using default port")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8216"
	}

	r := gin.Default()


	r.GET("/ping", func(c *gin.Context) {
		currentTime := time.Now().Format("2006-01-02 15:04:05")
		c.JSON(http.StatusOK, gin.H{
			"message": "pong " + currentTime,
		})
	})

	r.Use(AuthMiddleware())
	r.Use(gin.Recovery())
	r.Use(gin.Logger())


	r.GET("/users", handler.GetUsers)
	r.POST("/users", handler.CreateUser)

	fmt.Printf("🚀 Server is running at: http://localhost:%s\n", port)
	r.Run(":" + port)
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("x-authorization")
		if authHeader != "secret" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}
