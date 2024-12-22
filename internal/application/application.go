package application

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"unicode"

	"github.com/C1scoR/Go_Project_Yandex_Lyceum/calculator"
)

type Config struct {
	Addr string
}

func ConfigFromEnv() *Config {
	config := new(Config)
	//config.Addr = os.Getenv("PORT")
	//if config.Addr == "" {
	config.Addr = "8000"
	//}
	return config
}

type Application struct {
	config *Config
}

func New() *Application {
	return &Application{
		config: ConfigFromEnv(),
	}
}

func (a *Application) Run() error {
	for {
		// читаем выражение для вычисления из командной строки
		log.Println("input expression")
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Println("failed to read expression from console")
		}
		// убираем пробелы, чтобы оставить только вычислемое выражение
		text = strings.TrimSpace(text)
		// выходим, если ввели команду "exit"
		if text == "exit" {
			log.Println("aplication was successfully closed")
			return nil
		}
		//вычисляем выражение
		result, err := calculator.Calc(text)
		if err != nil {
			log.Println(text, " calculation failed wit error: ", err)
		} else {
			log.Println(text, "=", result)
		}
	}
}

type Request struct {
	Expression string `json:"expression"`
}

func CalcHandler(w http.ResponseWriter, r *http.Request) {
	var request Request
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Printf("Error decoding request: %v", err)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if len(request.Expression) == 0 {
		log.Println("The length of given expression is 0")
		http.Error(w, fmt.Sprintf("Unprocessable entity (The length of given expression is 0), error status: %d", http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	for _, exp := range request.Expression {
		if unicode.IsLetter(exp) {
			log.Println("The expression contains letters")
			http.Error(w, fmt.Sprintf("Unprocessable entity (The expression contains letters), error status: %d", http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
			return
		}
	}

	result, err := calculator.Calc(request.Expression)
	if err != nil {
		if errors.Is(err, calculator.ErrInvalidExpression) {
			http.Error(w, fmt.Sprintf("Invalid expression: %s", err.Error()), http.StatusBadRequest)
		} else {
			http.Error(w, "unknown error", http.StatusInternalServerError)
		}
		return
	}

	fmt.Fprintf(w, "result: %f, server status: %d", result, http.StatusOK)
}

func CalcMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r)
	}
}

func (a *Application) RunServer() error {
	http.HandleFunc("/api/v1/calculate", CalcMiddleware(CalcHandler))
	return http.ListenAndServe(":"+a.config.Addr, nil)
}
