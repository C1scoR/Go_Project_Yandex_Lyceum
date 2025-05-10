package orchestrator

import (
	"bufio"
<<<<<<< HEAD
	"database/sql"
	"encoding/json"
	"fmt"
=======
	"encoding/json"
	"fmt"
	"io"
>>>>>>> 686799b (Pushing SuperCalculator)
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
<<<<<<< HEAD
	"time"
	"unicode"

	"github.com/C1scoR/Go_Project_Yandex_Lyceum/internal/database/config"
	"github.com/C1scoR/Go_Project_Yandex_Lyceum/internal/database/handlers"
	middleware "github.com/C1scoR/Go_Project_Yandex_Lyceum/internal/database/middleware"
	configjwt "github.com/C1scoR/Go_Project_Yandex_Lyceum/internal/jwt/config"
=======
	"sync/atomic"
	"time"
	"unicode"

>>>>>>> 686799b (Pushing SuperCalculator)
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

//var id_expression_dictionary = make(map[string]string)

type Expressions_parametres struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Result string `json:"result"`
}

type Expressions_storage struct {
	Expressions []Expressions_parametres `json:"expressions"`
}

// массивчик с выражениями от пользователя
var Expressions_storage_variable Expressions_storage

// А это корень дерева. Дерево должно будет составляться для каждого выражения
var root_of_AST_TREE *Node

func AppendToQueue(expressions_parametres Expressions_parametres) {
	Expressions_storage_variable.Expressions = append(Expressions_storage_variable.Expressions, expressions_parametres)
}

func DeletefromQueue() Expressions_parametres {
	first_value_in_queue := Expressions_storage_variable.Expressions[0]
	Expressions_storage_variable.Expressions = Expressions_storage_variable.Expressions[1:]
	return first_value_in_queue
}

type Config struct {
	Addr string
	// Time_Addition_ms       time.Duration
	// Time_Subtraction_ms    time.Duration
	// Time_Multiplication_ms time.Duration
	// Time_Division_ms       time.Duration
}
type Orchestrator struct {
	config *Config
}

// структура для закодирования тасков, которые потом пойдут к агенту
<<<<<<< HEAD
=======
type Task struct {
	Id             string        `json:"id"`
	Arg1           string        `json:"Arg1"`
	Arg2           string        `json:"Arg2"`
	Operation      string        `json:"Operation"`
	Operation_time time.Duration `json:"Operation_time"`
}
>>>>>>> 686799b (Pushing SuperCalculator)

// Структура для распарсивания отета агента
type ResponseOfSecondServer struct {
	Id     string  `json:"id"`
	Result float64 `json:"result"`
}

func ConfigFromEnv() *Config {
	config := new(Config)
<<<<<<< HEAD
	//home, err := os.UserHomeDir()
	// if err != nil {
	// 	log.Fatal("Ошибка получения домашней директории:", err)
	// }

	err := godotenv.Load("./internal/orchestrator/.env")
	if err != nil {
		log.Fatalln("orchestrator/ConfigFromEnv(): Ошибка загрузки .env файла:", err)
=======
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Ошибка получения домашней директории:", err)
	}

	err = godotenv.Load(home + "/GO_projects/internal/orchestrator/.env")
	if err != nil {
		log.Println("Ошибка загрузки .env файла:", err)
>>>>>>> 686799b (Pushing SuperCalculator)
	}
	config.Addr = os.Getenv("PORT_ORCHESTRATOR")
	return config
}

func New() *Orchestrator {
	return &Orchestrator{
		config: ConfigFromEnv(),
	}
}

func (orch *Orchestrator) Run() error {
	for {
		log.Println("input expression")
		reader := bufio.NewReader(os.Stdin)
		readed_expression, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Failed to read expression from console")
			return err
		}
		readed_expression = strings.TrimSpace(readed_expression)
		if readed_expression == "exit" {
			log.Println("Exit from the loop of expressions")
			return nil
		}
		RPN_array, err := InfixToPostfix(readed_expression)
		if err != nil {
			return err
		}
		Orchestrator_result, err := evalRPN(RPN_array, "")
		if err != nil {
			log.Println(readed_expression, " calculation failed wit error: ", err)
		} else {
			log.Println(readed_expression, "=", Orchestrator_result)
		}
	}
}

// Вот отсюда начинается реализация польской нотации
func prec(c string) int {
	if c == "*" || c == "/" {
		return 2
	} else if c == "+" || c == "-" {
		return 1
	} else {
		return -1
	}
}

