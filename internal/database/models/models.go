package models

import (
	"errors"
	"regexp"
)

// используется для того, чтобы добавить пользователя в БД
type User struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
}

// UserLogin нужен для входа пользователя в систему
type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserRegister нужен для регистрации пользователя в системе
type UserRegister struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Validate проверяет верен ли формат email.
// Если нет, то возвращает ошибку
func (u *UserRegister) Validate() error {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(u.Email) {
		return errors.New("invalid email format")
	}
	return nil
}
