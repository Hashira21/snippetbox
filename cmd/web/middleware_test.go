package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashira21/snippetbox/internal/assert"
)

func TestCommonHeader(t *testing.T) {
	rr := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Создаем имитацию HTTP-обработчика, которую мы можем передать в commonHeaders
	// промежуточное ПО, которое возвращает код состояния 200 и тело ответа «ОК».
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Передаем имитацию HTTP-обработчика в наше промежуточное ПО commonHeaders. Поскольку
	// commonHeaders *возвращает* http.Handler, мы можем вызвать его метод ServeHTTP(),
	// передав http.ResponseRecorder и фиктивный http.Request для его выполнения.
	commonHeaders(next).ServeHTTP(rr, req)

	// Вызовите метод Result() в http.ResponseRecorder, чтобы получить результаты теста
	res := rr.Result()
	defer res.Body.Close()

	// Проверки
	expectedValue := "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com"
	assert.Equal(t, res.Header.Get("Content-Security-Policy"), expectedValue)

	expectedValue = "origin-when-cross-origin"
	assert.Equal(t, res.Header.Get("Referrer-Policy"), expectedValue)

	expectedValue = "nosniff"
	assert.Equal(t, res.Header.Get("X-Content-Type-Options"), expectedValue)

	expectedValue = "deny"
	assert.Equal(t, res.Header.Get("X-Frame-Options"), expectedValue)

	expectedValue = "0"
	assert.Equal(t, res.Header.Get("X-XSS-Protection"), expectedValue)

	expectedValue = "Go"
	assert.Equal(t, res.Header.Get("Server"), expectedValue)

	assert.Equal(t, res.StatusCode, http.StatusOK)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)

	assert.Equal(t, string(body), "OK")
}
