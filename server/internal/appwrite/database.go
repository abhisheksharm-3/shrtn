package appwrite

import (
	"context"
	"encoding/json"
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
			query.Equal("ShortCode", shortCode),
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
		ShortCode   string  `json:"ShortCode"`
		OriginalURL string  `json:"OriginalURL"`
		Clicks      float64 `json:"Clicks"`
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
			"Clicks":    clicks + 1,
			"UpdatedAt": time.Now(),
		}),
	)
	return err
}

// GetAllURLs retrieves all URLs
func (c *DatabaseClient) GetAllURLs(ctx context.Context, limit int, offset int) ([]model.URL, error) {
	documents, err := c.databases.ListDocuments(
		c.config.AppwriteDatabase,
		c.config.AppwriteCollection,
		c.databases.WithListDocumentsQueries([]string{}),
	)
	if err != nil {
		return nil, err
	}

	// Handle pagination manually
	totalDocs := documents.Documents
	startIdx := offset
	endIdx := offset + limit

	if startIdx >= len(totalDocs) {
		return []model.URL{}, nil
	}

	if endIdx > len(totalDocs) {
		endIdx = len(totalDocs)
	}

	paginatedDocs := totalDocs[startIdx:endIdx]

	urls := make([]model.URL, 0, len(paginatedDocs))
	for _, doc := range paginatedDocs {
		// Create a map to hold all document attributes
		var docMap map[string]interface{}

		// Convert document to JSON and then unmarshal into map
		docBytes, err := json.Marshal(doc)
		if err != nil {
			fmt.Printf("Error marshaling document %s: %v\n", doc.Id, err)
			continue
		}

		if err := json.Unmarshal(docBytes, &docMap); err != nil {
			fmt.Printf("Error unmarshaling document %s: %v\n", doc.Id, err)
			continue
		}

		// Extract values from the map
		shortCode, _ := docMap["ShortCode"].(string)
		originalURL, _ := docMap["OriginalURL"].(string)

		// Handle clicks with type assertion
		var clicks int
		if clicksVal, ok := docMap["Clicks"]; ok {
			switch v := clicksVal.(type) {
			case float64:
				clicks = int(v)
			case int:
				clicks = v
			}
		}

		// Parse dates
		var createdAt, updatedAt time.Time
		if createdAtStr, ok := docMap["CreatedAt"].(string); ok && createdAtStr != "" {
			createdAt, _ = time.Parse(time.RFC3339, createdAtStr)
		}

		if updatedAtStr, ok := docMap["UpdatedAt"].(string); ok && updatedAtStr != "" {
			updatedAt, _ = time.Parse(time.RFC3339, updatedAtStr)
		}

		url := model.URL{
			ID:          doc.Id,
			ShortCode:   shortCode,
			OriginalURL: originalURL,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
			Clicks:      clicks,
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
