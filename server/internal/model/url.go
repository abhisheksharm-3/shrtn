// Package model defines domain models for the URL shortener.
package model

import "time"

// URL represents a shortened URL.
type URL struct {
	ID          string    `json:"id"`
	ShortCode   string    `json:"shortCode"`
	OriginalURL string    `json:"originalUrl"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Clicks      int       `json:"clicks"`
	UserID      string    `json:"userId,omitempty"`
}

// URLInput represents the input to create a shortened URL.
type URLInput struct {
	OriginalURL string `json:"originalUrl" binding:"required"`
	CustomCode  string `json:"customCode,omitempty"`
}

// URLListResponse represents a paginated list of URLs.
type URLListResponse struct {
	URLs   []URL `json:"urls"`
	Total  int   `json:"total"`
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
}
