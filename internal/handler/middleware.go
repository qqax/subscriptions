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
	return loggingMiddleware(
		recoveryMiddleware(
			requestIDMiddleware(
				corsMiddleware(
					authMiddleware(
						rateLimitMiddleware(
							handler,
						),
					),
				),
			),
		),
	)
}

// loggingMiddleware логирует все HTTP запросы
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Создаем ResponseWriter для захвата статуса
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Выполняем запрос
		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)

		// Логируем запрос
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

// recoveryMiddleware обрабатывает паники и логирует ошибки
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

				// Отправляем JSON ошибку вместо plain text
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

// requestIDMiddleware добавляет Request ID к каждому запросу
func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			// Генерируем UUID если не предоставлен
			requestID = generateRequestID()
		}

		// Добавляем в заголовки ответа
		w.Header().Set("X-Request-ID", requestID)

		// Добавляем в контекст запроса
		ctx := withRequestID(r.Context(), requestID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// corsMiddleware добавляет CORS headers
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

// authMiddleware проверяет аутентификацию
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
	// Реализация проверки токена
	return strings.HasPrefix(token, "Bearer ")
}

// rateLimitMiddleware ограничивает частоту запросов
func rateLimitMiddleware(next http.Handler) http.Handler {
	limiter := rate.NewLimiter(rate.Every(time.Second), 10) // 10 запросов в секунду

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

// responseWriter оборачивает http.ResponseWriter для захвата статуса
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

// Вспомогательные функции

func getRequestID(r *http.Request) string {
	if requestID := r.Header.Get("X-Request-ID"); requestID != "" {
		return requestID
	}
	return "unknown"
}

func generateRequestID() string {
	return "req-" + time.Now().Format(uuid.New().String())
}

// contextKey для хранения request ID в контексте
type contextKey string

const requestIDKey contextKey = "request_id"

func withRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

func getRequestIDFromContext(ctx context.Context) string {
	if requestID, ok := ctx.Value(requestIDKey).(string); ok {
		return requestID
	}
	return "unknown"
}

//// IsHeaderSent проверяет были ли уже отправлены заголовки
//func IsHeaderSent(w http.ResponseWriter) bool {
//	// Для стандартного http.ResponseWriter всегда возвращаем false
//	// В реальном приложении можно использовать type assertion для проверки
//	return false
//}
