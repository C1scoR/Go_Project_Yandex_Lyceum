package orchestrator

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	pb "github.com/C1scoR/Go_Project_Yandex_Lyceum/proto"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
)

type OrchAgentServer struct {
	pb.UnimplementedOrchAgentServer
	mu                       sync.Mutex
	map_for_input_variables  map[string]*Node
	map_for_output_variables map[string]*Node
}

func newOrchAgentServer() *OrchAgentServer {
	return &OrchAgentServer{
		map_for_input_variables:  make(map[string]*Node),
		map_for_output_variables: make(map[string]*Node),
	}
}

// Итерратор, чтобы передвигаться по очереди выражений
var iterrator_global int32
var mapMutex sync.Mutex // Мьютекс для защиты map

// grpc Обработчик AgentOrchGet, который отвечает на запросы агента и выдаёт ему выражения, которые можно вычислить параллельно
func (s *OrchAgentServer) AgentOrchGet(ar *pb.AgentRequest, stream pb.OrchAgent_AgentOrchGetServer) error {
	s.map_for_input_variables = make(map[string]*Node)
	//log.Println("Длина массива выражений :", len(Expressions_storage_variable.Expressions))
	if len(Expressions_storage_variable.Expressions) == int(iterrator_global) {
		//в общем если вдруг так получается, что длина и итерратор равны, то будет паника и grpc сервер дальше не попрёт
		//Потому что len = 1 и iterrator = 1, при обращении будет паника об обращении к невыделенной области памяти
		//log.Println("grpc_orchestrator/AgentOrchGet(): будет забавно если не хватало только этого условия")
		return status.Error(codes.NotFound, "Не было найдено выражения для вычисления")
	}
	for index, exp_value := range Expressions_storage_variable.Expressions {
		if exp_value.Status == StatusCreated || exp_value.Status == StatusInWork {
			break
			//return status.Error(codes.NotFound, "Выражение вычисляется пока что не нужно его забирать")
		} else if (index == len(Expressions_storage_variable.Expressions)-1) && exp_value.Status != StatusCreated {
			log.Println("grpc_orchestrator/AgentOrchGet(): Нет выражений для вычисления")
			return status.Error(codes.NotFound, "Не было найдено выражения для вычисления")
		}
	}

	err := godotenv.Load("./internal/orchestrator/.env")
	if err != nil {
		log.Fatalln("grpc_orchestrator/AgentOrchGet(): Ошибка загрузки .env файла в Обработчике:", err)
		return status.Errorf(codes.Internal, "Ошибка загрузки .env файла в Обработчике: %q", err)
	}

	if Expressions_storage_variable.Expressions[iterrator_global].Status == StatusCreated {
		//Эта часть с созданием дерева выполнится 1 раз, для 1-го выражения
		var err error

		Postfix_array, err := InfixToPostfix(Expressions_storage_variable.Expressions[iterrator_global].Result)
		if err != nil {
			log.Println("grpc_orchestrator/AgentOrchGet():Something went wrong throw convertation of expression", err)
			Expressions_storage_variable.Expressions[iterrator_global].Status = StatusFailed
			return status.Errorf(codes.Internal, "Что-то пошло не так при конвертации выражения в постфикс: %q", err)
		}
		root_of_AST_TREE = TranslateToASTTree(Postfix_array)
		Expressions_storage_variable.Expressions[iterrator_global].Status = StatusInWork
	}

	if root_of_AST_TREE == nil {
		log.Println("У нас пустой корень дерева")
		Expressions_storage_variable.Expressions[iterrator_global].Status = StatusFailed
		return status.Error(codes.NotFound, "У нас пустой корень дерева")
	}
	var wg sync.WaitGroup
	queue := make(chan *Node, 100)

	wg.Add(1)
	go CollectComputableNodes(root_of_AST_TREE, queue, &wg)
	wg.Wait()
	close(queue)
	// mutex.Lock()
	// defer mutex.Unlock()
	s.mu.Lock()
	var got_value_from_channel bool = false
	//итеррируемся по каналу, в которой записываем значения из функции CollectComputableNodes
	for Atomic_Node := range queue {
		if Atomic_Node.Status == StatusFree {
			got_value_from_channel = true
			Atomic_Node.Status = StatusRestrict
			id := uuid.New().String()
			s.map_for_output_variables[id] = Atomic_Node //вот здесь я записываю в мапу id и ноды, чтобы потом по id серверу можно было взять ноду и поменять её значение на
			//результат присланный агентом
			s.map_for_input_variables[id] = Atomic_Node
		}
	}
	if !got_value_from_channel {
		s.mu.Unlock()
		return status.Errorf(codes.NotFound, "orchestrator/AgentOrchGet():Пока что отсюда нечего забирать %d", http.StatusNoContent)
	}
	s.mu.Unlock()
	//log.Print(s.map_for_output_variables) ////////////////////
	var tasks_array []*pb.Task
	s.mu.Lock()
	for key, node_value := range s.map_for_input_variables {
		task := &pb.Task{}
		task.ID = key
		task.Arg1 = node_value.left.value
		task.Arg2 = node_value.right.value
		task.Operation = node_value.value
		if task.Operation == "+" {
			Time_Addition_ms_int_value, err := strconv.Atoi(os.Getenv("TIME_ADDITION_MS"))
			if err != nil {
				log.Println("Ошибка преобразования TIME_ADDITION_MS:", err)
			}
			d := time.Duration(Time_Addition_ms_int_value) * time.Millisecond
			task.OperationTime = durationpb.New(d)
		} else if task.Operation == "-" {
			Time_Subtraction_ms_int_value, err := strconv.Atoi(os.Getenv("TIME_SUBTRACTION_MS"))
			if err != nil {
				log.Println("Ошибка преобразования TIME_SUBTRACTION_MS:", err)
			}
			d := time.Duration(Time_Subtraction_ms_int_value) * time.Millisecond
			task.OperationTime = durationpb.New(d)
		} else if task.Operation == "*" {
			Time_Multiplication_ms_int_value, err := strconv.Atoi(os.Getenv("TIME_MULTIPLICATIONS_MS"))
			if err != nil {
				log.Println("Ошибка преобразования TIME_MULTIPLICATIONS_MS:", err)
			}
			d := time.Duration(Time_Multiplication_ms_int_value) * time.Millisecond
			task.OperationTime = durationpb.New(d)
		} else if task.Operation == "/" {
			Time_Division_ms_int_value, err := strconv.Atoi(os.Getenv("TIME_DIVISIONS_MS"))
			if err != nil {
				log.Println("Ошибка преобразования TIME_DIVISIONS_MS:", err)
			}
			d := time.Duration(Time_Division_ms_int_value) * time.Millisecond
			task.OperationTime = durationpb.New(d)
		}
		tasks_array = append(tasks_array, task)
	}
	s.mu.Unlock()
	for _, task := range tasks_array {
		if err := stream.Send(task); err != nil {
			log.Println("orchestrator/AgentOrchGet():Ошибка при стриминге данных агенту", err)
			mapMutex.Lock()
			Expressions_storage_variable.Expressions[iterrator_global].Status = StatusFailed
			mapMutex.Unlock()
			return err

		}
	}
	return nil
}

