// Package service implements business logic for the URL shortener application
package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"strings"
	"time"

	"github.com/abhisheksharm-3/shrtn/internal/appwrite"
	"github.com/abhisheksharm-3/shrtn/internal/config"
	"github.com/abhisheksharm-3/shrtn/internal/model"
)

// Common error definitions
var (
	ErrShortCodeExists = errors.New("short code already in use")
	ErrInvalidURL      = errors.New("invalid URL format")
	ErrShortCodeEmpty  = errors.New("short code cannot be empty")
)

// URLService handles business logic for URL shortening
type URLService struct {
	config   *config.Config
	dbClient *appwrite.DatabaseClient
	rng      *rand.Rand // Thread-safe random number generator
}

// NewURLService creates a new URLService
func NewURLService(cfg *config.Config) *URLService {
	dbClient := appwrite.NewDatabaseClient(cfg)

	// Create a separate RNG with its own seed
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	return &URLService{
		config:   cfg,
		dbClient: dbClient,
		rng:      rng,
	}
}

// Create creates a new shortened URL
func (s *URLService) Create(ctx context.Context, input model.URLInput) (*model.URL, error) {
	// Validate the original URL
	if err := s.validateURL(input.OriginalURL); err != nil {
		return nil, err
	}

	// Generate short code if not provided
	shortCode := input.CustomCode
	if shortCode == "" {
		shortCode = s.generateShortCode()
	} else {
		// Validate custom short code
		if len(shortCode) < 3 {
			return nil, fmt.Errorf("custom short code must be at least 3 characters long")
		}
	}

	// Check if short code already exists
	existingURL, err := s.dbClient.GetURLByShortCode(ctx, shortCode)
	if err == nil && existingURL != nil {
		return nil, ErrShortCodeExists
	}

	// Only proceed if error was "not found"
	if err != nil && !errors.Is(err, appwrite.ErrURLNotFound) {
		return nil, fmt.Errorf("error checking short code availability: %w", err)
	}

	// Create the URL entry
	now := time.Now().UTC()
	url := model.URL{
		ShortCode:   shortCode,
		OriginalURL: s.ensureURLProtocol(input.OriginalURL),
		CreatedAt:   now,
		UpdatedAt:   now,
		Clicks:      0,
	}

	// Create document in Appwrite
	id, err := s.dbClient.CreateURL(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to create URL document: %w", err)
	}

	url.ID = id
	return &url, nil
}

// GetByShortCode retrieves a URL by its short code without incrementing clicks
func (s *URLService) GetByShortCode(ctx context.Context, shortCode string) (*model.URL, error) {
	if shortCode == "" {
		return nil, ErrShortCodeEmpty
	}

	// Query Appwrite for the document with the given short code
	url, err := s.dbClient.GetURLByShortCode(ctx, shortCode)
	if err != nil {
		return nil, err
	}

	return url, nil
}

// IncrementClicks updates the click count for a URL
func (s *URLService) IncrementClicks(ctx context.Context, urlID string, currentClicks int) error {
	if urlID == "" {
		return errors.New("URL ID cannot be empty")
	}

	return s.dbClient.UpdateURLClicks(ctx, urlID, currentClicks)
}

// GetAll retrieves all URLs
func (s *URLService) GetAll(ctx context.Context) ([]model.URL, error) {
	// Query Appwrite for all documents
	urls, err := s.dbClient.GetAllURLs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URLs: %w", err)
	}

	return urls, nil
}

// validateURL checks if a URL is valid
func (s *URLService) validateURL(inputURL string) error {
	if inputURL == "" {
		return ErrInvalidURL
	}

	// Ensure it has a protocol for parsing
	urlWithProtocol := s.ensureURLProtocol(inputURL)

	// Parse the URL to validate it
	_, err := url.Parse(urlWithProtocol)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidURL, err)
	}

	return nil
}

// ensureURLProtocol ensures the URL has a protocol
func (s *URLService) ensureURLProtocol(inputURL string) string {
	if !strings.HasPrefix(inputURL, "http://") && !strings.HasPrefix(inputURL, "https://") {
		return "https://" + inputURL
	}
	return inputURL
}

// generateShortCode creates a random short code
func (s *URLService) generateShortCode() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 6

	// Lock for random number generation if using shared RNG
	code := make([]byte, length)
	for i := range code {
		code[i] = charset[s.rng.Intn(len(charset))]
	}

	return string(code)
}

// GetStatsForURL retrieves statistics for a URL
func (s *URLService) GetStatsForURL(ctx context.Context, shortCode string) (*model.URLStats, error) {
	url, err := s.GetByShortCode(ctx, shortCode)
	if err != nil {
		return nil, err
	}

	// For now, just return basic stats
	stats := &model.URLStats{
		URL:            *url,
		TotalClicks:    url.Clicks,
		LastAccessedAt: url.UpdatedAt,
	}

	return stats, nil
}

// DeleteURL removes a URL by its ID
func (s *URLService) DeleteURL(ctx context.Context, urlID string) error {
	// This method would be implemented once you add delete functionality to the DatabaseClient
	return errors.New("delete functionality not implemented")
}
