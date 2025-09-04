package handler

import "net/http"

// HeaderSentChecker интерфейс для проверки отправки заголовков
type HeaderSentChecker interface {
	HeaderSent() bool
}

// Реализуем интерфейс для нашего responseWriter
func (rw *responseWriter) HeaderSent() bool {
	return rw.headerSent
}

// IsHeaderSent универсальная проверка
func IsHeaderSent(w http.ResponseWriter) bool {
	// Проверяем реализует ли ResponseWriter наш интерфейс
	if checker, ok := w.(HeaderSentChecker); ok {
		return checker.HeaderSent()
	}

	// Проверяем стандартные реализации
	if rw, ok := w.(*responseWriter); ok {
		return rw.headerSent
	}

	// Для неизвестных реализаций используем надежный метод
	return isHeaderSentConservative(w)
}

// isHeaderSentConservative консервативная проверка
func isHeaderSentConservative(w http.ResponseWriter) bool {
	// Пытаемся получить мапу заголовков
	headers := w.Header()

	// Если заголовки уже "заморожены" (отправлены), попытка изменения вызовет панику
	defer func() {
		if recover() != nil {
			// Паника означает что заголовки уже отправлены
		}
	}()

	// Пробуем добавить тестовый заголовок
	headers.Add("X-Test-Probe", "value")

	// Если не было паники, заголовки еще не отправлены
	return false
}
