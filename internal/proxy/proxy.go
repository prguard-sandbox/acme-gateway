package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// To forwards requests under prefix to the backend at base, without the prefix.
func To(base, prefix string) http.Handler {
	target, err := url.Parse(base)
	if err != nil {
		panic(err)
	}
	rp := httputil.NewSingleHostReverseProxy(target)
	return http.StripPrefix(strings.TrimSuffix(prefix, "/"), rp)
}
