package orchestrator

import (
	"context"
	"log"
	"net"
	"testing"

	pb "github.com/C1scoR/Go_Project_Yandex_Lyceum/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterOrchAgentServer(s, newOrchAgentServer()) // <-- Подразумевается, что ты реализовал это
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestAgentOrchGet(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := pb.NewOrchAgentClient(conn)

	tests := map[string]struct {
		in       *pb.AgentRequest
		expected []*pb.Task
		errMsg   string
	}{
		"Must Success with no expressions": {
			in:       &pb.AgentRequest{},
			expected: []*pb.Task{},
			errMsg:   "rpc error: code = NotFound desc = Не было найдено выражения для вычисления",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := client.AgentOrchGet(ctx, tt.in)

			if err != nil {
				if err.Error() != tt.errMsg {
					t.Errorf("unexpected error:\ngot:  %v\nwant: %v", err.Error(), tt.errMsg)
				}
				return
			}

		})
	}
}
