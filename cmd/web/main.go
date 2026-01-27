package main

import (
	"log"
	"net/http"
)

func main() {
	//Используем функцию http.NewServeMux() для инициализации нового роутера
	//затем регистрируем функцию home как обработчик для URL пути "/"
	mux := http.NewServeMux()

	// Создаем файл сервер, который обслуживает файлы из каталога ./us/static
	// Отметим, что путь связан с местом расположения проекта
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// Используем mux.Handle() для регистрации файлового сервера в качестве обработчика
	// для всех URL путей, которые начинаются с /static/. Для сопоставления путей мы
	// удаляем /static префикс перед тем, как запрос достигнет файлового сервера
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	mux.HandleFunc("GET /snippet/create", snippetCreate)
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)

	//Напечатать log сообщение о старте работы сервера
	log.Print("starting server on :4000")

	//Для запуска нового веб сервера используется функция http.ListenAndServe().
	//На вход подается два параметра: TCP адрес для прослущивания (:4000)
	//и роутер. Если функция вернёт ошибку выведется лог с помощью log.Fatal()
	//Каждая ошибка не nil
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
