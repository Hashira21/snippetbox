package main

import (
	"fmt"
	"log/slog"
	"net/http"
)

func commonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		w.Header().Set("Server", "Go")

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
		)

		app.logger.Info(
			"received request",
			slog.Any("ip", ip),
			slog.Any("proto", proto),
			slog.Any("method", method),
			slog.Any("uri", uri),
		)
		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			// Use the built-in recover() function to check if a panic occurred.
			// If a panic did happen, recover() will return the panic value.
			// If a panic didn't happen, it will return nil.
			pv := recover()

			// Если паника случилась
			if pv != nil {
				// Установим "Connection: close" заголовок в ответе
				w.Header().Set("Connection", "close")
				// Вызвать вспомогательный метод app.serverError для возврата ошибки 500
				app.serverError(w, r, fmt.Errorf("%v", pv))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
