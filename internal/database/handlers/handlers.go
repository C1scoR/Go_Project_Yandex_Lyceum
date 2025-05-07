package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/C1scoR/Go_Project_Yandex_Lyceum/internal/database/config"
	"github.com/C1scoR/Go_Project_Yandex_Lyceum/internal/database/models"
	"github.com/C1scoR/Go_Project_Yandex_Lyceum/internal/database/utils"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	db        *config.Litedb
	jwtSecret []byte
	// Add token expiration configuration
	tokenExpiration time.Duration
}

func NewAuthHandler(db *config.Litedb, jwtSecret []byte) *AuthHandler {
	return &AuthHandler{
		db:              db,
		jwtSecret:       jwtSecret,
		tokenExpiration: 24 * time.Hour, // Default 24 hour expiration
	}
}

// Обрабатывает регистрацию пользователя
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user models.UserRegister

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Fprintf(w, "Ошибка при декодировании тела запроса: %d", http.StatusBadRequest)
		log.Println("handlers/func Register()/ ошибка при декодировании тела запроса", err)
		return
	}

	if err := user.Validate(); err != nil {
		fmt.Fprintf(w, "%q", err)
		log.Println("handlers/func Register()/ошибка валидации данных", err)
		return
	}

	var exists bool
	err = h.db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", user.Email).Scan(&exists)
	if err != nil {
		fmt.Fprintf(w, "%q", err)
		log.Println("handlers/func Register(): ", err)
		return
	}
	if exists {
		fmt.Fprintf(w, "Пользователь уже существует %d", http.StatusConflict)
		log.Println("handlers/func Register(): Пользователь уже существует, и пытается залогиниться заново")
		return
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		fmt.Fprintf(w, "Ошибка сервера: %d", http.StatusInternalServerError)
		log.Println("handlers/func Register(): ошибка хэширования пароля: ", err)
		return
	}
	tx, err := h.db.DB.Begin()
	if err != nil {
		fmt.Fprintf(w, "Транзакция не была начата: %d", http.StatusInternalServerError)
		log.Println("handlers/func Register(): транзакция не началась", err)
		return
	}
	var id int
	err = tx.QueryRow(`INSERT INTO users (email, password_hash) 
        VALUES (?, ?) 
        RETURNING id`, user.Email, hashedPassword).Scan(&id)
	if err != nil {
		tx.Rollback()
		fmt.Fprintf(w, "Транзакция не была завершена корректно из-за неправильных данных: %d", http.StatusBadRequest)
		log.Println("handlers/func Register(): транзакция не была заверешена: ", err)
		return
	}
	if err = tx.Commit(); err != nil {
		fmt.Fprintf(w, "Ошибка коммита транзакции: %d", http.StatusInternalServerError)
		log.Println("handlers/func Register(): что-то произошло при коммите: ", err)
		return
	}

	response := map[string]interface{}{
		"message": "User registered successfully",
		"user_id": id,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
	log.Println("пользователь был зарегистрирован:")
}

// Обрабатывает Аутентификацию и jwt генерацию
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var login models.UserLogin
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		fmt.Fprintf(w, "Ошибка при декодировании тела запроса: %d", http.StatusBadRequest)
		log.Println("handlers/func Login()/ ошибка при декодировании тела запроса", err)
		return
	}

	var user models.User
	err = h.db.DB.QueryRow(`
        SELECT id, email, password_hash 
        FROM users 
        WHERE email = ?`,
		login.Email,
	).Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err == sql.ErrNoRows {
		fmt.Fprintf(w, "Такого пользователя нет, пожалуйста зарегистрируйтесь: %d", http.StatusUnauthorized)
		log.Println("handlers/func Login()/QueryRow(): пользователь не зареган", err)
		return
	}

	if err != nil {
		fmt.Fprintf(w, "Ошибка входа в систему: %d", http.StatusInternalServerError)
		log.Println("handlers/func Login(): что-то пошло не так при запросе данных в БД: ", err)
		return
	}

	if !utils.CheckPasswordHash(login.Password, user.PasswordHash) {
		// Используем это для безопасности
		fmt.Fprintf(w, "Пароль неверный %d", http.StatusUnprocessableEntity)
		log.Println("handlers/func Login()/CheckPasswordHash(): что-то пошло не так при сопоставлении паролей", err)
		return
	}
	//Генерируем jwt с данными
	now := time.Now()
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"iat":     now.Unix(),
		"exp":     now.Add(h.tokenExpiration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(h.jwtSecret)
	if err != nil {
		fmt.Fprintf(w, "Не удалось создать jwt токен %d", http.StatusUnprocessableEntity)
		log.Println("handlers/func Login()/SignedString(): не удалось создать jwt токен: ", err)
		return
	}
	response := struct {
		Token      string
		Expires_in float64
		Token_type string
	}{
		Token:      tokenString,
		Expires_in: h.tokenExpiration.Seconds(),
		Token_type: "Bearer",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// RefreshToken генерирует новый токен для залогининых пользователей
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	//Получаем ID пользоватея через контекст (от middleware)
	//	userID, exists := r.Context().Value("user_id").(string)
	ID := r.Context().Value("user_id").(float64)
	userID := int(ID)
	userEmail := r.Context().Value("email").(string)
	// if !exists {
	// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
	// 	log.Println("handlers/func RefreshToken(): ошибка доставания id из контекста")
	// 	return
	// }
	// Создаём новый токен
	now := time.Now()
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   userEmail,
		"iat":     now.Unix(),
		"exp":     now.Add(h.tokenExpiration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(h.jwtSecret)
	if err != nil {
		fmt.Fprintf(w, "Не удалось обновить jwt токен %d", http.StatusUnprocessableEntity)
		log.Println("handlers/func RefreshToken()/SignedString(): не удалось обновить jwt токен: ", err)
		return
	}
	response := struct {
		Token      string
		Expires_in float64
		Token_type string
	}{
		Token:      tokenString,
		Expires_in: h.tokenExpiration.Seconds(),
		Token_type: "Bearer",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetUserData(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("user_id").(float64)
	user_id := int(id)
	email := r.Context().Value("email").(string)
	response := struct {
		ID    int
		Email string
	}{
		ID:    user_id,
		Email: email,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
