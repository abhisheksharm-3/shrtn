package model

import (
	"time"
)

// URL represents a shortened URL
type URL struct {
	ID          string    `json:"ID"`
	ShortCode   string    `json:"ShortCode"`
	OriginalURL string    `json:"OriginalURL"`
	CreatedAt   time.Time `json:"CreatedAt"`
	UpdatedAt   time.Time `json:"UpdatedAt"`
	Clicks      int       `json:"Clicks"`
	UserID      string    `json:"UserID,omitempty"`
}

// URLInput represents the input to create a shortened URL
type URLInput struct {
	OriginalURL string `json:"originalURL" binding:"required,url"`
	CustomCode  string `json:"customCode,omitempty"`
}
