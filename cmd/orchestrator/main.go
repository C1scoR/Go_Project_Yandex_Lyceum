package main

import (
	"github.com/C1scoR/Go_Project_Yandex_Lyceum/internal/orchestrator"
)

func main() {
	app := orchestrator.New()
	
	//app.Run()
	app.RunServer()
}

