package tokenbucket

import (
	"fmt"
	"sync"
	"time"
)

type TokenBucket struct {
	capacity   int
	tokens     int
	refillRate time.Duration
	lastRefill time.Time
	mu         sync.Mutex
}

func NewTokenBucket(capacity int, refillRate time.Duration) (*TokenBucket, error) {
	if capacity <= 0 {
		return nil, fmt.Errorf("capacity must be positive, got %d", capacity)
	}
	if refillRate <= 0 {
		return nil, fmt.Errorf("refill rate must be positive, got %v", refillRate)
	}

	return &TokenBucket{
		capacity:   capacity,
		tokens:     capacity,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}, nil
}

// IsAllowed checks if a request should be allowed based on available tokens
func (tokenBucket *TokenBucket) IsAllowed() bool {
	tokenBucket.mu.Lock()
	defer tokenBucket.mu.Unlock()

	now := time.Now()
	tokenBucket.refillTokens(now)

	if tokenBucket.tokens > 0 {
		tokenBucket.tokens--
		return true
	}
	return false
}

// refillTokens adds tokens to the bucket based on elapsed time
func (tokenBucket *TokenBucket) refillTokens(now time.Time) {
	elapsed := now.Sub(tokenBucket.lastRefill)

	// Safeguard: Ensure tokensToAdd does not exceed the bucket's capacity.
	tokensToAdd := min(int(elapsed.Nanoseconds()/tokenBucket.refillRate.Nanoseconds()), tokenBucket.capacity)

	if tokensToAdd > 0 {
		tokenBucket.tokens += tokensToAdd
		// Only update lastRefill when we actually add tokens to the bucket.
		tokenBucket.lastRefill = now
	}
}

// GetAvailableTokens returns the current number of available tokens
func (tokenBucket *TokenBucket) GetAvailableTokens() int {
	tokenBucket.mu.Lock()
	defer tokenBucket.mu.Unlock()

	tokenBucket.refillTokens(time.Now())
	return tokenBucket.tokens
}
