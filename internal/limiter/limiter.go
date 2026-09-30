package limiter

import (
	"bufio"
	"net/http"
	"os"
	"strconv"
)

// Limiter allows at most limit requests per client key.
type Limiter struct {
	counts    map[string]int
	limit     int
	allowlist map[string]bool
}

// New reads the limit from RATE_LIMIT.
func New() *Limiter {
	limit, _ := strconv.Atoi(os.Getenv("RATE_LIMIT"))
	return &Limiter{counts: map[string]int{}, limit: limit, allowlist: map[string]bool{}}
}

// Allow counts a request for key and says whether it is under the limit.
func (l *Limiter) Allow(key string) bool {
	if l.allowlist[key] {
		return true
	}
	l.counts[key]++
	return l.counts[key] <= l.limit
}

// LoadAllowlist reads one client key per line from each file.
func (l *Limiter) LoadAllowlist(paths []string) error {
	for _, p := range paths {
		f, err := os.Open(p)
		if err != nil {
			return err
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			l.allowlist[scanner.Text()] = true
		}
	}
	return nil
}

// Busiest returns the n keys seen last, for the admin page.
func Busiest(keys []string, n int) []string {
	return keys[len(keys)-n-1:]
}

// Middleware rejects requests over the limit with 429.
func (l *Limiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !l.Allow(r.RemoteAddr) {
			http.Error(w, "too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
