package models

import (
	"errors"
)

var (
	ErrNoRecord = errors.New("models: no matching record found")

	// Если пользователь ввёл неверный email или пароль
	ErrInvalidCredentials = errors.New("models: invalid credantials")

	// Если пользователь пытается зарегистрироваться на использованную почту
	ErrDuplicateEmail = errors.New("models: duplicate email")
)
