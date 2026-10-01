package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/prguard-sandbox/acme-gateway/internal/proxy"
	"github.com/prguard-sandbox/acme-gateway/internal/ratelimit"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	mux := http.NewServeMux()
	mux.Handle("/store/", proxy.To(os.Getenv("STORE_URL"), "/store"))
	mux.Handle("/billing/", proxy.To(os.Getenv("BILLING_URL"), "/billing"))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })

	// Requests per second each client may make, with bursts of 20.
	rps, _ := strconv.ParseFloat(os.Getenv("RATE_LIMIT_RPS"), 64)
	limiter := ratelimit.New(rps, 20)

	log.Printf("listening on :%s (%.0f req/s per client)", port, rps)
	log.Fatal(http.ListenAndServe(":"+port, limiter.Middleware(mux)))
}
