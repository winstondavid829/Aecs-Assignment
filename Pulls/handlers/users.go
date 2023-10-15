package handlers

import (
	"time"

	"github.com/sony/gobreaker"
)

var breaker *gobreaker.CircuitBreaker

func init() {
	// Create a Circuit Breaker with your desired settings
	breaker = gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "GitHubAPI",
		MaxRequests: 3,               // Set an appropriate threshold
		Interval:    5 * time.Second, // Set an appropriate reset interval
	})
}
