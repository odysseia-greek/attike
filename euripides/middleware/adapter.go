package middleware

import (
	"net/http"
)

type Adapter func(http.Handler) http.Handler

// Adapt Iterate over adapters and run them one by one
func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

func SetCorsHeaders() Adapter {
	return func(f http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//allow all CORS
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST,OPTIONS")
			if r.Method == "OPTIONS" {
				return
			}
			f.ServeHTTP(w, r)
		})
	}
}
