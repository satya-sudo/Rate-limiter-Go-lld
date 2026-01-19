package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/satya-sudo/Rate-limiter-Go-lld/internal/limiter"
	"github.com/satya-sudo/Rate-limiter-Go-lld/internal/middleware"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"Status": "up"})
}

func apiGatewayHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"Message": "You have reached the endpoint"})
}

func main() {
	rateLimiter := limiter.NewTokenBucketRateLimiter(5, 0.5)
	apiGatewayHandler := http.HandlerFunc(apiGatewayHandle)

	http.HandleFunc("/health", healthHandler)
	http.Handle("/api", middleware.RateLimiterMiddleware(rateLimiter, apiGatewayHandler))
	fmt.Println("Listening on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
	rateLimiter.StartCleanUp(1*time.Minute, 30*time.Minute)
	defer rateLimiter.Stop()
}
