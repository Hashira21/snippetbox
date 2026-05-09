package main

import (
	"testing"
	"time"

	"github.com/hashira21/snippetbox/internal/assert"
)

func TestHumanDate(t *testing.T) {
	// Создаем слайс анонимных структур, содержащих название тестового сценария,
	// входные данные функции humanDate() (поле tm) и ожидаемый результат (поле want)
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2024, 3, 17, 10, 15, 0, 0, time.UTC),
			want: "17 Mar 2024 at 10:15",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2024, 3, 17, 10, 15, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "17 Mar 2024 at 09:15",
		},
	}

	// Перебираем тестовые сценарии.
	for _, tt := range tests {
		// Используем функцию t.Run() для запуска подтеста каждого тестового сценария.
		// Первый параметр - название теста, второй - анонимная функция, содержащий фактический
		// тест для каждого случая
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)
			assert.Equal(t, hd, tt.want)
		})
	}

}
