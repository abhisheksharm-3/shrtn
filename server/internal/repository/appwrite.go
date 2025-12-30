// Package repository provides Appwrite implementation for data persistence.
package repository

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/abhisheksharm-3/shrtn/internal/config"
	"github.com/abhisheksharm-3/shrtn/internal/model"

	"github.com/appwrite/sdk-for-go/appwrite"
	"github.com/appwrite/sdk-for-go/client"
	"github.com/appwrite/sdk-for-go/databases"
	"github.com/appwrite/sdk-for-go/id"
	"github.com/appwrite/sdk-for-go/models"
	"github.com/appwrite/sdk-for-go/query"
)

var (
	ErrURLNotFound = errors.New("url not found")
	ErrDecoding    = errors.New("error decoding response")
)

const (
	collectionAnalytics = "analytics"
	defaultTimeout      = 10 * time.Second
)

type urlDocument struct {
	ID          string  `json:"$id"`
	ShortCode   string  `json:"ShortCode"`
	OriginalURL string  `json:"OriginalURL"`
	CreatedAt   string  `json:"CreatedAt"`
	UpdatedAt   string  `json:"UpdatedAt"`
	Clicks      float64 `json:"Clicks"`
}

type urlDocumentList struct {
	*models.DocumentList
	Documents []urlDocument `json:"documents"`
}

// AppwriteClient manages singleton Appwrite client connection.
type AppwriteClient struct {
	client client.Client
}

var (
	appwriteInstance *AppwriteClient
	appwriteOnce     sync.Once
)

// GetAppwriteClient returns singleton Appwrite client instance.
func GetAppwriteClient(cfg *config.Config) *AppwriteClient {
	appwriteOnce.Do(func() {
		appwriteInstance = &AppwriteClient{
			client: appwrite.NewClient(
				appwrite.WithProject(cfg.AppwriteProjectID),
				appwrite.WithKey(cfg.AppwriteAPIKey),
			),
		}
	})
	return appwriteInstance
}

// AppwriteURLRepository implements URLRepository using Appwrite.
type AppwriteURLRepository struct {
	config    *config.Config
	databases *databases.Databases
}

// NewAppwriteURLRepository creates a new Appwrite URL repository.
func NewAppwriteURLRepository(cfg *config.Config) *AppwriteURLRepository {
	awClient := GetAppwriteClient(cfg)
	return &AppwriteURLRepository{
		config:    cfg,
		databases: databases.New(awClient.client),
	}
}

// Create inserts a new URL document and returns its ID.
func (r *AppwriteURLRepository) Create(ctx context.Context, url model.URL) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	uniqueID := id.Unique()
	document, err := r.databases.CreateDocument(
		r.config.AppwriteDatabase,
		r.config.AppwriteCollection,
		uniqueID,
		map[string]interface{}{
			"ID":          uniqueID,
			"ShortCode":   url.ShortCode,
			"OriginalURL": url.OriginalURL,
			"CreatedAt":   url.CreatedAt.Format(time.RFC3339),
			"UpdatedAt":   url.UpdatedAt.Format(time.RFC3339),
			"Clicks":      url.Clicks,
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to create URL document: %w", err)
	}

	return document.Id, nil
}

// GetByShortCode retrieves a URL by its short code.
func (r *AppwriteURLRepository) GetByShortCode(ctx context.Context, shortCode string) (*model.URL, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	if shortCode == "" {
		return nil, fmt.Errorf("short code cannot be empty")
	}

	response, err := r.databases.ListDocuments(
		r.config.AppwriteDatabase,
		r.config.AppwriteCollection,
		r.databases.WithListDocumentsQueries([]string{
			query.Equal("ShortCode", shortCode),
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query URL by short code: %w", err)
	}

	var urlList urlDocumentList
	if err := response.Decode(&urlList); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDecoding, err)
	}

	if len(urlList.Documents) == 0 {
		return nil, ErrURLNotFound
	}

	return documentToURL(urlList.Documents[0]), nil
}

// GetAll retrieves paginated URLs and total count.
func (r *AppwriteURLRepository) GetAll(ctx context.Context, limit, offset int) ([]model.URL, int, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	queries := []string{
		query.Limit(limit),
		query.Offset(offset),
		query.OrderDesc("CreatedAt"),
	}

	response, err := r.databases.ListDocuments(
		r.config.AppwriteDatabase,
		r.config.AppwriteCollection,
		r.databases.WithListDocumentsQueries(queries),
	)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list URL documents: %w", err)
	}

	var urlList urlDocumentList
	if err := response.Decode(&urlList); err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrDecoding, err)
	}

	urls := make([]model.URL, 0, len(urlList.Documents))
	for _, doc := range urlList.Documents {
		urls = append(urls, *documentToURL(doc))
	}

	return urls, urlList.Total, nil
}

