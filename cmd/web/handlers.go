package main

import (
	"errors"
	"fmt"

	"net/http"
	"strconv"

	"github.com/hashira21/snippetbox/internal/models"
	"github.com/hashira21/snippetbox/internal/validator"
)

// snippetCreateForm для представления данных формы и ошибок проверки полей формы.
// Все поля struct намеренно экспортированы (т.е. начинаются с заглавной буквы).
// Это связано с тем, что структурные поля должны быть экспортированы для того,
// чтобы быть прочитанными пакетом html/template при рендеринге шаблона.
type snippetCreateForm struct {
	Title   string
	Content string
	Expires int
	validator.Validator
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
		Title:   r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
		Expires: expires,
	}

	// Проверка правильности введённых данных

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	// Use the Valid() method to see if any of the checks failed. If they did,
	// then re-render the template passing in the form in the same way as
	// before.
	// Если были ошибки валидации, то в шаблон create.tmpl.html
	// передаются данные из экземпляра snippetCreateForm в качестве динамических данных
	// в поле Form. Используется HTTP код 422 Unprocessable Entity, когда отправляется
	// ответ, указывающий на ошибку при проверке.
	if !form.Valid() {
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
