package agent

import (
<<<<<<< HEAD
	"fmt"
	"log"
=======
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
>>>>>>> 686799b (Pushing SuperCalculator)
	"os"
	"strconv"
	"sync"
	"time"

<<<<<<< HEAD
	pb "github.com/C1scoR/Go_Project_Yandex_Lyceum/proto"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/durationpb"
=======
	"github.com/joho/godotenv"
>>>>>>> 686799b (Pushing SuperCalculator)
)

// Task — структура задания
type Task struct {
	Id             string        `json:"id"`
	Arg1           string        `json:"Arg1"`
	Arg2           string        `json:"Arg2"`
	Operation      string        `json:"Operation"`
	Operation_time time.Duration `json:"Operation_time"`
}

<<<<<<< HEAD
// Конфиг агента собирается в ConfigFromEnv()
type Config struct {
	OrchestratorURL     string //Это для общения с оркестратором по http
	PollInterval        time.Duration
	GrpcOrchestratorURL string //Это для общения с оркестратором по grpc
=======
// Конфиг агента
type Config struct {
	Addr            string
	OrchestratorURL string
	PollInterval    time.Duration
>>>>>>> 686799b (Pushing SuperCalculator)
}

// Agent — сам агент
type Agent struct {
	config *Config
<<<<<<< HEAD
	tasks  chan *pb.Task
=======
	tasks  chan Task
>>>>>>> 686799b (Pushing SuperCalculator)
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
<<<<<<< HEAD
		tasks:  make(chan *pb.Task, 100),
=======
		tasks:  make(chan Task, 100),
>>>>>>> 686799b (Pushing SuperCalculator)
	}
}

// Загружаем конфиг
func ConfigFromEnv() *Config {
	config := new(Config)
<<<<<<< HEAD
	// home, err := os.UserHomeDir()
	// if err != nil {
	// 	log.Fatal("Ошибка получения домашней директории:", err)
	// }

	err := godotenv.Load("./internal/orchestrator/.env")
	if err != nil {
		log.Fatalln("Ошибка загрузки .env файла:", err)
	}
	//Указываю адрес grpc Оркестратора
	config.GrpcOrchestratorURL = "localhost:5050"
	config.OrchestratorURL = "http://localhost:8080/internal/task"
	config.PollInterval = 5 * time.Second // Интервал между запросами задач
	return config
}

func Getenv(key, default_value string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return default_value
}

// Запрашиваем задания у оркестратора по http. Обычный REST API хэндлер, который оказался не нужен из-за наличия gRPC
// func fetchExpressions(agent *Agent) {
// 	for {
// 		req, err := http.NewRequest("GET", agent.config.OrchestratorURL, nil)
// 		if err != nil {
// 			log.Println("Agent/fetchExpressions(): Ошибка при создании GET запроса")
// 			time.Sleep(agent.config.PollInterval)
// 			continue
// 		}

// 		req.Header.Set("X-Agent-Key", agent.config.SuperSecretKey)
// 		resp, err := http.DefaultClient.Do(req)
// 		if err != nil {
// 			log.Println("Agent/fetchExpressions():Ошибка при запросе задач:", err)
// 			time.Sleep(agent.config.PollInterval)
// 			continue
// 		}

// 		if resp.StatusCode != http.StatusCreated {
// 			log.Println("Нет задач для выполнения")
// 			resp.Body.Close()
// 			time.Sleep(agent.config.PollInterval)
// 			continue
// 		}

// 		body, _ := io.ReadAll(resp.Body)
// 		resp.Body.Close()

// 		var tasks []Task
// 		err = json.Unmarshal(body, &tasks)
// 		if err != nil {
// 			log.Println("Ошибка при разборе JSON задач:", err)
// 			continue
// 		}

// 		// Отправляем задачи по одной и не дублируем
// 		for _, task := range tasks {
// 			select {
// 			case agent.tasks <- task:
// 				log.Println("Добавлена задача", task.Id)
// 			default:
// 				log.Println("Очередь задач заполнена, пропускаем")
// 			}
// 		}

// 		time.Sleep(agent.config.PollInterval) // Ждём перед следующим запросом
// 	}
// }

// Отправляем результат обратно в оркестратор - устаревшая функция для http запросов. (Оказалась не нужна с появлением GRPC)
// func sendResult(data DataForSend, agent *Agent) {
// 	payload, err := json.Marshal(data)
// 	if err != nil {
// 		log.Println("Ошибка маршалинга результата:", err)
// 		return
// 	}
// 	log.Println("Отправляю ответ...")
// 	req, err := http.NewRequest("POST", agent.config.OrchestratorURL, bytes.NewBuffer(payload))
// 	if err != nil {
// 		log.Println("Agent()/SendResult(): Ошибка создания клиента:", err)
// 		return
// 	}
// 	req.Header.Set("X-Agent-Key", agent.config.SuperSecretKey)
// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		log.Println("Ошибка отправки результата:", err)
// 		return
// 	}
// 	defer resp.Body.Close()
// }

func convertToTimeDuration(pbDuration *durationpb.Duration) (time.Duration, error) {
	if pbDuration == nil {
		return 0, fmt.Errorf("duration is nil")
	}
	return time.Duration(pbDuration.Seconds)*time.Second + time.Duration(pbDuration.Nanos)*time.Nanosecond, nil
}

// Воркер, который вычисляет выражения
func worker(id int, tasks <-chan *pb.Task, wg *sync.WaitGroup, grpcClient pb.OrchAgentClient) {
=======
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
>>>>>>> 686799b (Pushing SuperCalculator)
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
<<<<<<< HEAD
				log.Println("Ошибка: деление на ноль в задаче", task.ID)
=======
				log.Println("Ошибка: деление на ноль в задаче", task.Id)
>>>>>>> 686799b (Pushing SuperCalculator)
				continue
			}
		default:
			log.Println("Неизвестная операция:", task.Operation)
			continue
		}

<<<<<<< HEAD
		log.Printf("Воркер %d: вычислил %s = %.2f\n", id, task.ID, result)
		d, err := convertToTimeDuration(task.OperationTime)
		if err != nil {
			log.Println("Не сработало превращение Duration: ", err)
		}
		time.Sleep(d) // Эмулируем задержку операции

		grpcsendResult(grpcClient, &pb.ResponseOfSecondServer{ID: task.ID, Result: result})
=======
		log.Printf("Воркер %d: вычислил %s = %.2f\n", id, task.Id, result)

		time.Sleep(task.Operation_time) // Эмулируем задержку операции

		sendResult(DataForSend{Id: task.Id, Result: result}, agent)
>>>>>>> 686799b (Pushing SuperCalculator)
	}
}

// Запуск агента
func RunAgent() {
	agent := NewAgent()
<<<<<<< HEAD
	conn, err := grpc.NewClient(agent.config.GrpcOrchestratorURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("could not connect to the server")
	}

	defer conn.Close()
	grpcClient := pb.NewOrchAgentClient(conn)
=======

>>>>>>> 686799b (Pushing SuperCalculator)
	// Запускаем 10 воркеров
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
<<<<<<< HEAD
		go worker(i, agent.tasks, &wg, grpcClient)
	}

	// Запускаем постоянный опрос задач
	grpcfetchExpression(grpcClient, agent)
=======
		go worker(i, agent.tasks, &wg, agent)
	}

	// Запускаем постоянный опрос задач
	fetchExpressions(agent)
>>>>>>> 686799b (Pushing SuperCalculator)

	wg.Wait()
}
