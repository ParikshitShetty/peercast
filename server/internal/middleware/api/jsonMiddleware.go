package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// responseCaptureWriter captures the response body and status code.
type responseCaptureWriter struct {
	http.ResponseWriter
	body       bytes.Buffer
	statusCode int
}

func (rw *responseCaptureWriter) WriteHeader(code int) {
	if rw.statusCode == 0 { // prevent duplicate calls
		rw.statusCode = code
	}
}

func (rw *responseCaptureWriter) Write(b []byte) (int, error) {
	return rw.body.Write(b)
}

// JSONResponseHandler wraps responses in a standard JSON format.
func JSONResponseHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rcw := &responseCaptureWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Execute the original handler (but capture its writes)
		next(rcw, r)

		// Always set JSON header
		w.Header().Set("Content-Type", "application/json")

		// Handle errors (status >= 400)
		if rcw.statusCode >= 400 {
			w.WriteHeader(rcw.statusCode)
			resp := map[string]interface{}{
				"success": false,
				"error":   rcw.body.String(),
			}
			_ = json.NewEncoder(w).Encode(resp)
			return
		}

		// Try to decode valid JSON
		var data interface{}
		if err := json.Unmarshal(rcw.body.Bytes(), &data); err == nil {
			w.WriteHeader(rcw.statusCode)
			resp := map[string]interface{}{
				"success": true,
				"data":    data,
			}
			_ = json.NewEncoder(w).Encode(resp)
			return
		}

		// Fallback: non-JSON responses
		w.WriteHeader(rcw.statusCode)
		w.Write(rcw.body.Bytes())
	}
}
