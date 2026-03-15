package main

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"net/http"
	"strconv"

	"github.com/hashira21/snippetbox/internal/models"
)

// snippetCreateForm для представления данных формы и ошибок проверки полей формы.
// Все поля struct намеренно экспортированы (т.е. начинаются с заглавной буквы).
// Это связано с тем, что структурные поля должны быть экспортированы для того,
// чтобы быть прочитанными пакетом html/template при рендеринге шаблона.
type snippetCreateForm struct {
	Title       string
	Content     string
	Expires     int
	FieldErrors map[string]string
}

// Объявляем домашний обработчик, который печатает слайс байтов, содержащий
// "Hello from Snippetbox" в качестве тела ответа
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Call the newTemplateData() helper to get a templateData struct containing
	// the 'default' data (which for now is just the current year), and add the
	// snippets slice to it.
	data := app.newTemplateData(r)
	data.Snippets = snippets

	// Use the new render helper.
	app.render(w, r, http.StatusOK, "home.html", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	// And do the same thing again here...
	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, r, http.StatusOK, "view.tmpl.html", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.Form = snippetCreateForm{
		Expires: 365,
	}

	app.render(w, r, http.StatusOK, "create.tmpl.html", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// Сначала вызываем r.ParseForm() который добавляет данные из POST запроса
	// в r.PostForm map. Так же применимо к PUT и PATCH запросам.
	// в случае ошибки используем app.ClientError() хелпер для отправки ответа 400 Bad Request пользователю
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Используем метод r.PostForm.Get() для получения title, content из r.PostForm map
	// Метод r.PostForm.Get() всегда возращает данные формы как строку
	// Однако мы ожидаем значение expires как число
	// Поэтому нужно конвентировать используя strconv.Atoi()
	// В случае ошибки конвертации отправить ошибку 400 Bad Request
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := snippetCreateForm{
		Title:       r.PostForm.Get("title"),
		Content:     r.PostForm.Get("content"),
		Expires:     expires,
		FieldErrors: map[string]string{}, // мапа для хранения ошибок
	}

	// Проверка правильности введённых данных

	// title не пустой и не больше 100 символов
	if strings.TrimSpace(form.Title) == "" {
		form.FieldErrors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(form.Title) > 100 {
		form.FieldErrors["title"] = "This field cannot be more than 100 characters long"
	}

	// content не пустой
	if strings.TrimSpace(form.Content) == "" {
		form.FieldErrors["content"] = "This field cannot be blank"
	}

	// Значение expires 1, 7 или 365
	if form.Expires != 1 && form.Expires != 7 && form.Expires != 365 {
		form.FieldErrors["expires"] = "This field must equal 1, 7 or 365"
	}

	// Если были ошибки валидации, то в шаблон create.tmpl.html
	// передаются данные из экземпляра snippetCreateForm в качестве динамических данных
	// в поле Form. Используется HTTP код 422 Unprocessable Entity, когда отправляется
	// ответ, указывающий на ошибку при проверке.
	if len(form.FieldErrors) > 0 {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "create.tmpl.html", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