// AgentOrchPost принимает решение задачи от агента и записывает его в AST-дерево.
// Если решение уравнения окончательно, то записывает решение в базу данных.
func (s *OrchAgentServer) AgentOrchPost(_ context.Context, response *pb.ResponseOfSecondServer) (*pb.OrchResponse, error) {
	//PrintOrder(root_of_AST_TREE)
	s.mu.Lock()
	defer s.mu.Unlock()
	if node, ok := s.map_for_output_variables[response.ID]; ok {
		node.value = fmt.Sprint(response.Result)
		delete(s.map_for_output_variables, response.ID)
		//log.Println("Результат успешно обработан для задачи:", response.ID)
	} else {
		log.Println("Задача не найдена:", response.ID)
	}
	if root_of_AST_TREE.value == "+" || root_of_AST_TREE.value == "-" || root_of_AST_TREE.value == "*" || root_of_AST_TREE.value == "/" {
		log.Println("Пока преобразований корня не требуется")
	} else {
		_, err := strconv.ParseFloat(root_of_AST_TREE.value, 64)
		if err == nil {
			log.Println("Записываем выражение в корень в result массива больших выражений")
			Expressions_storage_variable.Expressions[iterrator_global].Result = root_of_AST_TREE.value
			Expressions_storage_variable.Expressions[iterrator_global].Status = StatusExecuted
			//Добавление результата выражения в БД
			log.Println("ID, который я вытаскиваю:", Expressions_storage_variable.Expressions[iterrator_global].ID)
			_, err := DB.Exec("UPDATE statements SET result = ? WHERE statement_id = ?", root_of_AST_TREE.value, Expressions_storage_variable.Expressions[iterrator_global].ID)
			if err != nil {
				log.Println("orchestrator/HandlerForCommunicationToOtherServer():\n Ошибка вставки результата выражения в БД", err)
			}
			for k := range s.map_for_output_variables {
				delete(s.map_for_output_variables, k)
			}
			log.Println("Глобальный итерратор: ", iterrator_global)
			//FIX
			if iterrator_global != int32(len(Expressions_storage_variable.Expressions)) {
				atomic.AddInt32(&iterrator_global, 1)
				log.Println("Глобальный итерратор: ", iterrator_global)
			} else {
				log.Println("Если дальше проитерируемся, то функция не выполнится")
			}
			return &pb.OrchResponse{}, nil
		} else {
			log.Println("grpc_orchestrator.go/AgentOrchPost():Ошибка с преобразованием корня дерева в число")
			return &pb.OrchResponse{}, status.Errorf(codes.Internal, "ошибка преобразования корня в число")
		}
	}
	return &pb.OrchResponse{}, nil
}

func RunGRPCServer() {
	listener, err := net.Listen("tcp", "localhost:5050")
	if err != nil {
		log.Fatalf("orchestrator/RunGRPCServer(): GRPC сервер не запустился: %q", err)
	}
	log.Printf("Сервер слушает на: %s", listener.Addr())
	grpcServer := grpc.NewServer()
	pb.RegisterOrchAgentServer(grpcServer, newOrchAgentServer())
	grpcServer.Serve(listener)
}
