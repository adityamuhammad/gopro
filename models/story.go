package models

import "time"

type Story struct {
	ID        uint      `gorm:"primaryKey"`
	Status    string    `json:"status"`
	UserID    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      User      `gorm:"foreignKey:UserID"`
}

type CreateStoryRequest struct {
	Status string `json:"status"`
}

type CreateStoryResponse struct {
	ID        uint      `json:"id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateStoryRequest struct {
	Status string `json:"status"`
}

type UpdateStoryResponse struct {
	ID        uint      `json:"id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GetStoriesResponse struct {
	ID        uint      `json:"id"`
	Status    string    `json:"status"`
	UserID    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
