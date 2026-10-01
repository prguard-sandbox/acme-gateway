// Package ratelimit limits each client to a steady request rate with short bursts.
package ratelimit

import (
	"net"
	"net/http"
	"strings"
	"time"
)

type bucket struct {
	tokens float64
	last   time.Time
}

// Limiter allows `rate` requests per second per client, with bursts of up to `burst`.
type Limiter struct {
	rate    float64
	burst   float64
	clients map[string]*bucket
}

func New(rate, burst float64) *Limiter {
	return &Limiter{rate: rate, burst: burst, clients: map[string]*bucket{}}
}

// clientID is the caller's address: the first X-Forwarded-For hop when a load balancer
// put one there, the connection's address otherwise.
func clientID(r *http.Request) string {
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		return strings.TrimSpace(strings.Split(fwd, ",")[0])
	}
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	return host
}

func (l *Limiter) allow(id string) bool {
	now := time.Now()
	b, ok := l.clients[id]
	if !ok {
		b = &bucket{tokens: l.burst, last: now}
		l.clients[id] = b
	}
	b.tokens += now.Sub(b.last).Seconds() * l.rate
	if b.tokens > l.burst {
		b.tokens = l.burst
	}
	b.last = now
	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}

// Middleware answers 429 to clients over their limit and passes everyone else through.
func (l *Limiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !l.allow(clientID(r)) {
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
