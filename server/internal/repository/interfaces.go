// Package repository defines interfaces for data persistence operations.
package repository

import (
	"context"

	"github.com/abhisheksharm-3/shrtn/internal/model"
)

// URLRepository defines operations for URL persistence.
type URLRepository interface {
	Create(ctx context.Context, url model.URL) (string, error)
	GetByShortCode(ctx context.Context, shortCode string) (*model.URL, error)
	GetAll(ctx context.Context, limit, offset int) ([]model.URL, int, error)
	UpdateClicks(ctx context.Context, docID string, clicks int) error
	Delete(ctx context.Context, docID string) error
}

// AnalyticsRepository defines operations for analytics persistence.
type AnalyticsRepository interface {
	Create(ctx context.Context, entry model.AnalyticsEntry) (string, error)
	GetByURLID(ctx context.Context, urlID string, limit, offset int) ([]model.AnalyticsEntry, error)
}
