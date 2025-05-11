package main

import (
	"github.com/C1scoR/Go_Project_Yandex_Lyceum/internal/orchestrator"
)

// В этой функции main будет производиться запуск оркестратора и миграции баз данных
func main() {
	app := orchestrator.New()
	go orchestrator.RunGRPCServer()
	app.RunServer()
}