func InfixToPostfix(infix_string string) ([]string, error) {
	infix_string = strings.ReplaceAll(infix_string, " ", "")
	operatorsTicker := 0
	operandsTicker := 0
	wasLastOperand := false

	for i := 0; i < len(infix_string); i++ {
		char := infix_string[i]

		if char == '+' || char == '-' || char == '*' || char == '/' {
			operatorsTicker++
			wasLastOperand = false
		} else if char >= '0' && char <= '9' || char == '.' {
			if !wasLastOperand {
				operandsTicker++
				wasLastOperand = true
			}
		} else if char == '(' || char == ')' {
			continue
		}
	}

	if operandsTicker <= operatorsTicker {
		return []string{}, ErrInvalidExpression
	}

	if infix_string[len(infix_string)-1] == '.' {
		return []string{}, ErrDotEndOfOperand
	}

	wasLastOperand = false
	var stack []rune
	var result []string
	temporary_string := ""
	for index := 0; index < len(infix_string); index++ {
		//temp_value := rune(infix_string[index])
		if infix_string[index] >= '0' && infix_string[index] <= '9' || infix_string[index] == '.' {
			//fmt.Print(temporary_string)
			for ; index < len(infix_string); index++ {
				if infix_string[index] >= '0' && infix_string[index] <= '9' && index == len(infix_string)-1 {
					temporary_string += string(infix_string[index])
					result = append(result, temporary_string)
					temporary_string = ""
				} else if infix_string[index] >= '0' && infix_string[index] <= '9' || infix_string[index] == '.' {
					temporary_string += string(infix_string[index])
				} else {
					result = append(result, temporary_string)
					temporary_string = ""
					index--
					break
				}
			}
		} else if infix_string[index] == '(' {
			stack = append(stack, rune(infix_string[index]))
		} else if infix_string[index] == ')' {
			for stack[len(stack)-1] != '(' {
				result = append(result, string(stack[len(stack)-1]))
				stack = stack[:len(stack)-1]
			}
			stack = stack[:len(stack)-1]
		} else if infix_string[index] == '+' || infix_string[index] == '-' || infix_string[index] == '/' || infix_string[index] == '*' {
			for len(stack) > 0 && (prec(string(infix_string[index])) < prec(string(stack[len(stack)-1])) || prec(string(infix_string[index])) == prec(string(stack[len(stack)-1]))) {
				result = append(result, string(stack[len(stack)-1]))
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, rune(infix_string[index]))
		}
	}
	for len(stack) > 0 {
		result = append(result, string(stack[len(stack)-1]))
		stack = stack[:len(stack)-1]
	}
	if len(result) <= 2 {
		return []string{}, ErrInvalidExpression
	}
	for _, c := range result {
		log.Printf("[%s] ", string(c))
	}
	return result, nil
}

func evalRPN(tokens []string, id string) (float64, error) {
	var stack []float64
	var mutex sync.Mutex
	var wg sync.WaitGroup
	for i := range Expressions_storage_variable.Expressions {
		if Expressions_storage_variable.Expressions[i].ID == id {
			Expressions_storage_variable.Expressions[i].Status = StatusInWork
		}
	}
	for _, token := range tokens {
		switch token {
		case "+", "-", "*", "/":
			if len(stack) < 2 {
				return 0.0, ErrInvalidExpression
			}

			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			if b == 0 && token == "/" {
				log.Println("Someone decided to divide by zero")
				return 0.0, ErrDivisionByZero
			}
			wg.Add(1)
			go func(op string, a, b float64) {
				defer wg.Done()
				temp_str := fmt.Sprintf("%.2f%s%.2f", a, op, b)
				fmt.Printf("temp_str on every iterration: %s\n", temp_str)
				var res float64
				switch op {
				case "+":
					res = a + b
				case "-":
					res = a - b
				case "*":
					res = a * b
				case "/":
					res = a / b
				}
				mutex.Lock()
				stack = append(stack, res)
				mutex.Unlock()
			}(token, a, b)
			wg.Wait()
		default:
			num, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return 0.0, ErrInvalidExpression
			}
			stack = append(stack, float64(num))
		}
	}

	if len(stack) != 1 {
		return 0.0, ErrInvalidExpression
	}

	return stack[0], nil
}

type Request struct {
	Expression string `json:"expression"`
}

<<<<<<< HEAD
var DB *sql.DB

func init() {
	var err error
	DB, err = sql.Open("sqlite3", "./internal/database/app.db")
	if err != nil {
		log.Fatalf("Orchestrator/Init(): произошла ошибка подключения к БД: %q", err)
	}
}

