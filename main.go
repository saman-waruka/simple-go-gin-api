package main

import (
	"fmt"
	"go-gin-api/handler"
	"go-gin-api/middleware"
	"net/http"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
  sentrygin "github.com/getsentry/sentry-go/gin"
  "go.uber.org/zap"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/files"
	_ "go-gin-api/docs"
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

	sentryDNS := os.Getenv("SENTRY_DNS")
	if sentryDNS == "" || sentryDNS == "null" || sentryDNS == "undefined"  {
		panic("⚠️  SENTRY_DSN is not set")
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()


	// To initialize Sentry's handler, you need to initialize Sentry itself beforehand
if err := sentry.Init(sentry.ClientOptions{
  Dsn: sentryDNS,
	EnableTracing: true,
  // Set TracesSampleRate to 1.0 to capture 100%
  // of transactions for tracing.
  // We recommend adjusting this value in production,
  TracesSampleRate: 1.0,
}); err != nil {
  fmt.Printf("Sentry initialization failed: %v\n", err)
}

	r := gin.Default()
	// Once it's done, you can attach the handler as one of your middleware
	r.Use(sentrygin.New(sentrygin.Options{}))
	r.Use(middleware.PrometheusMiddleware())
	r.Use(middleware.RecoveryWithZap(logger))
	r.Use(gin.Logger())

	r.GET("/ping", func(c *gin.Context) {
		currentTime := time.Now().Format("2006-01-02 15:04:05")
		c.JSON(http.StatusOK, gin.H{
			"message": "pong " + currentTime,
		})
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// @title           Simple GO GIN API Example for thrown error
	// @version         1.0
	// @description     This is a simple example of a GO GIN API with integrated monitoring tools (Sentry, Prometheus, Grafana, Zap logging).
	// @host            localhost:8216
	// @BasePath        /
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
