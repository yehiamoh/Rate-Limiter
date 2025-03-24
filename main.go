package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	tokenbucket "github.com/yehiamoh/rate-limiter/token-bucket"
)

type Message struct {
	Status string `json:"status"`
	Body   string `json:"body"`
}

func newMessage(status string, resBody string) *Message {
	return &Message{Status: status, Body: resBody}
}

var tokenBucketlimiter = tokenbucket.NewTokenBucket(5, 1*time.Second)

func endpointHandler(w http.ResponseWriter, r *http.Request) {
	if !tokenBucketlimiter.IsAllow() {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	message := newMessage("successful", "You have reaced the api successfully")
	if err := json.NewEncoder(w).Encode(message); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func main() {
	http.HandleFunc("/ping", endpointHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalln("Error in initalizing the server on port 8080")
	}
}
