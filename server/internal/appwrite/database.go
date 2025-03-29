// Package appwrite provides Appwrite database integration for the URL shortener service
package appwrite

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/abhisheksharm-3/shrtn/internal/config"
	"github.com/abhisheksharm-3/shrtn/internal/model"

	"github.com/appwrite/sdk-for-go/databases"
	"github.com/appwrite/sdk-for-go/id"
	"github.com/appwrite/sdk-for-go/models"
	"github.com/appwrite/sdk-for-go/query"
)

// Common error definitions
var (
	ErrURLNotFound = errors.New("URL not found")
	ErrDecoding    = errors.New("error decoding Appwrite response")
)

// Collection names
const (
	CollectionAnalytics = "analytics"
)

// URLDocument represents the URL document structure in Appwrite
type URLDocument struct {
	ID          string  `json:"$id"`
	ShortCode   string  `json:"ShortCode"`
	OriginalURL string  `json:"OriginalURL"`
	CreatedAt   string  `json:"CreatedAt"`
	UpdatedAt   string  `json:"UpdatedAt"`
	Clicks      float64 `json:"Clicks"`
}

// URLList represents a list of URL documents from Appwrite
type URLList struct {
	*models.DocumentList
	Documents []URLDocument `json:"documents"`
}

// DatabaseClient encapsulates Appwrite database operations
type DatabaseClient struct {
	config    *config.Config
	databases *databases.Databases
}

// NewDatabaseClient creates a new DatabaseClient with the provided configuration
func NewDatabaseClient(cfg *config.Config) *DatabaseClient {
	client := GetClient(cfg)
	databasesService := databases.New(client)

	return &DatabaseClient{
		config:    cfg,
		databases: databasesService,
	}
}

// CreateURL creates a new URL document in Appwrite and returns its ID
func (c *DatabaseClient) CreateURL(ctx context.Context, url model.URL) (string, error) {
	uniqueID := id.Unique()
	document, err := c.databases.CreateDocument(
		c.config.AppwriteDatabase,
		c.config.AppwriteCollection,
		uniqueID,
		map[string]interface{}{
			"ID":          uniqueID,
			"ShortCode":   url.ShortCode,
			"OriginalURL": url.OriginalURL,
			"CreatedAt":   url.CreatedAt,
			"UpdatedAt":   url.UpdatedAt,
			"Clicks":      url.Clicks,
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to create URL document: %w", err)
	}

	return document.Id, nil
}

// GetURLByShortCode retrieves a URL by its short code
func (c *DatabaseClient) GetURLByShortCode(ctx context.Context, shortCode string) (*model.URL, error) {
	if shortCode == "" {
		return nil, fmt.Errorf("short code cannot be empty")
	}

	response, err := c.databases.ListDocuments(
		c.config.AppwriteDatabase,
		c.config.AppwriteCollection,
		c.databases.WithListDocumentsQueries([]string{
			query.Equal("ShortCode", shortCode),
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query URL by short code: %w", err)
	}

	var urlList URLList
	if err := response.Decode(&urlList); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDecoding, err)
	}

	if len(urlList.Documents) == 0 {
		return nil, ErrURLNotFound
	}

	return documentToURL(urlList.Documents[0]), nil
}

// UpdateURLClicks increments the click count for a URL
func (c *DatabaseClient) UpdateURLClicks(ctx context.Context, docID string, clicks int) error {
	if docID == "" {
		return fmt.Errorf("document ID cannot be empty")
	}

	_, err := c.databases.UpdateDocument(
		c.config.AppwriteDatabase,
		c.config.AppwriteCollection,
		docID,
		c.databases.WithUpdateDocumentData(map[string]interface{}{
			"Clicks":    clicks + 1,
			"UpdatedAt": time.Now().Format(time.RFC3339),
		}),
	)
	if err != nil {
		return fmt.Errorf("failed to update URL clicks: %w", err)
	}

	return nil
}

// GetAllURLs retrieves all URLs from the database
func (c *DatabaseClient) GetAllURLs(ctx context.Context) ([]model.URL, error) {
	response, err := c.databases.ListDocuments(
		c.config.AppwriteDatabase,
		c.config.AppwriteCollection,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list URL documents: %w", err)
	}

	var urlList URLList
	if err := response.Decode(&urlList); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDecoding, err)
	}

	urls := make([]model.URL, 0, len(urlList.Documents))
	for _, doc := range urlList.Documents {
		urls = append(urls, *documentToURL(doc))
	}

	return urls, nil
}

// CreateAnalyticsEntry records a click analytics entry and returns its ID
func (c *DatabaseClient) CreateAnalyticsEntry(ctx context.Context, entry model.AnalyticsEntry) (string, error) {
	if entry.URLId == "" {
		return "", fmt.Errorf("URL ID cannot be empty for analytics entry")
	}

	document, err := c.databases.CreateDocument(
		c.config.AppwriteDatabase,
		CollectionAnalytics,
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

// documentToURL converts an Appwrite URLDocument to a model.URL
func documentToURL(doc URLDocument) *model.URL {
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
