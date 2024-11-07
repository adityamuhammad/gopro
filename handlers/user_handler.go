package handlers

import (
	"fmt"
	"gopro/config"
	"gopro/models"
	"gopro/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func UpdateProfileImage(c *gin.Context) {
	file, err := c.FormFile("profile_image")
	if err != nil {
		c.String(http.StatusBadRequest, "Could not get the file: %v", err)
		return
	}

	if !utils.IsAllowedImageExtension(file.Filename) {
		c.String(http.StatusBadRequest, "Only JPEG, PNG, or GIF images are allowed")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	fileFolder := fmt.Sprintf("%s%d%s", "profile_image/", time.Now().UnixNano(), file.Filename)
	filePath := "./storage/" + fileFolder
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.String(http.StatusInternalServerError, "Could not save the file: %v", err)
		return
	}
	user.ProfileImage = fileFolder

	config.DB.Save(&user)
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
