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
	Addr            string
	OrchestratorURL string
	PollInterval    time.Duration
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
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Ошибка получения домашней директории:", err)
	}

	err = godotenv.Load(home + "/GO_projects/internal/orchestrator/.env")
	if err != nil {
		log.Println("Ошибка загрузки .env файла:", err)
	}

	config.Addr = os.Getenv("PORT_AGENT")
	config.OrchestratorURL = "http://localhost:8080/internal/task"
	config.PollInterval = 1 * time.Second // Интервал между запросами задач
	return config
}

// Запрашиваем задания у оркестратора
func fetchExpressions(agent *Agent) {
	for {
		resp, err := http.Get(agent.config.OrchestratorURL)
		if err != nil {
			log.Println("Ошибка при запросе задач:", err)
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

	resp, err := http.Post(agent.config.OrchestratorURL, "application/json", bytes.NewBuffer(payload))
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
