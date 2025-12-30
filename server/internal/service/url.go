// Package service implements business logic for the URL shortener.
package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"net"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/abhisheksharm-3/shrtn/internal/model"
	"github.com/abhisheksharm-3/shrtn/internal/repository"
)

var (
	ErrShortCodeExists   = errors.New("short code already in use")
	ErrInvalidURL        = errors.New("invalid URL format")
	ErrShortCodeEmpty    = errors.New("short code cannot be empty")
	ErrShortCodeTooShort = errors.New("short code must be at least 3 characters")
	ErrShortCodeInvalid  = errors.New("short code contains invalid characters")
	ErrURLBlocked        = errors.New("URL is not allowed")
)

const (
	shortCodeCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	shortCodeLength  = 6
	minCustomLength  = 3
	maxCustomLength  = 20
)

var (
	shortCodeRegex  = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	blockedPrefixes = []string{"javascript:", "data:", "vbscript:", "file:"}
	reservedCodes   = map[string]bool{
		"api": true, "admin": true, "health": true, "www": true,
		"static": true, "assets": true, "favicon": true,
	}
)

// URLService handles business logic for URL shortening.
type URLService struct {
	repo repository.URLRepository
}

// NewURLService creates a new URLService with the given repository.
func NewURLService(repo repository.URLRepository) *URLService {
	return &URLService{repo: repo}
}

// Create creates a new shortened URL.
func (s *URLService) Create(ctx context.Context, input model.URLInput) (*model.URL, error) {
	normalizedURL, err := s.validateAndNormalizeURL(input.OriginalURL)
	if err != nil {
		return nil, err
	}

	shortCode := input.CustomCode
	if shortCode == "" {
		shortCode, err = s.generateShortCode()
		if err != nil {
			return nil, fmt.Errorf("failed to generate short code: %w", err)
		}
	} else {
		if err := s.validateCustomCode(shortCode); err != nil {
			return nil, err
		}
	}

	existingURL, err := s.repo.GetByShortCode(ctx, shortCode)
	if err == nil && existingURL != nil {
		return nil, ErrShortCodeExists
	}
	if err != nil && !errors.Is(err, repository.ErrURLNotFound) {
		return nil, fmt.Errorf("error checking short code availability: %w", err)
	}

	now := time.Now().UTC()
	newURL := model.URL{
		ShortCode:   shortCode,
		OriginalURL: normalizedURL,
		CreatedAt:   now,
		UpdatedAt:   now,
		Clicks:      0,
	}

	id, err := s.repo.Create(ctx, newURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create URL: %w", err)
	}

	newURL.ID = id
	return &newURL, nil
}

// GetByShortCode retrieves a URL by its short code.
func (s *URLService) GetByShortCode(ctx context.Context, shortCode string) (*model.URL, error) {
	if shortCode == "" {
		return nil, ErrShortCodeEmpty
	}
	return s.repo.GetByShortCode(ctx, shortCode)
}

// GetAll retrieves paginated URLs.
func (s *URLService) GetAll(ctx context.Context, limit, offset int) (*model.URLListResponse, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	urls, total, err := s.repo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URLs: %w", err)
	}

	return &model.URLListResponse{
		URLs:   urls,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}

// IncrementClicks updates the click count for a URL.
func (s *URLService) IncrementClicks(ctx context.Context, urlID string, currentClicks int) error {
	if urlID == "" {
		return errors.New("URL ID cannot be empty")
	}
	return s.repo.UpdateClicks(ctx, urlID, currentClicks)
}

// Delete removes a URL by its ID.
func (s *URLService) Delete(ctx context.Context, urlID string) error {
	if urlID == "" {
		return errors.New("URL ID cannot be empty")
	}
	return s.repo.Delete(ctx, urlID)
}

func (s *URLService) validateAndNormalizeURL(inputURL string) (string, error) {
	if inputURL == "" {
		return "", ErrInvalidURL
	}

	normalized := strings.TrimSpace(inputURL)
	if !strings.HasPrefix(normalized, "http://") && !strings.HasPrefix(normalized, "https://") {
		normalized = "https://" + normalized
	}

	lowerURL := strings.ToLower(normalized)
	for _, prefix := range blockedPrefixes {
		if strings.HasPrefix(lowerURL, prefix) {
			return "", ErrURLBlocked
		}
	}

	parsed, err := url.Parse(normalized)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidURL, err)
	}

	if parsed.Host == "" {
		return "", ErrInvalidURL
	}

	if ip := net.ParseIP(parsed.Hostname()); ip != nil {
		if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() {
			return "", ErrURLBlocked
		}
	}

	return normalized, nil
}

func (s *URLService) validateCustomCode(code string) error {
	if len(code) < minCustomLength {
		return ErrShortCodeTooShort
	}
	if len(code) > maxCustomLength {
		return fmt.Errorf("short code must not exceed %d characters", maxCustomLength)
	}
	if !shortCodeRegex.MatchString(code) {
		return ErrShortCodeInvalid
	}
	if reservedCodes[strings.ToLower(code)] {
		return fmt.Errorf("short code '%s' is reserved", code)
	}
	return nil
}

func (s *URLService) generateShortCode() (string, error) {
	code := make([]byte, shortCodeLength)
	charsetLen := big.NewInt(int64(len(shortCodeCharset)))

	for i := range code {
		n, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", err
		}
		code[i] = shortCodeCharset[n.Int64()]
	}

	return string(code), nil
}
