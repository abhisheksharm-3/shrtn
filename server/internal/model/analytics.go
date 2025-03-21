package model

import "time"

type AnalyticsEntry struct {
	ID        string    `json:"id"`
	URLId     string    `json:"urlId"`
	Timestamp time.Time `json:"timestamp"`
	UserAgent string    `json:"userAgent"`
	IPAddress string    `json:"ipAddress"`
	Referer   string    `json:"referer"`
}
