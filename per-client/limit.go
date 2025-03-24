package perclient

import (
	"sync"
	"time"

	tokenbucket "github.com/yehiamoh/Rate-Limiter/token-bucket"
)

// PerClientRateLimiter manages rate limiting for multiple clients.
type PerClientRateLimiter struct {
	buckets    sync.Map      // Stores TokenBucket instances for each client.
	capacity   int           // Maximum tokens per bucket.
	refillRate time.Duration // Time to add one token.
}

// NewPerClientLimiter creates a new PerClientRateLimiter.
func NewPerClientLimiter(capacity int, refillRate time.Duration) *PerClientRateLimiter {
	return &PerClientRateLimiter{
		capacity:   capacity,
		refillRate: refillRate,
	}
}

// GetBuckets retrieves or creates a TokenBucket for a client.
func (perClientRateLimiter *PerClientRateLimiter) GetBuckets(clientID string) *tokenbucket.TokenBucket {
	bucket, ok := perClientRateLimiter.buckets.Load(clientID)
	if !ok {
		newBucket := tokenbucket.NewTokenBucket(perClientRateLimiter.capacity, perClientRateLimiter.refillRate)
		perClientRateLimiter.buckets.Store(clientID, newBucket)
		return newBucket
	}
	return bucket.(*tokenbucket.TokenBucket) // Type assertion to *TokenBucket.
}

// IsAllow checks if a request from a client is allowed.
func (perClientRateLimiter *PerClientRateLimiter) IsAllow(clientID string) bool {
	bucket := perClientRateLimiter.GetBuckets(clientID)
	return bucket.IsAllow()
}
