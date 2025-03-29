package service

import (
	"context"
	"errors"
	"math/rand"

	"github.com/abhisheksharm-3/shrtn/internal/appwrite"
	"github.com/abhisheksharm-3/shrtn/internal/config"

	"time"

	"github.com/abhisheksharm-3/shrtn/internal/model"
)

// URLService handles business logic for URL shortening
type URLService struct {
	config   *config.Config
	dbClient *appwrite.DatabaseClient
}

// NewURLService creates a new URLService
func NewURLService(cfg *config.Config) *URLService {
	dbClient := appwrite.NewDatabaseClient(cfg)

	return &URLService{
		config:   cfg,
		dbClient: dbClient,
	}
}

// Create creates a new shortened URL
func (s *URLService) Create(ctx context.Context, input model.URLInput) (*model.URL, error) {
	// Generate short code if not provided
	shortCode := input.CustomCode
	if shortCode == "" {
		shortCode = generateShortCode()
	}

	// Check if short code already exists
	_, err := s.GetByShortCode(ctx, shortCode)
	if err == nil {
		// Short code already exists
		return nil, errors.New("short code already in use")
	}

	now := time.Now()
	url := model.URL{
		ShortCode:   shortCode,
		OriginalURL: input.OriginalURL,
		CreatedAt:   now,
		UpdatedAt:   now,
		Clicks:      0,
	}

	// Create document in Appwrite
	id, err := s.dbClient.CreateURL(ctx, url)
	if err != nil {
		return nil, err
	}

	url.ID = id

	return &url, nil
}

// GetByShortCode retrieves a URL by its short code
func (s *URLService) GetByShortCode(ctx context.Context, shortCode string) (*model.URL, error) {
	// Query Appwrite for the document with the given short code
	url, err := s.dbClient.GetURLByShortCode(ctx, shortCode)
	if err != nil {
		return nil, err
	}

	// Increment click count
	err = s.dbClient.UpdateURLClicks(ctx, url.ID, url.Clicks)
	if err != nil {
		// Non-critical error, just log it
		// Log error here
	} else {
		url.Clicks++
	}

	return url, nil
}

// GetAll retrieves all URLs
func (s *URLService) GetAll(ctx context.Context) ([]model.URL, error) {
	// Query Appwrite for all documents
	return s.dbClient.GetAllURLs(ctx)
}

// Generate a random short code
func generateShortCode() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 6

	rand.Seed(time.Now().UnixNano())
	code := make([]byte, length)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}

	return string(code)
}
