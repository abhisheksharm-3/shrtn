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
type URLStats struct {
	URL            URL       `json:"url"`
	TotalClicks    int       `json:"totalClicks"`
	LastAccessedAt time.Time `json:"lastAccessedAt"`
	CreatedAt      time.Time `json:"createdAt"`

	// Additional statistics fields
	Today         int            `json:"today"`
	ThisWeek      int            `json:"thisWeek"`
	ThisMonth     int            `json:"thisMonth"`
	ReferrerStats map[string]int `json:"referrers,omitempty"`
	BrowserStats  map[string]int `json:"browsers,omitempty"`
	CountryStats  map[string]int `json:"countries,omitempty"`
	DailyClicks   map[string]int `json:"dailyClicks,omitempty"`
	DeviceStats   map[string]int `json:"devices,omitempty"`
}

// StatsPeriod represents a time period for statistics
type StatsPeriod string

const (
	PeriodDay   StatsPeriod = "day"
	PeriodWeek  StatsPeriod = "week"
	PeriodMonth StatsPeriod = "month"
	PeriodYear  StatsPeriod = "year"
	PeriodAll   StatsPeriod = "all"
)

// StatsRequest represents a request for URL statistics
type StatsRequest struct {
	ShortCode string      `json:"shortCode"`
	Period    StatsPeriod `json:"period"`
	From      time.Time   `json:"from,omitempty"`
	To        time.Time   `json:"to,omitempty"`
}
