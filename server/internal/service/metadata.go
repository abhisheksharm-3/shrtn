// Package service implements business logic for the URL shortener.
package service

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// LinkPreview contains Open Graph metadata for a URL.
type LinkPreview struct {
	URL         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
	SiteName    string `json:"siteName"`
	Favicon     string `json:"favicon"`
}

// MetadataService fetches Open Graph metadata from URLs.
type MetadataService struct {
	client *http.Client
}

// NewMetadataService creates a new MetadataService.
func NewMetadataService() *MetadataService {
	return &MetadataService{
		client: &http.Client{
			Timeout: 10 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= 3 {
					return fmt.Errorf("too many redirects")
				}
				return nil
			},
		},
	}
}

// FetchPreview fetches Open Graph metadata for a URL.
func (s *MetadataService) FetchPreview(ctx context.Context, targetURL string) (*LinkPreview, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; ShrtnBot/1.0)")
	req.Header.Set("Accept", "text/html")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("URL returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1024*1024))
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	html := string(body)
	preview := &LinkPreview{
		URL: targetURL,
	}

	preview.Title = extractMeta(html, `og:title`) 
	if preview.Title == "" {
		preview.Title = extractTitle(html)
	}

	preview.Description = extractMeta(html, `og:description`)
	if preview.Description == "" {
		preview.Description = extractMeta(html, `description`)
	}

	preview.Image = extractMeta(html, `og:image`)
	preview.SiteName = extractMeta(html, `og:site_name`)
	preview.Favicon = extractFavicon(html, targetURL)

	return preview, nil
}

func extractMeta(html, property string) string {
	patterns := []string{
		fmt.Sprintf(`<meta[^>]*property=["']%s["'][^>]*content=["']([^"']*)["']`, property),
		fmt.Sprintf(`<meta[^>]*content=["']([^"']*)["'][^>]*property=["']%s["']`, property),
		fmt.Sprintf(`<meta[^>]*name=["']%s["'][^>]*content=["']([^"']*)["']`, property),
		fmt.Sprintf(`<meta[^>]*content=["']([^"']*)["'][^>]*name=["']%s["']`, property),
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(html)
		if len(matches) > 1 {
			return strings.TrimSpace(matches[1])
		}
	}
	return ""
}

func extractTitle(html string) string {
	re := regexp.MustCompile(`<title[^>]*>([^<]*)</title>`)
	matches := re.FindStringSubmatch(html)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}

func extractFavicon(html, baseURL string) string {
	patterns := []string{
		`<link[^>]*rel=["'](?:shortcut )?icon["'][^>]*href=["']([^"']*)["']`,
		`<link[^>]*href=["']([^"']*)["'][^>]*rel=["'](?:shortcut )?icon["']`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(html)
		if len(matches) > 1 {
			favicon := matches[1]
			if strings.HasPrefix(favicon, "//") {
				return "https:" + favicon
			}
			if strings.HasPrefix(favicon, "/") {
				parts := strings.SplitN(baseURL, "/", 4)
				if len(parts) >= 3 {
					return parts[0] + "//" + parts[2] + favicon
				}
			}
			if !strings.HasPrefix(favicon, "http") {
				parts := strings.SplitN(baseURL, "/", 4)
				if len(parts) >= 3 {
					return parts[0] + "//" + parts[2] + "/" + favicon
				}
			}
			return favicon
		}
	}

	parts := strings.SplitN(baseURL, "/", 4)
	if len(parts) >= 3 {
		return parts[0] + "//" + parts[2] + "/favicon.ico"
	}
	return ""
}
