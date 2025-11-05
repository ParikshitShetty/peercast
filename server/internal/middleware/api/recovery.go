package middleware

import (
	"encoding/json"
	"log"
	"net/http"
)

// RecoverPanic catches panics and returns a 500 JSON error response.
func RecoverPanic(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("⚠️ Panic recovered: %v", rec)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"success": false,
					"error":   "Internal Server Error",
				})
			}
		}()
		next(w, r)
	}
}
