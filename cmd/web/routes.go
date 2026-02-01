package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// Используем mux.Handle() для регистрации файлового сервера в качестве обработчика
	// для всех URL путей, которые начинаются с /static/. Для сопоставления путей мы
	// удаляем /static префикс перед тем, как запрос достигнет файлового сервера
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippet/create", app.snippetCreate)
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

	return mux
}