// OrchestratorHandler Получает выражение на вход и добалвяет его в очередь для вычисления
=======
>>>>>>> 686799b (Pushing SuperCalculator)
func OrchestratorHandler(w http.ResponseWriter, r *http.Request) {
	var request Request
	//проверка на json
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Printf("Error decoding request (Orchestrator): %v", err)
		http.Error(w, "Invalid json format: ", http.StatusUnprocessableEntity)
		return
	}
	defer r.Body.Close()
	//проверка на длину
	if len(request.Expression) == 0 {
		log.Println("The length of given expression is 0")
		http.Error(w, "Unprocessable entity (The length of given expression is 0), error status: ", http.StatusUnprocessableEntity)
		return
	}
	if request.Expression[len(request.Expression)-1] == '+' || request.Expression[len(request.Expression)-1] == '-' || request.Expression[len(request.Expression)-1] == '*' || request.Expression[len(request.Expression)-1] == '/' {
		log.Println("The expression contains operator in end of expression")
		http.Error(w, "Unprocessable entity (The expression contains operator in end of expression), error status: ", http.StatusUnprocessableEntity)
		return
	}
	//проверка на буквы
	for _, exp := range request.Expression {
		if unicode.IsLetter(exp) {
			log.Println("The expression contains letters")
			http.Error(w, fmt.Sprintf("Unprocessable entity (The expression contains letters), error status: %d", http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
			return
		}
	}
<<<<<<< HEAD
	id_float64 := r.Context().Value("user_id").(float64)
	user_id := int(id_float64)
	id := uuid.New().String()
	//fmt.Println("ID который я должен вытащить: ", id)
=======

	id := uuid.New().String()
>>>>>>> 686799b (Pushing SuperCalculator)

	//Вот тут мы отправляем пользователю ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id": id,
	})
<<<<<<< HEAD

	//ВАЖНО!!! ДОБАВЛЯЮ ЭЛЕМЕНТ В ОЧЕРЕДЬ
	AppendToQueue(Expressions_parametres{id, StatusCreated, request.Expression})
	//log.Printf("orchestrator/OrchestratorHandler(): элемент был добавлен в очередь, %T длина очереди: %d", Expressions_storage_variable.Expressions, len(Expressions_storage_variable.Expressions))
	//atomic.AddInt32(&iterrator_global, 1)
	//Добавляю выражение в БД
	_, err = DB.Exec("INSERT INTO statements (user_id, statement_id, statement) VALUES (?, ?, ?)", user_id, id, request.Expression)
	if err != nil {
		log.Println("Ошибка вставки выражения в Базу Данных", err)
		http.Error(w, "Ошибка вставки в БД", http.StatusInternalServerError)
		return
	}

}

// Эта функция параллельно и рекурсивно раскладывает дерево
// Если она находит свободную связку для вычисления (по типу 4+5), то передаёт её в канал
=======
	//ВАЖНО!!! ДОБАВЛЯЮ ЭЛЕМЕНТ В ОЧЕРЕДЬ
	AppendToQueue(Expressions_parametres{id, StatusCreated, request.Expression})
	//atomic.AddInt32(&iterrator_global, 1)

}

>>>>>>> 686799b (Pushing SuperCalculator)
func CollectComputableNodes(node *Node, queue chan *Node, wg *sync.WaitGroup) {
	defer wg.Done() // Уменьшаем счетчик wg при выходе из горутины

	if node == nil {
		return
	}

	// Проверяем, можно ли вычислить узел
	if (node.value == "+" || node.value == "-" || node.value == "*" || node.value == "/") &&
		node.left != nil && node.right != nil &&
		(node.left.value != "+" && node.left.value != "-" && node.left.value != "*" && node.left.value != "/") &&
		(node.right.value != "+" && node.right.value != "-" && node.right.value != "*" && node.right.value != "/") &&
		node.Status == StatusFree {

		select {
		case queue <- node: // Отправляем узел в канал
		default:
			log.Println("Очередь заполнена, узел не был добавлен")
		}
		return
	}

	// Обходим поддеревья только если они не nil
	if node.left != nil {
		wg.Add(1)
		go CollectComputableNodes(node.left, queue, wg)
	}
	if node.right != nil {
		wg.Add(1)
		go CollectComputableNodes(node.right, queue, wg)
	}
}

<<<<<<< HEAD
//глобальные переменные для старого оркестратора
// var (
// 	map_for_output_variables = make(map[string]*Node)
// 	mutex                    sync.Mutex // Мьютекс для других синхронизаций
// )

type Task struct {
	Id             string        `json:"id"`
	Arg1           string        `json:"Arg1"`
	Arg2           string        `json:"Arg2"`
	Operation      string        `json:"Operation"`
	Operation_time time.Duration `json:"Operation_time"`
}

//Оркестратор для общения по http - не нужен с появлением gRPC
// func HandlerForCommunicationToOtherServer(w http.ResponseWriter, r *http.Request) {
// 	map_for_input_variables := make(map[string]*Node)
// 	for index, exp_value := range Expressions_storage_variable.Expressions {
// 		if exp_value.Status == StatusCreated || exp_value.Status == StatusInWork {
// 			break
// 		} else if (index == len(Expressions_storage_variable.Expressions)-1) && exp_value.Status != StatusCreated {
// 			log.Println("Нет выражений для вычисления")
// 			http.Error(w, "Нет выражений для вычисления абсолютно", http.StatusNotFound)
// 		}
// 	}
// 	// home, err := os.UserHomeDir()
// 	// if err != nil {
// 	// 	log.Fatal("Ошибка получения домашней директории в Обработчике:", err)
// 	// 	return
// 	// }

// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatalln("orchestrator/HandlerForCommunicationToOtherServer(): Ошибка загрузки .env файла в Обработчике:", err)
// 		return
// 	}

// 	iterrator := atomic.LoadInt32(&iterrator_global)
// 	//log.Print(iterrator)
// 	if Expressions_storage_variable.Expressions[iterrator].Status == StatusCreated {
// 		//Эта часть с созданием дерева выполнится 1 раз, для 1-го выражения
// 		var err error
// 		Postfix_array, err := InfixToPostfix(Expressions_storage_variable.Expressions[iterrator].Result)
// 		if err != nil {
// 			log.Println("Something went wrong throw convertation of expression")
// 			Expressions_storage_variable.Expressions[iterrator].Status = StatusFailed
// 			http.Error(w, fmt.Sprintf("Something went wrong throw convertation of expression, error status: %d", http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
// 			return
// 		}
// 		root_of_AST_TREE = TranslateToASTTree(Postfix_array)
// 		Expressions_storage_variable.Expressions[iterrator].Status = StatusInWork
// 	}
// 	//Если приходит Get запрос, то мы разбиваем наше дерево. Находим задачи, которые можно вычислить параллельно, через канал добавляем их в 2 очереди
// 	//и в тело post запроса кладём сразу столько элементов, сколько накопилось в очереди
// 	//Одну из очередей чистим, вторую чистим когда придёт post-запрос с ответом
// 	if r.Method == http.MethodGet {
// 		if root_of_AST_TREE == nil {
// 			log.Println("У нас пустой корень дерева")
// 			http.Error(w, "Пустой корень дерева", http.StatusBadRequest)
// 			Expressions_storage_variable.Expressions[iterrator].Status = StatusFailed
// 			return
// 		}
// 		var wg sync.WaitGroup
// 		queue := make(chan *Node, 100)

// 		wg.Add(1)
// 		go CollectComputableNodes(root_of_AST_TREE, queue, &wg)
// 		wg.Wait()
// 		close(queue)
// 		mutex.Lock()
// 		defer mutex.Unlock()
// 		mapMutex.Lock()
// 		//итеррируемся по каналу, в которой записываем значения из функции CollectComputableNodes
// 		for Atomic_Node := range queue {
// 			if Atomic_Node.Status == StatusFree {
// 				Atomic_Node.Status = StatusRestrict
// 				id := uuid.New().String()
// 				map_for_output_variables[id] = Atomic_Node //вот здесь я записываю в мапу id и ноды, чтобы потом по id агента сервера можно было взять ноду и поменять её значение на
// 				//результат присланный агентом
// 				map_for_input_variables[id] = Atomic_Node
// 			}
// 		}
// 		mapMutex.Unlock()
// 		log.Print(map_for_output_variables) ////////////////////
// 		var tasks_array []Task
// 		mapMutex.Lock()
// 		for key, node_value := range map_for_input_variables {
// 			task := Task{}
// 			task.Id = key
// 			task.Arg1 = node_value.left.value
// 			task.Arg2 = node_value.right.value
// 			task.Operation = node_value.value
// 			if task.Operation == "+" {
// 				Time_Addition_ms_int_value, err := strconv.Atoi(os.Getenv("TIME_ADDITION_MS"))
// 				if err != nil {
// 					log.Println("Ошибка преобразования TIME_ADDITION_MS:", err)
// 				}
// 				task.Operation_time = time.Duration(Time_Addition_ms_int_value) * time.Millisecond
// 			} else if task.Operation == "-" {
// 				Time_Subtraction_ms_int_value, err := strconv.Atoi(os.Getenv("TIME_SUBTRACTION_MS"))
// 				if err != nil {
// 					log.Println("Ошибка преобразования TIME_SUBTRACTION_MS:", err)
// 				}
// 				task.Operation_time = time.Duration(Time_Subtraction_ms_int_value) * time.Millisecond
// 			} else if task.Operation == "*" {
// 				Time_Multiplication_ms_int_value, err := strconv.Atoi(os.Getenv("TIME_MULTIPLICATIONS_MS"))
// 				if err != nil {
// 					log.Println("Ошибка преобразования TIME_MULTIPLICATIONS_MS:", err)
// 				}
// 				task.Operation_time = time.Duration(Time_Multiplication_ms_int_value) * time.Millisecond

