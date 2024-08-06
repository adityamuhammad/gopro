package handlers

import (
	"encoding/json"
	"fmt"
	"gopro/config"
	"gopro/models"
	"gopro/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func RealTimeMessage(c *gin.Context) {
	conn, err := utils.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set WebSocket upgrade"})
		return
	}
	defer conn.Close()

	userID, _ := c.Get("userID")
	userIDUint := uint(userID.(float64))

	utils.Registry.Add(userIDUint, conn)
	defer utils.Registry.Remove(userIDUint)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("Error: %v\n", err)
			}
			break
		}
		fmt.Printf("Received message from user %v: %s\n", userIDUint, message)

		var msg utils.Message
		if err := json.Unmarshal(message, &msg); err != nil {
			fmt.Printf("Error unmarshalling message: %v\n", err)
			continue
		}

		var messageModel models.Message
		messageModel.FromUserId = userIDUint
		messageModel.ToUserID = msg.ToUserID
		messageModel.Content = msg.Content

		config.DB.Create(&messageModel)

		if toConn, exists := utils.Registry.Get(msg.ToUserID); exists {
			if err := toConn.WriteMessage(websocket.TextMessage, []byte(msg.Content)); err != nil {
				fmt.Printf("Error sending message to user %v: %v\n", msg.ToUserID, err)
			}
		}
	}
}
