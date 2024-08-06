package handlers

import (
	"gopro/config"
	"gopro/models"
	"gopro/models/converter"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateStory(c *gin.Context) {
	var createStoryRequest models.CreateStoryRequest
	if err := c.ShouldBindJSON(&createStoryRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var story models.Story
	story.UserID = uint(userID.(float64))
	story.Status = createStoryRequest.Status

	config.DB.Create(&story)
	var createStoryResponse = converter.MapStoryToCreateStoryResponse(&story)
	c.JSON(http.StatusOK, gin.H{"data": createStoryResponse})
}

func GetStories(c *gin.Context) {
	var stories []models.Story
	if err := config.DB.Find(&stories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch storys"})
		return
	}
	var getStoriesResponse = converter.MapStoriesToGetStoriesResponse(&stories)
	c.JSON(http.StatusOK, gin.H{"data": getStoriesResponse})
}

func GetStory(c *gin.Context) {
	var story models.Story
	if err := config.DB.First(&story, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Story not found"})
		return
	}
	var getStoryResponse = converter.MapStoryToGetStoriesResponse(&story)
	c.JSON(http.StatusOK, gin.H{"data": getStoryResponse})
}

func UpdateStory(c *gin.Context) {
	var storyRequest models.UpdateStoryRequest
	var story models.Story
	if err := config.DB.First(&story, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Story not found"})
		return
	}

	if err := c.ShouldBindJSON(&storyRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var userIDInt = uint(userID.(float64))
	if story.UserID != userIDInt {

		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}
	story.Status = storyRequest.Status

	config.DB.Save(&story)
	var updateStoryResponse = converter.MapStoryToCreateStoryResponse(&story)
	c.JSON(http.StatusOK, gin.H{"data": updateStoryResponse})
}

func DeleteStory(c *gin.Context) {
	var story models.Story
	if err := config.DB.First(&story, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Story not found"})
		return
	}

	config.DB.Delete(&story)
	c.JSON(http.StatusOK, gin.H{"data": true})
}
