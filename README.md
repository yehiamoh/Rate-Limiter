# Go Rate Limiter

[![Go Reference](https://pkg.go.dev/badge/github.com/yehiamoh/rate-limiter.svg)](https://pkg.go.dev/github.com/yehiamoh/rate-limiter)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/yehiamoh/rate-limiter)](https://golang.org/dl/)

A flexible and efficient rate limiter implementation in Go, supporting **Token Bucket** and **Per-Client** strategies.

---

## Features

- 🪣 Token Bucket Algorithm
- 🖥️ Per-Client Rate Limiting
- 🔒 Thread-Safe Implementation
- 📦 Zero Dependencies

---

## Installation

To install the package, use the following command:

```bash
go get github.com/yehiamoh/rate-limiter
```

---

## Usage

### 1. Token Bucket Strategy

The **Token Bucket Strategy** applies a global rate limit for all requests.

```go
package main

import (
    "net/http"
    "time"
    "github.com/yehiamoh/rate-limiter/tokenbucket"
)

func main() {
    // Create a token bucket with a capacity of 10 tokens and a refill rate of 1 token per second.
    limiter, err := tokenbucket.NewTokenBucket(10, 1*time.Second)
    if err != nil {
        panic(err)
    }

    http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
        if !limiter.IsAllowed() {
            http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
            return
        }
        // Handler logic
        w.Write([]byte("Request successful"))
    })

    http.ListenAndServe(":8080", nil)
}
```

---

### 2. Per-Client Strategy

The **Per-Client Strategy** applies rate limits individually for each client, identified by their unique client ID (e.g., IP address).

```go
package main

import (
    "net/http"
    "time"
    "github.com/yehiamoh/rate-limiter/perclient"
)

func main() {
    // Create a per-client rate limiter with a capacity of 5 tokens and a refill rate of 1 token every 2 seconds.
    limiter := perclient.NewPerClientLimiter(5, 2*time.Second)

    http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
        clientID := r.RemoteAddr // Use the client's IP address as the unique identifier.
        if !limiter.IsAllowed(clientID) {
            http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
            return
        }
        // Handler logic
        w.Write([]byte("Request successful"))
    })

    http.ListenAndServe(":8080", nil)
}
```

---

## Testing

To test the rate limiter, follow these steps:

1. Run the demo server:

   ```bash
   go run main.go
   ```

2. Test the rate limits using `curl`:

   ```bash
   # Send 6 requests to the /api endpoint
   for i in {1..6}; do curl http://localhost:8080/api; done
   ```

   - The first few requests will succeed.
   - Once the rate limit is exceeded, the server will respond with `429 Too Many Requests`.

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

## Author

Created and maintained by [Yehiamoh](https://github.com/yehiamoh).
