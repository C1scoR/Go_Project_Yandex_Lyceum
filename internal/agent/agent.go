package agent

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

// Task — структура задания
type Task struct {
	Id             string        `json:"id"`
	Arg1           string        `json:"Arg1"`
	Arg2           string        `json:"Arg2"`
	Operation      string        `json:"Operation"`
	Operation_time time.Duration `json:"Operation_time"`
}

// Конфиг агента
type Config struct {
	OrchestratorURL string
	PollInterval    time.Duration
	SuperSecretKey  string
}

// Agent — сам агент
type Agent struct {
	config *Config
	tasks  chan Task
}

// Структура ответа с результатом
type DataForSend struct {
	Id     string  `json:"id"`
	Result float64 `json:"result"`
}

// Создаём агента
func NewAgent() *Agent {
	return &Agent{
		config: ConfigFromEnv(),
		tasks:  make(chan Task, 100),
	}
}

// Загружаем конфиг
func ConfigFromEnv() *Config {
	config := new(Config)
	// home, err := os.UserHomeDir()
	// if err != nil {
	// 	log.Fatal("Ошибка получения домашней директории:", err)
	// }

	err := godotenv.Load("C:/Users/G3eb/all_go_projects/GO_projects/internal/orchestrator/.env")
	if err != nil {
		log.Println("Ошибка загрузки .env файла:", err)
	}

	config.OrchestratorURL = "http://localhost:8080/internal/task"
	config.PollInterval = 5 * time.Second // Интервал между запросами задач
	config.SuperSecretKey = Getenv("SUPER_SECRET_KEY", "super-secret")
	return config
}

func Getenv(key, default_value string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return default_value
}

// Запрашиваем задания у оркестратора
func fetchExpressions(agent *Agent) {
	for {
		req, err := http.NewRequest("GET", agent.config.OrchestratorURL, nil)
		if err != nil {
			log.Println("Agent/fetchExpressions(): Ошибка при создании GET запроса")
			time.Sleep(agent.config.PollInterval)
			continue
		}

		req.Header.Set("X-Agent-Key", agent.config.SuperSecretKey)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println("Agent/fetchExpressions():Ошибка при запросе задач:", err)
			time.Sleep(agent.config.PollInterval)
			continue
		}

		if resp.StatusCode != http.StatusCreated {
			log.Println("Нет задач для выполнения")
			resp.Body.Close()
			time.Sleep(agent.config.PollInterval)
			continue
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		var tasks []Task
		err = json.Unmarshal(body, &tasks)
		if err != nil {
			log.Println("Ошибка при разборе JSON задач:", err)
			continue
		}

		// Отправляем задачи по одной и не дублируем
		for _, task := range tasks {
			select {
			case agent.tasks <- task:
				log.Println("Добавлена задача", task.Id)
			default:
				log.Println("Очередь задач заполнена, пропускаем")
			}
		}

		time.Sleep(agent.config.PollInterval) // Ждём перед следующим запросом
	}
}

// Отправляем результат обратно в оркестратор
func sendResult(data DataForSend, agent *Agent) {
	payload, err := json.Marshal(data)
	if err != nil {
		log.Println("Ошибка маршалинга результата:", err)
		return
	}
	log.Println("Отправляю ответ...")
	req, err := http.NewRequest("POST", agent.config.OrchestratorURL, bytes.NewBuffer(payload))
	if err != nil {
		log.Println("Agent()/SendResult(): Ошибка создания клиента:", err)
		return
	}
	req.Header.Set("X-Agent-Key", agent.config.SuperSecretKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Ошибка отправки результата:", err)
		return
	}
	defer resp.Body.Close()
}

// Воркер, который вычисляет выражения
func worker(id int, tasks <-chan Task, wg *sync.WaitGroup, agent *Agent) {
	defer wg.Done()

	for task := range tasks {
		Arg1, _ := strconv.ParseFloat(task.Arg1, 64)
		Arg2, _ := strconv.ParseFloat(task.Arg2, 64)
		var result float64

		switch task.Operation {
		case "+":
			result = Arg1 + Arg2
		case "-":
			result = Arg1 - Arg2
		case "*":
			result = Arg1 * Arg2
		case "/":
			if Arg2 != 0 {
				result = Arg1 / Arg2
			} else {
				log.Println("Ошибка: деление на ноль в задаче", task.Id)
				continue
			}
		default:
			log.Println("Неизвестная операция:", task.Operation)
			continue
		}

		log.Printf("Воркер %d: вычислил %s = %.2f\n", id, task.Id, result)

		time.Sleep(task.Operation_time) // Эмулируем задержку операции

		sendResult(DataForSend{Id: task.Id, Result: result}, agent)
	}
}

// Запуск агента
func RunAgent() {
	agent := NewAgent()

	// Запускаем 10 воркеров
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go worker(i, agent.tasks, &wg, agent)
	}

	// Запускаем постоянный опрос задач
	fetchExpressions(agent)

	wg.Wait()
}
