package http

import (
	"github.com/go-chi/chi/middleware"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
	"time"
	"todo/internal/trace"
	"todo/user"
)

const traceIDHeader = "X-Trace-ID"

func profilingMiddleware(log *zap.SugaredLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			defer func() {
				userID := uint(0)
				usr, ok := r.Context().Value(user.UserContextKey).(user.User)
				if ok {
					userID = usr.ID
				}
				start := time.Now()
				log.With("status_code", rw.Status()).
					With("http_verb", r.Method).
					With("bytes", rw.BytesWritten()).
					With("latency", time.Since(start).Seconds()).
					With("uri", r.URL.String()).
					With("trace_id", r.Header.Get(traceIDHeader)).
					With("userID", userID).
					Info("router: http request")
			}()
			next.ServeHTTP(rw, r)
		})
	}
}

func traceMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			traceID := r.Header.Get(traceIDHeader)
			if traceID == "" {
				traceID = uuid.New().String()
			}
			r.WithContext(trace.WithValue(r.Context(), traceID))
			w.Header().Set(traceIDHeader, traceID)
			next.ServeHTTP(w, r)
		})
	}
}
