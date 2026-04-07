package middleware

import (
	"net/http"

	coreHttp "github.com/fedotovmax/medods-test/internal/core/transport/http"
	"github.com/google/uuid"
)

func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(coreHttp.HeaderRequestID)

			if requestID == "" {
				requestID = uuid.NewString()
			}

			r.Header.Set(coreHttp.HeaderRequestID, requestID)
			w.Header().Set(coreHttp.HeaderRequestID, requestID)

			next.ServeHTTP(w, r)

		})
	}
}
