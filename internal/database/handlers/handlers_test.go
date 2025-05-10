package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/C1scoR/Go_Project_Yandex_Lyceum/internal/database/models"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	test_data := []struct {
		name   string
		user   any
		result any
	}{
		//Если вдруг решите запустить несколько раз тесты, то в первом поменяйте почту
		//Тесты интеграционные, а значит каждый пользователь из теста записывается в БД
		{
			name: "правильный ответ",
			user: models.UserRegister{
				Email:    "johndoe6@gmail.com",
				Password: "12345678",
			},
			result: "User registered successfully",
		},
		{
			name: "неверный формат email",
			user: models.UserRegister{
				Email:    "johnDoe@gmail.com",
				Password: "12345678",
			},
			result: "\"invalid email format\"",
		},
		{
			name: "Ошибка регистрации существующего пользователя",
			user: models.UserRegister{
				Email:    "johndoe@gmail.com",
				Password: "12345678",
			},
			result: "Пользователь уже существует 409",
		},
	}

	for _, tc := range test_data {
		t.Run(tc.name, func(t *testing.T) {

			body, err := json.Marshal(tc.user)
			if err != nil {
				t.Fatalf("Ошибка маршалинга: %v", err)
			}

			req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/register", bytes.NewBuffer(body))
			if err != nil {
				t.Fatalf("Ошибка создания запроса: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Ошибка выполнения запроса: %v", err)
			}
			defer resp.Body.Close()

			bodyResp, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Ошибка чтения тела ответа: %v", err)
			}

			var response map[string]interface{}
			err = json.Unmarshal(bodyResp, &response)
			if err != nil {
				// Если не удается распарсить ответ как JSON, проверяем как строку
				responseText := string(bodyResp)
				if responseText != tc.result {
					t.Errorf("Ожидался ответ %q, но получен %q", tc.result, responseText)
				}
			} else {
				// Если ответ успешно распарсился как JSON, проверяем его поле "message"
				if response["message"] != tc.result {
					t.Errorf("Ожидался ответ %q, но получен %v", tc.result, response["message"])
				}
			}
		})
	}
}

func TestLogin(t *testing.T) {
	test_data := []struct {
		name   string
		user   models.UserLogin
		result any
	}{
		{
			name: "правильный ответ",
			user: models.UserLogin{Email: "johndoe@gmail.com", Password: "12345678"},
			result: struct {
				Expires_in string `json:"Expires_in"`
				Token_type string `json:"Token_type"`
			}{
				Expires_in: "86400",
				Token_type: "Bearer",
			},
		},
		{
			name:   "пользователя не существует",
			user:   models.UserLogin{Email: "notevenexist@gmail.com", Password: "12345678"},
			result: fmt.Sprintf("Такого пользователя нет, пожалуйста зарегистрируйтесь: %d", http.StatusUnauthorized),
		},
		{
			name:   "неверный пароль",
			user:   models.UserLogin{Email: "johndoe@gmail.com", Password: "87654321"},
			result: fmt.Sprintf("Пароль неверный %d", http.StatusUnprocessableEntity),
		},
	}
	for _, tt := range test_data {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.user)
			if err != nil {
				t.Fatalf("Ошибка маршалинга: %v", err)
			}
			req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/login", bytes.NewBuffer(body))
			if err != nil {
				t.Fatalf("Не удалось создать запрос на login")
			}
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Не удалось выполнить запрос на логин")
			}
			defer resp.Body.Close()
			bodyResp, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Не удалось прочитать тело ответа")
			}
			var responseJSON map[string]interface{}
			err = json.Unmarshal(bodyResp, &responseJSON)
			if err != nil {
				//"Не удалось распарсить тело ответа как json"
				responseText := string(bodyResp)
				assert.Equal(t, responseText, tt.result.(string))
			} else {
				assert.Equal(t, float64(86400), responseJSON["Expires_in"])
				assert.Equal(t, "Bearer", responseJSON["Token_type"])
			}

		})
	}
}
