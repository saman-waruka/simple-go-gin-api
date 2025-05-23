package main

import (
	"fmt"
	"go-gin-api/handler"
	"go-gin-api/middlewares"
	"net/http"
	"os"
	"time"
  "go.uber.org/zap"
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

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	r := gin.Default()
	r.Use(middleware.RecoveryWithZap(logger))
	r.Use(gin.Logger())

	r.GET("/ping", func(c *gin.Context) {
		currentTime := time.Now().Format("2006-01-02 15:04:05")
		c.JSON(http.StatusOK, gin.H{
			"message": "pong " + currentTime,
		})
	})

	r.GET("/panic", func(c *gin.Context) {
		panic("😱 Something exploded!")
	})

	authorized := r.Group("/", AuthMiddleware())
	{
		authorized.GET("/users", handler.GetUsers)
		authorized.POST("/users", handler.CreateUser)
	}

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
