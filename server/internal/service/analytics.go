package service

import (
	"context"
	"net/http"
	"time"

	"github.com/abhisheksharm-3/shrtn/internal/appwrite"
	"github.com/abhisheksharm-3/shrtn/internal/config"

	"github.com/abhisheksharm-3/shrtn/internal/model"
)

// AnalyticsService handles URL click analytics
type AnalyticsService struct {
	config   *config.Config
	dbClient *appwrite.DatabaseClient
}

// NewAnalyticsService creates a new AnalyticsService
func NewAnalyticsService(cfg *config.Config) *AnalyticsService {
	dbClient := appwrite.NewDatabaseClient(cfg)

	return &AnalyticsService{
		config:   cfg,
		dbClient: dbClient,
	}
}

// RecordClick records a click on a URL
func (s *AnalyticsService) RecordClick(ctx context.Context, urlID string, req *http.Request) error {
	// Create analytics entry
	entry := model.AnalyticsEntry{
		URLId:     urlID,
		Timestamp: time.Now(),
		UserAgent: req.UserAgent(),
		IPAddress: getIPAddress(req),
		Referer:   req.Referer(),
	}

	// Create document in Appwrite
	_, err := s.dbClient.CreateAnalyticsEntry(ctx, entry)
	return err
}

// getIPAddress extracts the client IP address from a request
func getIPAddress(r *http.Request) string {
	// Check for X-Forwarded-For header first
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		return forwarded
	}

	// Fall back to RemoteAddr
	return r.RemoteAddr
}
