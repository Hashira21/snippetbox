package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// Объявляем домашний обработчик, который печатает слайс байтов, содержащий
// "Hello from Snippetbox" в качестве тела ответа
func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from Snippetbox"))
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	msg := fmt.Sprintf("Display a specific snippet with ID %d...", id)
	w.Write([]byte(msg))
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new snippet..."))
}

func snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Save a new snippet..."))
}

func main() {
	//Используем функцию http.NewServeMux() для инициализации нового роутера
	//затем регистрируем функцию home как обработчик для URL пути "/"
	mux := http.NewServeMux()
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
