package main

import (
	"log"
	"net/http"
	"os"

	"github.com/prguard-sandbox/acme-gateway/internal/proxy"
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
	log.Printf("listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
