package main

import (
	_ "github.com/joho/godotenv/autoload"
	"kifuan.me/hello-board-server/models"
	"kifuan.me/hello-board-server/routes"
)

func main() {
	models.Init()
	defer models.Cleanup()

	routes.Run()
}
