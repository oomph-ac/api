package utils

import "net/http"

// ClientIP returns the client IP forwarded by Cloudflare.
func ClientIP(r *http.Request) string {
	return r.Header.Get("CF-Connecting-IP")
}
