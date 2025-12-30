// Package service implements business logic for the URL shortener.
package service

import (
	"context"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/abhisheksharm-3/shrtn/internal/model"
	"github.com/abhisheksharm-3/shrtn/internal/repository"
)

// AnalyticsService handles URL click analytics.
type AnalyticsService struct {
	repo         repository.AnalyticsRepository
	trustedProxy string
}

// NewAnalyticsService creates a new AnalyticsService with the given repository.
func NewAnalyticsService(repo repository.AnalyticsRepository, trustedProxy string) *AnalyticsService {
	return &AnalyticsService{
		repo:         repo,
		trustedProxy: trustedProxy,
	}
}

// RecordClick records a click event for a URL.
func (s *AnalyticsService) RecordClick(ctx context.Context, urlID string, req *http.Request) error {
	entry := model.AnalyticsEntry{
		URLId:     urlID,
		Timestamp: time.Now().UTC(),
		UserAgent: req.UserAgent(),
		IPAddress: s.extractClientIP(req),
		Referer:   req.Referer(),
	}

	_, err := s.repo.Create(ctx, entry)
	return err
}

func (s *AnalyticsService) extractClientIP(r *http.Request) string {
	if s.trustedProxy != "" {
		remoteIP, _, _ := net.SplitHostPort(r.RemoteAddr)
		if remoteIP == s.trustedProxy || remoteIP == "127.0.0.1" {
			if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
				parts := strings.Split(xff, ",")
				if len(parts) > 0 {
					return strings.TrimSpace(parts[0])
				}
			}
			if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
				return xrip
			}
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
