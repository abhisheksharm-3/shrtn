package model

import (
	"time"
)

// URL represents a shortened URL
type URL struct {
	ID          string    `json:"id"`
	ShortCode   string    `json:"shortCode"`
	OriginalURL string    `json:"originalURL"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Clicks      int       `json:"clicks"`
	UserID      string    `json:"userId,omitempty"`
}

// URLInput represents the input to create a shortened URL
type URLInput struct {
	OriginalURL string `json:"originalURL" binding:"required,url"`
	CustomCode  string `json:"customCode,omitempty"`
}
