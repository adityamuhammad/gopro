package main

import (
	"gopro/config"
	"gopro/handlers"
	middleware "gopro/middeware"
	"gopro/models"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path"},
	)
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Histogram of HTTP request durations",
		},
		[]string{"method", "path"},
	)
)

func init() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found")
	}
	prometheus.MustRegister(requestCount)
	prometheus.MustRegister(requestDuration)
}

func main() {
	r := gin.Default()

	config.ConnectDatabase()

	config.DB.AutoMigrate(&models.User{}, &models.Story{}, &models.Message{})

	authPrometheus := gin.BasicAuth(gin.Accounts{
		os.Getenv("PROMETHEUS_BASIC_AUTH_USERNAME"): os.Getenv("PROMETHEUS_BASIC_AUTH_PASSWORD"),
	})

	r.GET("/metrics", authPrometheus, gin.WrapH(promhttp.Handler()))
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	authorized := r.Group("/")
	authorized.Use(middleware.AuthMiddleware())
	{
		authorized.GET("/realtime-message", handlers.RealTimeMessage)
		authorized.GET("/message/:id", handlers.GetMessage)
		authorized.POST("/stories", handlers.CreateStory)
		authorized.GET("/stories", handlers.GetStories)
		authorized.GET("/stories/:id", handlers.GetStory)
		authorized.PUT("stories/:id", handlers.UpdateStory)
		authorized.DELETE("stories/:id", handlers.DeleteStory)
	}

	r.Run(":8080")
}
