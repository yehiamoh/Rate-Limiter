package tokenbucket

import (
	"sync"
	"time"
)

type TokenBucket struct {
	capacity    int
	tokens      int
	refilleRate time.Duration
	lastRefill  time.Time
	mu          sync.Mutex
}

func NewTokenBucket(capacity int, refillRate time.Duration) *TokenBucket {
	return &TokenBucket{
		capacity:    capacity,
		tokens:      capacity,
		lastRefill:  time.Now(),
		refilleRate: refillRate,
	}
}

func (tokenBucket *TokenBucket) IsAllow() bool {
	// Lock the mutex to ensure thread-safe access to the token bucket.
	tokenBucket.mu.Lock()
	defer tokenBucket.mu.Unlock() // Ensure the mutex is unlocked when the function exits.

	// Get the current time.
	now := time.Now()
	// Calculate the time elapsed since the last refill.
	elapsed := now.Sub(tokenBucket.lastRefill)

	// Add tokens to the bucket based on the elapsed time and refill rate.
	tokenBucket.tokens += int(elapsed / tokenBucket.refilleRate)

	// Ensure the number of tokens does not exceed the bucket's capacity.
	if tokenBucket.tokens > tokenBucket.capacity {
		tokenBucket.tokens = tokenBucket.capacity
	}
	// Update the last refill time to the current time.
	tokenBucket.lastRefill = now

	// Check if there are tokens available in the bucket.
	if tokenBucket.tokens > 0 {
		// Consume one token and allow the request.
		tokenBucket.tokens--
		return true
	}
	// If no tokens are available, deny the request.
	return false
}