// UpdateClicks increments the click count for a URL.
func (r *AppwriteURLRepository) UpdateClicks(ctx context.Context, docID string, clicks int) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	if docID == "" {
		return fmt.Errorf("document ID cannot be empty")
	}

	_, err := r.databases.UpdateDocument(
		r.config.AppwriteDatabase,
		r.config.AppwriteCollection,
		docID,
		r.databases.WithUpdateDocumentData(map[string]interface{}{
			"Clicks":    clicks + 1,
			"UpdatedAt": time.Now().UTC().Format(time.RFC3339),
		}),
	)
	if err != nil {
		return fmt.Errorf("failed to update URL clicks: %w", err)
	}

	return nil
}

// Delete removes a URL document by ID.
func (r *AppwriteURLRepository) Delete(ctx context.Context, docID string) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	if docID == "" {
		return fmt.Errorf("document ID cannot be empty")
	}

	_, err := r.databases.DeleteDocument(
		r.config.AppwriteDatabase,
		r.config.AppwriteCollection,
		docID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete URL document: %w", err)
	}

	return nil
}

// AppwriteAnalyticsRepository implements AnalyticsRepository using Appwrite.
type AppwriteAnalyticsRepository struct {
	config    *config.Config
	databases *databases.Databases
}

// NewAppwriteAnalyticsRepository creates a new Appwrite analytics repository.
func NewAppwriteAnalyticsRepository(cfg *config.Config) *AppwriteAnalyticsRepository {
	awClient := GetAppwriteClient(cfg)
	return &AppwriteAnalyticsRepository{
		config:    cfg,
		databases: databases.New(awClient.client),
	}
}

// Create inserts a new analytics entry and returns its ID.
func (r *AppwriteAnalyticsRepository) Create(ctx context.Context, entry model.AnalyticsEntry) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	if entry.URLId == "" {
		return "", fmt.Errorf("URL ID cannot be empty for analytics entry")
	}

	document, err := r.databases.CreateDocument(
		r.config.AppwriteDatabase,
		collectionAnalytics,
		id.Unique(),
		map[string]interface{}{
			"urlId":     entry.URLId,
			"timestamp": entry.Timestamp.Format(time.RFC3339),
			"userAgent": entry.UserAgent,
			"ipAddress": entry.IPAddress,
			"referer":   entry.Referer,
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to create analytics entry: %w", err)
	}

	return document.Id, nil
}

// GetByURLID retrieves analytics entries for a URL with pagination.
func (r *AppwriteAnalyticsRepository) GetByURLID(ctx context.Context, urlID string, limit, offset int) ([]model.AnalyticsEntry, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	if urlID == "" {
		return nil, fmt.Errorf("URL ID cannot be empty")
	}

	queries := []string{
		query.Equal("urlId", urlID),
		query.Limit(limit),
		query.Offset(offset),
		query.OrderDesc("timestamp"),
	}

	response, err := r.databases.ListDocuments(
		r.config.AppwriteDatabase,
		collectionAnalytics,
		r.databases.WithListDocumentsQueries(queries),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query analytics: %w", err)
	}

	var result struct {
		Documents []struct {
			ID        string `json:"$id"`
			URLId     string `json:"urlId"`
			Timestamp string `json:"timestamp"`
			UserAgent string `json:"userAgent"`
			IPAddress string `json:"ipAddress"`
			Referer   string `json:"referer"`
		} `json:"documents"`
	}
	if err := response.Decode(&result); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDecoding, err)
	}

	entries := make([]model.AnalyticsEntry, 0, len(result.Documents))
	for _, doc := range result.Documents {
		timestamp, _ := time.Parse(time.RFC3339, doc.Timestamp)
		entries = append(entries, model.AnalyticsEntry{
			ID:        doc.ID,
			URLId:     doc.URLId,
			Timestamp: timestamp,
			UserAgent: doc.UserAgent,
			IPAddress: doc.IPAddress,
			Referer:   doc.Referer,
		})
	}

	return entries, nil
}

func documentToURL(doc urlDocument) *model.URL {
	var createdAt, updatedAt time.Time
	if doc.CreatedAt != "" {
		createdAt, _ = time.Parse(time.RFC3339, doc.CreatedAt)
	}
	if doc.UpdatedAt != "" {
		updatedAt, _ = time.Parse(time.RFC3339, doc.UpdatedAt)
	}

	return &model.URL{
		ID:          doc.ID,
		ShortCode:   doc.ShortCode,
		OriginalURL: doc.OriginalURL,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		Clicks:      int(doc.Clicks),
	}
}
