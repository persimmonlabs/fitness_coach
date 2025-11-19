package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimitConfig holds rate limiter configuration
type RateLimitConfig struct {
	RequestsPerMinute int
	CleanupInterval   time.Duration
}

// DefaultRateLimitConfig returns default rate limit configuration
func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		RequestsPerMinute: 100,
		CleanupInterval:   5 * time.Minute,
	}
}

// rateLimitEntry holds request count and timestamp for a user
type rateLimitEntry struct {
	count     int
	resetTime time.Time
}

// inMemoryRateLimiter implements in-memory rate limiting
type inMemoryRateLimiter struct {
	mu      sync.RWMutex
	entries map[string]*rateLimitEntry
	config  RateLimitConfig
}

// newInMemoryRateLimiter creates a new in-memory rate limiter
func newInMemoryRateLimiter(config RateLimitConfig) *inMemoryRateLimiter {
	limiter := &inMemoryRateLimiter{
		entries: make(map[string]*rateLimitEntry),
		config:  config,
	}

	// Start cleanup goroutine
	go limiter.cleanup()

	return limiter
}

// cleanup periodically removes expired entries
func (l *inMemoryRateLimiter) cleanup() {
	ticker := time.NewTicker(l.config.CleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		l.mu.Lock()
		now := time.Now()
		for key, entry := range l.entries {
			if now.After(entry.resetTime) {
				delete(l.entries, key)
			}
		}
		l.mu.Unlock()
	}
}

// allow checks if a request is allowed for the given key
func (l *inMemoryRateLimiter) allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	entry, exists := l.entries[key]

	if !exists || now.After(entry.resetTime) {
		// Create new entry
		l.entries[key] = &rateLimitEntry{
			count:     1,
			resetTime: now.Add(time.Minute),
		}
		return true
	}

	// Check if limit exceeded
	if entry.count >= l.config.RequestsPerMinute {
		return false
	}

	// Increment count
	entry.count++
	return true
}

// RateLimiter creates a middleware that limits requests per user
// Limits to 100 requests per minute per user by default
func RateLimiter(config RateLimitConfig) gin.HandlerFunc {
	limiter := newInMemoryRateLimiter(config)

	return func(c *gin.Context) {
		// Try to get userID from context (set by AuthJWT middleware)
		userID, exists := c.Get("userID")

		// Use IP address as fallback if user is not authenticated
		key := c.ClientIP()
		if exists && userID != nil {
			key = userID.(string)
		}

		// Check rate limit
		if !limiter.allow(key) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Maximum 100 requests per minute allowed.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
