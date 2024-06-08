package middlewares

import (
	"log"
	"net/http"

	"github.com/CSPF-Founder/shieldsup/scanner-api/helpers"
)

// Use allows us to stack middleware to process the request
func Use(handler http.HandlerFunc, mid ...func(http.Handler) http.HandlerFunc) http.HandlerFunc {
	for _, m := range mid {
		handler = m(handler)
	}
	return handler
}

// ApplySecurityHeaders applies various security headers according to best-
// practices.
func ApplySecurityHeaders(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		csp := "frame-ancestors 'none';"
		w.Header().Set("Content-Security-Policy", csp)
		w.Header().Set("X-Frame-Options", "DENY")
		next.ServeHTTP(w, r)
	}
}

func RequireStaticAPIKey(serverAPIKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			inputAPIKey, err := helpers.GetBearerValue(r.Header)
			if err != nil {
				helpers.SendError(w, http.StatusUnauthorized, "Invalid API Key")
				return
			}

			if inputAPIKey != serverAPIKey {
				log.Printf("Unable to get the user by API Key: %v", inputAPIKey)
				helpers.SendError(w, http.StatusUnauthorized, "Invalid API Key")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
