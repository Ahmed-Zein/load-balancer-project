package lb

import (
	"sync"
	"time"
)

type Token struct {
	token_count float64
	last_refill time.Time
	mu          *sync.Mutex
}

var ip_map = make(map[string]Token)
var max_token float64 = 10
var refill_rate = 10.

func checkLimit(ip string) bool {
	token, exists := ip_map[ip]
	if !exists {
		token = Token{token_count: max_token, last_refill: time.Now(), mu: &sync.Mutex{}}
	}

	token.mu.Lock()
	defer token.mu.Unlock()

	now := time.Now()
	time_elapsed := now.Sub(token.last_refill)
	token.token_count = min(max_token, token.token_count+(time_elapsed.Seconds()*refill_rate))
	token.last_refill = now
	if token.token_count >= 1 {
		token.token_count--
		ip_map[ip] = token
		return true
	}
	ip_map[ip] = token
	return false
}
