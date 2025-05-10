package main

import (
	"github.com/C1scoR/Go_Project_Yandex_Lyceum/internal/orchestrator"
)

<<<<<<< HEAD
// В этой функции main будет производиться запуск оркестратора и миграции баз данных
func main() {
	app := orchestrator.New()
	go orchestrator.RunGRPCServer()
	app.RunServer()
}
=======
func main() {
	app := orchestrator.New()
	
	//app.Run()
	app.RunServer()
}

>>>>>>> 686799b (Pushing SuperCalculator)
