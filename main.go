package main

import (
	"log"
	"net/http"
)

// Объявляем домашний обработчик, который печатает слайс байтов, содержащий
// "Hello from Snippetbox" в качестве тела ответа
func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from Snippetbox"))
}

func main() {
	//Используем функцию http.NewServeMux() для инициализации нового роутера
	//затем регистрируем функцию home как обработчик для URL пути "/"
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

	//Напечатать log сообщение о старте работы сервера
	log.Print("starting server on :4000")

	//Для запуска нового веб сервера используется функция http.ListenAndServe().
	//На вход подается два параметра: TCP адрес для прослущивания (:4000)
	//и роутер. Если функция вернёт ошибку выведется лог с помощью log.Fatal()
	//Каждая ошибка не nil
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
