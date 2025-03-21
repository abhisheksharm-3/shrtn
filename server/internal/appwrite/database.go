package appwrite

import (
	"context"
	"fmt"
	"time"

	"github.com/abhisheksharm-3/shrtn/internal/config"
	"github.com/abhisheksharm-3/shrtn/internal/model"

	"github.com/appwrite/sdk-for-go/databases"
	"github.com/appwrite/sdk-for-go/id"
	"github.com/appwrite/sdk-for-go/query"
)

// DatabaseClient encapsulates Appwrite database operations
type DatabaseClient struct {
	config    *config.Config
	databases *databases.Databases
}

// NewDatabaseClient creates a new DatabaseClient
func NewDatabaseClient(cfg *config.Config) *DatabaseClient {
	client := GetClient(cfg)
	databasesService := databases.New(client)

	return &DatabaseClient{
		config:    cfg,
		databases: databasesService,
	}
}

// CreateURL creates a new URL document in Appwrite
func (c *DatabaseClient) CreateURL(ctx context.Context, url model.URL) (string, error) {
	document, err := c.databases.CreateDocument(
		c.config.AppwriteDatabase,
		c.config.AppwriteCollection,
		id.Unique(),
		map[string]interface{}{
			"shortCode":   url.ShortCode,
			"originalURL": url.OriginalURL,
			"createdAt":   url.CreatedAt,
			"updatedAt":   url.UpdatedAt,
			"clicks":      url.Clicks,
		},
	)
	if err != nil {
		return "", err
	}

	return document.Id, nil
}

// GetURLByShortCode retrieves a URL by its short code
func (c *DatabaseClient) GetURLByShortCode(ctx context.Context, shortCode string) (*model.URL, error) {
	documents, err := c.databases.ListDocuments(
		c.config.AppwriteDatabase,
		c.config.AppwriteCollection,
		c.databases.WithListDocumentsQueries([]string{
			query.Equal("shortCode", shortCode),
		}),
	)
	if err != nil {
		return nil, err
	}

	if len(documents.Documents) == 0 {
		return nil, fmt.Errorf("URL not found")
	}

	doc := documents.Documents[0]
	var urlData struct {
		ShortCode   string  `json:"shortCode"`
		OriginalURL string  `json:"originalURL"`
		Clicks      float64 `json:"clicks"`
	}

	doc.Decode(&urlData)

	url := model.URL{
		ID:          doc.Id,
		ShortCode:   urlData.ShortCode,
		OriginalURL: urlData.OriginalURL,
		Clicks:      int(urlData.Clicks),
	}

	return &url, nil
}

// UpdateURLClicks increments the click count for a URL
func (c *DatabaseClient) UpdateURLClicks(ctx context.Context, docID string, clicks int) error {
	_, err := c.databases.UpdateDocument(
		c.config.AppwriteDatabase,
		c.config.AppwriteCollection,
		docID,
		c.databases.WithUpdateDocumentData(map[string]interface{}{
			"clicks":    clicks + 1,
			"updatedAt": time.Now(),
		}),
	)
	return err
}

// GetAllURLs retrieves all URLs
func (c *DatabaseClient) GetAllURLs(ctx context.Context, limit int, offset int) ([]model.URL, error) {
	// Since WithListDocumentsLimit and WithListDocumentsOffset aren't available,
	// we need to handle pagination differently or accept the default values

	documents, err := c.databases.ListDocuments(
		c.config.AppwriteDatabase,
		c.config.AppwriteCollection,
		c.databases.WithListDocumentsQueries([]string{}),
	)
	if err != nil {
		return nil, err
	}

	// Since we can't set limit and offset directly, we'll handle it manually in-memory
	// This is not ideal for large collections but works as a workaround
	totalDocs := documents.Documents
	startIdx := offset
	endIdx := offset + limit

	if startIdx >= len(totalDocs) {
		return []model.URL{}, nil
	}

	if endIdx > len(totalDocs) {
		endIdx = len(totalDocs)
	}

	// Get the paginated subset
	paginatedDocs := totalDocs[startIdx:endIdx]

	urls := make([]model.URL, 0, len(paginatedDocs))
	for _, doc := range paginatedDocs {
		var data struct {
			ShortCode   string  `json:"shortCode"`
			OriginalURL string  `json:"originalURL"`
			CreatedAt   string  `json:"createdAt"`
			UpdatedAt   string  `json:"updatedAt"`
			Clicks      float64 `json:"clicks"`
		}

		doc.Decode(&data)

		createdAt, _ := time.Parse(time.RFC3339, data.CreatedAt)
		updatedAt, _ := time.Parse(time.RFC3339, data.UpdatedAt)

		url := model.URL{
			ID:          doc.Id,
			ShortCode:   data.ShortCode,
			OriginalURL: data.OriginalURL,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
			Clicks:      int(data.Clicks),
		}
		urls = append(urls, url)
	}

	return urls, nil
}

// CreateAnalyticsEntry records a click analytics entry
func (c *DatabaseClient) CreateAnalyticsEntry(ctx context.Context, entry model.AnalyticsEntry) (string, error) {
	document, err := c.databases.CreateDocument(
		c.config.AppwriteDatabase,
		"analytics", // Assuming you have an "analytics" collection
		id.Unique(),
		map[string]interface{}{
			"urlId":     entry.URLId,
			"timestamp": entry.Timestamp,
			"userAgent": entry.UserAgent,
			"ipAddress": entry.IPAddress,
			"referer":   entry.Referer,
		},
	)
	if err != nil {
		return "", err
	}

	return document.Id, nil
}
