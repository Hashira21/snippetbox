package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/justinas/nosurf"
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

func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Если пользователь не аутентифицирован, редирект на страницу входа
		// И выйти из цепочки middleware
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		// В противном случае установите заголовок "Cache-Control: no-store",
		// чтобы страницы требующие аутентификации, не сохранялись в кэше браузера пользователя
		// или в другом промежуточном кэше.
		w.Header().Add("Cache-Control", "no-store")

		// Вызвать следующий обработчик
		next.ServeHTTP(w, r)
	})
}

// Функция preventCSRF, которая использует специальный файл cookie CSRF с заданными атрибутами Secure, Path и HttpOnly.
func preventCSRF(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})

	return csrfHandler
}
