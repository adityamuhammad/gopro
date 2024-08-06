package handlers

import (
	"gopro/config"
	"gopro/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetMessage(c *gin.Context) {
	var messages []models.GetMessageResponse
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var currentUserId = uint(userID.(float64))
	var targetUserId = c.Param("id")

	query := `
        select 
            m.content, 
            m.created_at, 
            ur.name as receiver, 
            us.name as sender 
        from messages m
        join users ur on m.to_user_id = ur.id
        join users us on m.from_user_id = us.id
        where (m.from_user_id = ? and m.to_user_id = ?)
        or (m.from_user_id = ? and m.to_user_id = ?)
    `

	if err := config.DB.Raw(query, currentUserId, targetUserId, targetUserId, currentUserId).Scan(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": messages})
}
