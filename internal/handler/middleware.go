package handler

import (
	"context"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"time"

	"golang.org/x/time/rate"
	"subscription/internal/logger"
)

func AddMiddleware(handler http.Handler) http.Handler {
	return requestIDMiddleware(
		loggingMiddleware(
			recoveryMiddleware(
				requestIDMiddleware(
					corsMiddleware(
						//authMiddleware(
						rateLimitMiddleware(
							handler,
						),
						//),
					),
				),
			),
		),
	)
}

// loggingMiddleware logs all HTTP requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)

		logger.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("query", r.URL.RawQuery).
			Str("ip", r.RemoteAddr).
			Str("user_agent", r.UserAgent()).
			Int("status", wrapped.statusCode).
			Dur("duration", duration).
			Int64("duration_ms", duration.Milliseconds()).
			Str("request_id", getRequestID(r)).
			Msg("HTTP request")
	})
}

// recoveryMiddleware handles panics and logs errors
func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error().
					Interface("panic", err).
					Str("method", r.Method).
					Str("path", r.URL.Path).
					Str("request_id", getRequestID(r)).
					Msg("Recovered from panic")

				// Проверяем не был ли уже отправлен заголовок
				if IsHeaderSent(w) {
					logger.Warn().
						Str("request_id", getRequestID(r)).
						Msg("Headers already sent, cannot send error response")
					return
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)

				errorResponse := []byte(`{"error":"internal_error","message":"Internal Server Error"}`)
				if _, writeErr := w.Write(errorResponse); writeErr != nil {
					logger.Error().
						Err(writeErr).
						Str("request_id", getRequestID(r)).
						Msg("Failed to write JSON error response")
				}
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// requestIDMiddleware adds Request ID to each request
func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		w.Header().Set("X-Request-ID", requestID)

		ctx := withRequestID(r.Context(), requestID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// corsMiddleware adds CORS headers
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// authMiddleware checks authentication
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			logger.Warn().
				Str("path", r.URL.Path).
				Str("request_id", getRequestID(r)).
				Msg("Unauthorized request")

			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Проверка токена (упрощенно)
		if !isValidToken(token) {
			logger.Warn().
				Str("path", r.URL.Path).
				Str("request_id", getRequestID(r)).
				Msg("Invalid token")

			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func isValidToken(token string) bool {
	return strings.HasPrefix(token, "Bearer ")
}

// rateLimitMiddleware limits the frequency of requests
func rateLimitMiddleware(next http.Handler) http.Handler {
	limiter := rate.NewLimiter(rate.Every(time.Second), 10) // 10 requests per second

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			logger.Warn().
				Str("ip", r.RemoteAddr).
				Str("path", r.URL.Path).
				Str("request_id", getRequestID(r)).
				Msg("Rate limit exceeded")

			w.WriteHeader(http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// responseWriter wraps http.ResponseWriter to capture status
type responseWriter struct {
	http.ResponseWriter
	statusCode  int
	headerSent  bool
	writeCalled bool
}

func (rw *responseWriter) WriteHeader(code int) {
	if !rw.headerSent {
		rw.statusCode = code
		rw.ResponseWriter.WriteHeader(code)
		rw.headerSent = true
	}
}

func (rw *responseWriter) Write(data []byte) (int, error) {
	if !rw.headerSent {
		rw.WriteHeader(http.StatusOK)
	}
	rw.writeCalled = true
	return rw.ResponseWriter.Write(data)
}

// Helper functions

func getRequestID(r *http.Request) string {
	if requestID := r.Header.Get("X-Request-ID"); requestID != "" {
		return requestID
	}
	return "unknown"
}

func generateRequestID() string {
	return "req-" + time.Now().Format(uuid.New().String())
}

func withRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, "request_id", requestID)
}
