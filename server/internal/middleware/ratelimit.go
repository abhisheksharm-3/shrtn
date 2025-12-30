// Package middleware provides HTTP middleware for the API server.
package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiterConfig configures the rate limiter behavior.
type RateLimiterConfig struct {
	RequestsPerMinute int
	BurstSize         int
	CleanupInterval   time.Duration
}

// DefaultRateLimiterConfig returns sensible default configuration.
func DefaultRateLimiterConfig() RateLimiterConfig {
	return RateLimiterConfig{
		RequestsPerMinute: 60,
		BurstSize:         10,
		CleanupInterval:   5 * time.Minute,
	}
}

type clientRecord struct {
	tokens     float64
	lastRefill time.Time
}

type rateLimiter struct {
	mu       sync.RWMutex
	clients  map[string]*clientRecord
	config   RateLimiterConfig
	stopChan chan struct{}
}

func newRateLimiter(cfg RateLimiterConfig) *rateLimiter {
	rl := &rateLimiter{
		clients:  make(map[string]*clientRecord),
		config:   cfg,
		stopChan: make(chan struct{}),
	}
	go rl.cleanup()
	return rl
}

func (rl *rateLimiter) allow(clientIP string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	record, exists := rl.clients[clientIP]

	if !exists {
		rl.clients[clientIP] = &clientRecord{
			tokens:     float64(rl.config.BurstSize - 1),
			lastRefill: now,
		}
		return true
	}

	elapsed := now.Sub(record.lastRefill)
	refillRate := float64(rl.config.RequestsPerMinute) / 60.0
	record.tokens += elapsed.Seconds() * refillRate

	if record.tokens > float64(rl.config.BurstSize) {
		record.tokens = float64(rl.config.BurstSize)
	}

	record.lastRefill = now

	if record.tokens >= 1 {
		record.tokens--
		return true
	}

	return false
}

func (rl *rateLimiter) cleanup() {
	ticker := time.NewTicker(rl.config.CleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.mu.Lock()
			threshold := time.Now().Add(-rl.config.CleanupInterval)
			for ip, record := range rl.clients {
				if record.lastRefill.Before(threshold) {
					delete(rl.clients, ip)
				}
			}
			rl.mu.Unlock()
		case <-rl.stopChan:
			return
		}
	}
}

func (rl *rateLimiter) stop() {
	close(rl.stopChan)
}

// RateLimiter returns middleware that limits requests per client IP.
func RateLimiter(cfg RateLimiterConfig) gin.HandlerFunc {
	limiter := newRateLimiter(cfg)

	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		if !limiter.allow(clientIP) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
				"code":  "rate_limited",
			})
			return
		}

		c.Next()
	}
}
