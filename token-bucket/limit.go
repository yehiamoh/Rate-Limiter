package tokenbucket

import (
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

func NewTokenBucket(capacity int, refillRate time.Duration) *TokenBucket {
	return &TokenBucket{
		capacity:   capacity,
		tokens:     capacity,
		lastRefill: time.Now(),
		refillRate: refillRate,
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
	if elapsed > 0 {
		tokensToAdd := int(elapsed / tokenBucket.refillRate)
		// Safeguard: Ensure tokensToAdd does not exceed the bucket's capacity.
		if tokensToAdd > tokenBucket.capacity {
			tokensToAdd = tokenBucket.capacity
		}
		if tokensToAdd > 0 {
			tokenBucket.tokens += tokensToAdd
			// Only update lastRefill when we actually add tokens to the bucket.
			tokenBucket.lastRefill = now
		}
	}

	// Check if there are tokens available in the bucket.
	if tokenBucket.tokens > 0 {
		// Consume one token and allow the request.
		tokenBucket.tokens--
		return true
	}
	// If no tokens are available, deny the request.
	return false
}
