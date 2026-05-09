package main

import (
	"net/http"
	"testing"

	"github.com/hashira21/snippetbox/internal/assert"
)

/*
func TestPing(t *testing.T) {
	// Инициализируем httptest.ResponseRecoder
	rr := httptest.NewRecorder()

	// Инициализируем новый фиктивный http.Request
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Вызываем функцию ping, передавая httptest.ResponseRecorder и http.Request
	ping(rr, req)

	// Вызываем метод Result() для httptest.ResponseRecorder для получения
	// http.Response сгенерированным обработчиком ping
	res := rr.Result()
	defer res.Body.Close()

	// Проверка, что код состояния 200
	assert.Equal(t, res.StatusCode, http.StatusOK)

	// Проверка, что тело ответа = "OK"
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)

	assert.Equal(t, string(body), "OK")
}
*/

func TestPing(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	res := ts.get(t, "/ping")
	assert.Equal(t, res.status, http.StatusOK)
	assert.Equal(t, res.body, "OK")
}
