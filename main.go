package main

import (
	_ "github.com/joho/godotenv/autoload"
	"kifuan.me/hello-board-server/models"
)

func main() {
	models.Init()

	defer models.Cleanup()
}
