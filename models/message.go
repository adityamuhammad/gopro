package models

import "time"

type Message struct {
	ID         uint      `gorm:"primaryKey"`
	ToUserID   uint      `json:"to_user_id"`
	FromUserId uint      `json:"from_user_id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type GetMessageResponse struct {
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	Receiver  string    `json:"receiver"`
	Sender    string    `json:"sender"`
}
