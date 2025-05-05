// Этот пакет используется для хэширования пароля с помощью bcrypt
package utils

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword превращает обычную строку в захэшированную версию
func HashPassword(password string) (string, error) {
	//Функция GenerateFromPassword превращает строку в массив байтов,
	//и используя cost factor = 12 хэширует пароль (это нормальный фактор, пароль будет хэшироваться примерно 100-300мсек)
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", errors.New("database/utils/func HashPassword(): failed to hash password")
	}
	return string(bytes), nil
}

// Сравнивает пароль с его хэшем
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Проверяет соответствует ли пароль всем условиям
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	return nil
}
