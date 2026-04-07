package middleware

import (
	"net/http"

	"github.com/fedotovmax/medods-test/internal/core/logger"
	coreHttp "github.com/fedotovmax/medods-test/internal/core/transport/http"
)

func Logger(log logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(coreHttp.HeaderRequestID)

			l := log.With(
				logger.String("request_id", requestID),
				logger.String("url", r.URL.String()),
			)

			ctx := logger.ToContext(r.Context(), l)

			next.ServeHTTP(w, r.WithContext(ctx))

		})
	}
}
