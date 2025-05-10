package agent

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/C1scoR/Go_Project_Yandex_Lyceum/proto"
)

func grpcfetchExpression(client pb.OrchAgentClient, agent *Agent) {
	stubb := &pb.AgentRequest{}
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		stream, err := client.AgentOrchGet(ctx, stubb)
		if err != nil {
			log.Println("grpc_agent/grpcfetchExpression(): Ошибка при вызове AgentOrchGet:", err)
			time.Sleep(agent.config.PollInterval)
			continue
		}

		// Получаем задачи от stream
		var tasks []*pb.Task
		for {
			expression, err := stream.Recv()
			if err == io.EOF {
				break // Завершаем при конце потока
			}
			if err != nil {
				log.Println("grpc_agent/grpcfetchExpression(): Ошибка при получении задачи:", err)
				time.Sleep(agent.config.PollInterval)
				break
			}

			// Обрабатываем полученное выражение
			tasks = append(tasks, expression)
		}

		// Если задач нет, ждём перед следующим запросом
		if len(tasks) == 0 {
			log.Println("grpc_agent/grpcfetchExpression(): Нет задач для обработки")
			time.Sleep(agent.config.PollInterval)
			continue
		}

		// Отправляем задачи в очередь, чтобы не дублировать
		for _, task := range tasks {
			select {
			case agent.tasks <- task:
				log.Println("Добавлена задача", task.ID)
			default:
				log.Println("Очередь задач заполнена, пропускаем")
			}
		}

		// Ждём перед следующим запросом
		time.Sleep(agent.config.PollInterval)
	}
}

func grpcsendResult(client pb.OrchAgentClient, data *pb.ResponseOfSecondServer) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("Отправляю ответ через gRPC...")
	_, err := client.AgentOrchPost(ctx, data)
	if err != nil {
		log.Println("grpc_agent/grpcsendResult(): Ошибка отправки результата:", err)
		return
	}
}
