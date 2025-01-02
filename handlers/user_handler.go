package handlers

import (
	"context"
	"fmt"
	"gopro/config"
	"gopro/models"
	"gopro/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
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

	/*
		THIS IS FOR METHOD STORE IN LOCAL STORAGE OF PROJECT
		BEGIN
	*/
	// fileFolder := fmt.Sprintf("%s%d%s", "profile_image/", time.Now().UnixNano(), file.Filename)
	// filePath := "./storage/" + fileFolder
	// if err := c.SaveUploadedFile(file, filePath); err != nil {
	// 	c.String(http.StatusInternalServerError, "Could not save the file: %v", err)
	// 	return
	// }
	// user.ProfileImage = fileFolder
	/*
		END
	*/

	// Generate a unique file name
	fileName := fmt.Sprintf("profile_image/%d_%s", time.Now().UnixNano(), file.Filename)

	// Open the file for reading
	srcFile, err := file.Open()
	if err != nil {
		c.String(http.StatusInternalServerError, "Could not open the file: %v", err)
		return
	}
	defer srcFile.Close()

	// Upload the file to MinIO
	bucketName := "gopro"
	contentType := file.Header.Get("Content-Type")
	_, err = config.MinIOClient.PutObject(context.Background(), bucketName, fileName, srcFile, file.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		c.String(http.StatusInternalServerError, "Could not upload the file to MinIO: %v", err)
		return
	}
	user.ProfileImage = fileName

	config.DB.Save(&user)
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