// 			} else if task.Operation == "/" {
// 				Time_Division_ms_int_value, err := strconv.Atoi(os.Getenv("TIME_DIVISIONS_MS"))
// 				if err != nil {
// 					log.Println("Ошибка преобразования TIME_DIVISIONS_MS:", err)
// 				}
// 				task.Operation_time = time.Duration(Time_Division_ms_int_value) * time.Millisecond
// 			}
// 			tasks_array = append(tasks_array, task)
// 		}
// 		mapMutex.Unlock()
// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusCreated)
// 		err := json.NewEncoder(w).Encode(tasks_array)
// 		if err != nil {
// 			log.Println("Не удалось закодировать json в Главном Хэндлере")
// 			Expressions_storage_variable.Expressions[iterrator].Status = StatusFailed
// 			http.Error(w, "unknown error", http.StatusInternalServerError)
// 			return
// 		}
// 		//Отправили массив с тасками, и такая отправка может быть на 2+ итерациях. То есть дерево реально обходится несколько раз.
// 	}
// 	if r.Method == http.MethodPost {
// 		//Вот к нам и пришёл запросик с данными
// 		body, err := io.ReadAll(r.Body)
// 		if err != nil {
// 			log.Println("Что-то не так с чтением тела запроса", err)
// 			http.Error(w, "Произошла ошибка с распаршиванием данных", http.StatusUnprocessableEntity)
// 			Expressions_storage_variable.Expressions[iterrator].Status = StatusFailed
// 			return
// 		}
// 		defer r.Body.Close()
// 		var response ResponseOfSecondServer
// 		err = json.Unmarshal(body, &response)
// 		log.Println("Результат, который пришёл (надеюсь): ", response.Id, response.Result)
// 		if err != nil {
// 			log.Println("Что-то не так с unmarshal json", err)
// 			Expressions_storage_variable.Expressions[iterrator].Status = StatusFailed
// 			return
// 		}
// 		PrintOrder(root_of_AST_TREE)
// 		//мы распарсили post-запрос, дальше нам нужно вставить resultы на место Node, которые мы храним в output_map
// 		mutex.Lock()
// 		defer mutex.Unlock()
// 		if node, ok := map_for_output_variables[response.Id]; ok {
// 			node.value = fmt.Sprint(response.Result)
// 			delete(map_for_output_variables, response.Id)
// 			log.Println("Результат успешно обработан для задачи:", response.Id)
// 		} else {
// 			log.Println("Задача не найдена:", response.Id)
// 		}
// 	}
// 	if root_of_AST_TREE.value == "+" || root_of_AST_TREE.value == "-" || root_of_AST_TREE.value == "*" || root_of_AST_TREE.value == "/" {
// 		log.Println("Пока преобразований корня не требуется")
// 	} else {
// 		_, err = strconv.ParseFloat(root_of_AST_TREE.value, 64)
// 		if err == nil {
// 			log.Println("Записываем выражение в корень в result массива больших выражений")
// 			Expressions_storage_variable.Expressions[iterrator].Result = root_of_AST_TREE.value
// 			Expressions_storage_variable.Expressions[iterrator].Status = StatusExecuted
// 			//Добавление результата выражения в БД
// 			_, err := DB.Exec("UPDATE statements SET result = ? WHERE statement_id = ?", root_of_AST_TREE.value, Expressions_storage_variable.Expressions[iterrator].ID)
// 			if err != nil {
// 				log.Println("orchestrator/HandlerForCommunicationToOtherServer():\n Ошибка вставки результата выражения в БД", err)
// 			}
// 			for k := range map_for_output_variables {
// 				delete(map_for_output_variables, k)
// 			}
// 			if iterrator != int32(len(Expressions_storage_variable.Expressions)-1) {
// 				atomic.AddInt32(&iterrator_global, 1)
// 			} else {
// 				log.Println("Если дальше проитерируемся, то функция не выполнится")
// 			}
// 			return
// 		} else {
// 			http.Error(w, "Ошибка с преобразованием корня дерева в число", http.StatusUnprocessableEntity)
// 			log.Println("Ошибка с преобразованием корня дерева в число")
// 			return
// 		}
// 	}
// }

func GetExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("user_id").(float64)
	user_id := int(id)
	var statement_id, statement, result string
	rows, err := DB.Query("SELECT statement_id, statement, result FROM statements WHERE user_id = ? AND result IS NOT NULL", user_id)
	if err == sql.ErrNoRows {
		log.Println("orchestrator/GetExpressionsHandler(): В Базу данных ещё не было добавлено ни одного выражения: ErrNoRows")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
		jsonData := struct {
			Error string
		}{
			Error: "Вы ещё не делали запросы с выражениями",
		}
		json.NewEncoder(w).Encode(jsonData)
		return
	}
	if err != nil {
		log.Println("orchestrator/GetExpressionsHandler(): Неизвестная ошибка при получении строк БД", err)
		http.Error(w, " Неизвестная ошибка при получении строк БД", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	type Expression struct {
		Statement_id string `json:"statement_id"`
		Statement    string `json:"statement"`
		Result       string `json:"result"`
	}

	var expressions []Expression

	for rows.Next() {
		err := rows.Scan(&statement_id, &statement, &result)
		if err != nil {
			log.Println("orchestrator/GetExpressionsHandler(): Неизвестная ошибка итерации по строкам", err)
			http.Error(w, " Неизвестная ошибка", http.StatusInternalServerError)
			return
		}
		expressions = append(expressions, Expression{Statement_id: statement_id, Statement: statement, Result: result})
	}
	if len(expressions) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
		jsonData := struct {
			Error string
		}{
			Error: "Вы ещё не делали запросы на выражения",
		}
		json.NewEncoder(w).Encode(jsonData)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(expressions)
=======
var (
	map_for_output_variables = make(map[string]*Node)
	mapMutex                 sync.Mutex // Мьютекс для защиты map
	iterrator_global         int32
	mutex                    sync.Mutex // Мьютекс для других синхронизаций
)

func HandlerForCommunicationToOtherServer(w http.ResponseWriter, r *http.Request) {
	map_for_input_variables := make(map[string]*Node)
	for index, exp_value := range Expressions_storage_variable.Expressions {
		if exp_value.Status == StatusCreated || exp_value.Status == StatusInWork {
			break
		} else if (index == len(Expressions_storage_variable.Expressions)-1) && exp_value.Status != StatusCreated {
			log.Println("Нет выражений для вычисления")
			http.Error(w, "Нет выражений для вычисления абсолютно", http.StatusNotFound)
		}
	}
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Ошибка получения домашней директории в Обработчике:", err)
		return
	}

	err = godotenv.Load(home + "/GO_projects/internal/orchestrator/.env")
	if err != nil {
		log.Println("Ошибка загрузки .env файла в Обработчике:", err)
		return
	}

	iterrator := atomic.LoadInt32(&iterrator_global)
	//log.Print(iterrator)
	if Expressions_storage_variable.Expressions[iterrator].Status == StatusCreated {
		//Эта часть с созданием дерева выполнится 1 раз, для 1-го выражения
		var err error
		Postfix_array, err := InfixToPostfix(Expressions_storage_variable.Expressions[iterrator].Result)
		if err != nil {
			log.Println("Something went wrong throw convertation of expression")
			Expressions_storage_variable.Expressions[iterrator].Status = StatusFailed
			http.Error(w, fmt.Sprintf("Something went wrong throw convertation of expression, error status: %d", http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
			return
		}
		root_of_AST_TREE = TranslateToASTTree(Postfix_array)
		Expressions_storage_variable.Expressions[iterrator].Status = StatusInWork
	}
	//Если приходит Get запрос, то мы разбиваем наше дерево. Находим задачи, которые можно вычислить параллельно, через канал добавляем их в 2 очереди
	//и в тело post запроса кладём сразу столько элементов, сколько накопилось в очереди
	//Одну из очередей чистим, вторую чистим когда придёт post-запрос с ответом
	if r.Method == http.MethodGet {
		if root_of_AST_TREE == nil {
			log.Println("У нас пустой корень дерева")
			http.Error(w, "Пустой корень дерева", http.StatusBadRequest)
			Expressions_storage_variable.Expressions[iterrator].Status = StatusFailed
			return
		}
		var wg sync.WaitGroup
		queue := make(chan *Node, 100)

		wg.Add(1)
		go CollectComputableNodes(root_of_AST_TREE, queue, &wg)
		wg.Wait()
		close(queue)
		mutex.Lock()
		defer mutex.Unlock()
		mapMutex.Lock()
		for Atomic_Node := range queue {
			if Atomic_Node.Status == StatusFree {
				Atomic_Node.Status = StatusRestrict
				id := uuid.New().String()
				map_for_output_variables[id] = Atomic_Node //вот здесь я записываю в мапу id и ноды, чтобы потом по id агента сервера можно было взять ноду и поменять её значение на
				//результат присланный агентом
				map_for_input_variables[id] = Atomic_Node
			}
		}
		mapMutex.Unlock()
		log.Print(map_for_output_variables) ////////////////////
		var tasks_array []Task
		mapMutex.Lock()
		for key, node_value := range map_for_input_variables {
			task := Task{}
			task.Id = key
			task.Arg1 = node_value.left.value
			task.Arg2 = node_value.right.value
			task.Operation = node_value.value
			if task.Operation == "+" {
				Time_Addition_ms_int_value, err := strconv.Atoi(os.Getenv("TIME_ADDITION_MS"))
				if err != nil {
					log.Println("Ошибка преобразования TIME_ADDITION_MS:", err)
				}
				task.Operation_time = time.Duration(Time_Addition_ms_int_value) * time.Millisecond
			} else if task.Operation == "-" {
				Time_Subtraction_ms_int_value, err := strconv.Atoi(os.Getenv("TIME_SUBTRACTION_MS"))
				if err != nil {
					log.Println("Ошибка преобразования TIME_SUBTRACTION_MS:", err)
				}
				task.Operation_time = time.Duration(Time_Subtraction_ms_int_value) * time.Millisecond
			} else if task.Operation == "*" {
				Time_Multiplication_ms_int_value, err := strconv.Atoi(os.Getenv("TIME_MULTIPLICATIONS_MS"))
				if err != nil {
					log.Println("Ошибка преобразования TIME_MULTIPLICATIONS_MS:", err)
				}
				task.Operation_time = time.Duration(Time_Multiplication_ms_int_value) * time.Millisecond

			} else if task.Operation == "/" {
				Time_Division_ms_int_value, err := strconv.Atoi(os.Getenv("TIME_DIVISIONS_MS"))
				if err != nil {
					log.Println("Ошибка преобразования TIME_DIVISIONS_MS:", err)
				}
				task.Operation_time = time.Duration(Time_Division_ms_int_value) * time.Millisecond
			}
			tasks_array = append(tasks_array, task)
		}
		mapMutex.Unlock()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		err := json.NewEncoder(w).Encode(tasks_array)
		if err != nil {
			log.Println("Не удалось закодировать json в Главном Хэндлере")
			Expressions_storage_variable.Expressions[iterrator].Status = StatusFailed
			http.Error(w, "unknown error", http.StatusInternalServerError)
			return
		}
		//Отправили массив с тасками, и такая отправка может быть на 2+ итерациях. То есть дерево реально обходится несколько раз.
	}
	if r.Method == http.MethodPost {
		//Вот к нам и пришёл запросик с данными
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("Что-то не так с чтением тела запроса", err)
			http.Error(w, "Произошла ошибка с распаршиванием данных", http.StatusUnprocessableEntity)
			Expressions_storage_variable.Expressions[iterrator].Status = StatusFailed
			return
		}
		defer r.Body.Close()
		var response ResponseOfSecondServer
		err = json.Unmarshal(body, &response)
		if err != nil {
			log.Println("Что-то не так с unmarshal json", err)
			Expressions_storage_variable.Expressions[iterrator].Status = StatusFailed
			return
		}
		PrintOrder(root_of_AST_TREE)
		//мы распарсили post-запрос, дальше нам нужно вставить resultы на место Node, которые мы храним в output_map
		mutex.Lock()
		defer mutex.Unlock()
		if node, ok := map_for_output_variables[response.Id]; ok {
			node.value = fmt.Sprint(response.Result)
			delete(map_for_output_variables, response.Id)
			log.Println("Результат успешно обработан для задачи:", response.Id)
		} else {
			log.Println("Задача не найдена:", response.Id)
		}
	}
	if root_of_AST_TREE.value == "+" || root_of_AST_TREE.value == "-" || root_of_AST_TREE.value == "*" || root_of_AST_TREE.value == "/" {
		log.Println("Пока преобразований корня не требуется")
	} else {
		_, err = strconv.ParseFloat(root_of_AST_TREE.value, 64)
		if err == nil {
			log.Println("Записываем выражение в корень в result массива больших выражений")
			Expressions_storage_variable.Expressions[iterrator].Result = root_of_AST_TREE.value
			Expressions_storage_variable.Expressions[iterrator].Status = StatusExecuted
			for k := range map_for_output_variables {
				delete(map_for_output_variables, k)
			}
			if iterrator != int32(len(Expressions_storage_variable.Expressions)-1) {
				atomic.AddInt32(&iterrator_global, 1)
			} else {
				log.Println("Если дальше проитерируемся, то функция не выполнится")
			}
			return
		} else {
			http.Error(w, "Ошибка с преобразованием корня дерева в число", http.StatusUnprocessableEntity)
			log.Println("Ошибка с преобразованием корня дерева в число")
			return
		}
	}
	//если ошибки нет, значит в корне уже стоит число
	//если стоит число, то мы в массив expression подставляем result.
	//Только вот вопрос а как дальше итерироваться, потому что это же пизда, ну типа либо заводить глобальный счётчик, либо я хуй знает, как стэк эту залупу юзать, где
	//мы вставляем элемент и берём последний вставленный на вычисления, а как только досчитали, то прихуячиваем его в самое начало массива

}

func GetExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	if len(Expressions_storage_variable.Expressions) == 0 {
		log.Println("Запрос на данные, но данных нет")

		response := map[string]string{"error": "There are no expressions in Database"}
		jsonData, err := json.Marshal(response)
		if err != nil {
			http.Error(w, `{"error": "unknown error"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write(jsonData)
		return
	} else {
		for i := 0; i < len(Expressions_storage_variable.Expressions); i++ {
			jsonData, err := json.MarshalIndent(Expressions_storage_variable.Expressions[i], "", " ")
			if err != nil {
				log.Println("Что-то пошло не так при маршализации json")
				http.Error(w, "unknown error", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, err = w.Write(jsonData)
			if err != nil {
				http.Error(w, "Ошибка при отправке данных", http.StatusInternalServerError)
				log.Println("Ошибка при отправке данных:", err)
				return
			}
		}
	}
>>>>>>> 686799b (Pushing SuperCalculator)
}

func GetExpressionByIdHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
<<<<<<< HEAD
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	id := parts[4]
	user_id_float64 := r.Context().Value("user_id").(float64)
	userID := int(user_id_float64)
	var uID int
	var statement_id, statement, result string
	err := DB.QueryRow("SELECT user_id, statement_id, statement, result FROM statements WHERE statement_id = ?", id).Scan(&uID, &statement_id, &statement, &result)
	if userID != uID {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, "It's not yours expression: ", http.StatusForbidden)
		return
	}
	if err == sql.ErrNoRows {
		log.Println("orchestrator/GetExpressionByIdHandler(): В Базу данных ещё не было добавлено такого выражения: ErrNoRows")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
		jsonData := struct {
			Error string
		}{
			Error: "Такого выражения не существует",
		}
		json.NewEncoder(w).Encode(jsonData)
		return
	}
	if err != nil {
		log.Println("orchestrator/GetExpressionByIdHandler() Неизвестная ошибка при получении строки БД", err)
		http.Error(w, " Неизвестная ошибка при получении строки БД", http.StatusInternalServerError)
		return
	}

	type Expression struct {
		Statement_id string `json:"statement_id"`
		Statement    string `json:"statement"`
		Result       string `json:"result"`
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Expression{Statement_id: statement_id, Statement: statement, Result: result})
}

func (orch *Orchestrator) RunServer() error {
	DBcfg := config.ConfigDB()
	defer DBcfg.DB.Close()
	JWTcfg, err := configjwt.LoadJWT()
	authHandler := handlers.NewAuthHandler(DBcfg, []byte(JWTcfg.Secret))
	if err != nil {
		return err
	}
	log.Println("Сервера маму люблю, порт: ", orch.config.Addr)
	mux := http.NewServeMux()
	timeout := time.Second * 10
	//mux.Handle("/internal/task", middleware.AgentKeyMiddleware("super-secret")(http.HandlerFunc(HandlerForCommunicationToOtherServer)))
	mux.Handle("/api/v1/login", middleware.ContextMiddleware(timeout)(http.HandlerFunc(authHandler.Login)))
	mux.Handle("/api/v1/register", middleware.ContextMiddleware(timeout)(http.HandlerFunc(authHandler.Register)))
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("/api/v1/refresh-token", authHandler.RefreshToken)
	protectedMux.HandleFunc("/api/v1/calculate", OrchestratorHandler)
	protectedMux.HandleFunc("/api/v1/expressions", GetExpressionsHandler)
	protectedMux.HandleFunc("/api/v1/expressions/", GetExpressionByIdHandler)
	protectedMux.HandleFunc("/api/v1/profile", authHandler.GetUserData)
	protectedHandler := middleware.AuthMiddleware([]byte(JWTcfg.Secret))(protectedMux)
	finalHandler := SplitHandler(mux, protectedHandler)
	err = http.ListenAndServe(":"+orch.config.Addr, finalHandler)
=======
		log.Println()
		http.Error(w, "Invalid URL", http.StatusBadRequest)
	}
	id := parts[4]
	for _, expression := range Expressions_storage_variable.Expressions {
		if expression.ID == id {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			jsonData, err := json.MarshalIndent(expression, "", " ")
			if err != nil {
				http.Error(w, "Ошибка при отправке данных", http.StatusInternalServerError)
				log.Println("Данные не закодировались в json, в GetExpressionByIdHandler")
				return
			}
			_, err = w.Write([]byte(jsonData))
			if err != nil {
				http.Error(w, "Ошибка при отправке данных", http.StatusInternalServerError)
				log.Println("Что-то случилось при отправке json-данных пользователю", err)
				return
			}
			return
		}
	}
	http.Error(w, "Такого выражения нет", http.StatusNotFound)
	log.Println("Пользователь захотел несуществующее выражение")
}

func (orch *Orchestrator) RunServer() error {
	log.Println("Сервера маму люблю, порт: ", orch.config.Addr)
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/calculate", OrchestratorHandler)
	mux.HandleFunc("/api/v1/expressions", GetExpressionsHandler)
	mux.HandleFunc("/api/v1/expressions/", GetExpressionByIdHandler)
	mux.HandleFunc("/internal/task", HandlerForCommunicationToOtherServer)
	err := http.ListenAndServe(":"+orch.config.Addr, mux)
>>>>>>> 686799b (Pushing SuperCalculator)
	if err != nil {
		return err
	}
	return nil
}
<<<<<<< HEAD

func SplitHandler(public, protected http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/v1/refresh-token") ||
			strings.HasPrefix(r.URL.Path, "/api/v1/profile") ||
			strings.HasPrefix(r.URL.Path, "/api/v1/calculate") ||
			strings.HasPrefix(r.URL.Path, "/api/v1/expressions") {
			protected.ServeHTTP(w, r)
			return
		}
		public.ServeHTTP(w, r)
	})
}
=======
>>>>>>> 686799b (Pushing SuperCalculator)
