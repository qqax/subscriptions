package handler

import (
	"net/http"
	"subscription/internal/logger"
)

// HeaderSentChecker interface for checking sending headers
type HeaderSentChecker interface {
	HeaderSent() bool
}

func (rw *responseWriter) HeaderSent() bool {
	return rw.headerSent
}

// IsHeaderSent - universal check
func IsHeaderSent(w http.ResponseWriter) bool {
	if checker, ok := w.(HeaderSentChecker); ok {
		return checker.HeaderSent()
	}

	if rw, ok := w.(*responseWriter); ok {
		return rw.headerSent
	}

	return isHeaderSentConservative(w)
}

// isHeaderSentConservative - conservative check
func isHeaderSentConservative(w http.ResponseWriter) bool {
	headers := w.Header()

	defer func() {
		if recover() != nil {
			logger.Info().Msg("headers already sent")
		}
	}()

	headers.Add("X-Test-Probe", "value")

	return false
}
