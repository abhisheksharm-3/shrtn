// Package model defines domain models for the URL shortener.
package model

import "time"

// AnalyticsEntry represents a single click event.
type AnalyticsEntry struct {
	ID        string    `json:"id"`
	URLId     string    `json:"urlId"`
	Timestamp time.Time `json:"timestamp"`
	UserAgent string    `json:"userAgent"`
	IPAddress string    `json:"ipAddress"`
	Referer   string    `json:"referer"`
}

// URLStats represents aggregated statistics for a URL.
type URLStats struct {
	URL            URL            `json:"url"`
	TotalClicks    int            `json:"totalClicks"`
	LastAccessedAt time.Time      `json:"lastAccessedAt"`
	Today          int            `json:"today"`
	ThisWeek       int            `json:"thisWeek"`
	ThisMonth      int            `json:"thisMonth"`
	ReferrerStats  map[string]int `json:"referrers,omitempty"`
	BrowserStats   map[string]int `json:"browsers,omitempty"`
	CountryStats   map[string]int `json:"countries,omitempty"`
	DailyClicks    map[string]int `json:"dailyClicks,omitempty"`
	DeviceStats    map[string]int `json:"devices,omitempty"`
}
