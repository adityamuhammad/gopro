package main

import (
	"gopro/config"
	"gopro/handlers"
	middleware "gopro/middeware"
	"gopro/models"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found")
	}
}

func main() {
	r := gin.Default()

	config.ConnectDatabase()

	config.DB.AutoMigrate(&models.User{}, &models.Story{}, &models.Message{})

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
