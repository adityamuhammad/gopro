package main

import (
	"fmt"
	"gopro/config"
	"gopro/handlers"
	middleware "gopro/middeware"
	"gopro/models"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/ulule/limiter/v3"
	ginLimiter "github.com/ulule/limiter/v3/drivers/middleware/gin"
	storeRedis "github.com/ulule/limiter/v3/drivers/store/redis"
)

var (
	totalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Number of GET requests.",
		},
		[]string{"path"},
	)

	responseStatus = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "response_status",
			Help: "Status of HTTP response",
		},
		[]string{"status"},
	)

	httpDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_time_seconds",
			Help:    "Duration of HTTP requests.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path"},
	)
)

func init() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found")
	}
}

func prometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(startTime).Seconds()
		statusCode := strconv.Itoa(c.Writer.Status())
		path := c.FullPath()

		// Record metrics
		httpDuration.WithLabelValues(path).Observe(duration)
		totalRequests.WithLabelValues(path).Inc()
		responseStatus.WithLabelValues(statusCode).Inc()
	}
}

func main() {
	r := gin.Default()
	r.Use(prometheusMiddleware())

	config.ConnectDatabase()

	// Create a rate limiter
	rate, errLimiter := limiter.NewRateFromFormatted("10-M") // 1000 requests/hour
	if errLimiter != nil {
		panic(errLimiter)
	}

	store, errStore := storeRedis.NewStoreWithOptions(config.RedisPool, limiter.StoreOptions{
		Prefix: "gopro_ratelimit",
	})

	if errStore != nil {
		panic(fmt.Sprintf("Failed to create Redis store: %v", errStore))
	}
	rateLimitMiddleware := ginLimiter.NewMiddleware(limiter.New(store, rate))

	errMinio := config.InitMinIO()
	if errMinio != nil {
		log.Fatalf("Failed to initialize MinIO client: %v", errMinio)
	}

	config.DB.AutoMigrate(&models.User{}, &models.Story{}, &models.Message{})

	authPrometheus := gin.BasicAuth(gin.Accounts{
		os.Getenv("PROMETHEUS_BASIC_AUTH_USERNAME"): os.Getenv("PROMETHEUS_BASIC_AUTH_PASSWORD"),
	})

	r.GET("/metrics", authPrometheus, gin.WrapH(promhttp.Handler()))
	r.POST("/register", handlers.Register)
	r.POST("/login", rateLimitMiddleware, handlers.Login)

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
		authorized.POST("/user/profile-image", handlers.UpdateProfileImage)
	}

	r.Run(":8080")
}
