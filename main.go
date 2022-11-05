package main

import (
	"io"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"kifuan.me/hello-board-server/models"
	"kifuan.me/hello-board-server/routes"
)

var loggerFile = &lumberjack.Logger{
	Filename:   "./server.log",
	MaxSize:    5, // MBs
	MaxBackups: 3,
	MaxAge:     28, // Days
}

func initLogger() {
	logrus.SetOutput(io.MultiWriter(os.Stdout, loggerFile))
	gin.DisableConsoleColor()
	gin.DefaultWriter = io.MultiWriter(os.Stdout, loggerFile)
	gin.DefaultErrorWriter = io.MultiWriter(os.Stdout, loggerFile)
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetFormatter(&nested.Formatter{
		HideKeys:        true,
		TimestampFormat: "2006/01/02 15:04:05",
		NoColors:        true,
	})

	if gin.Mode() == gin.DebugMode {
		logrus.SetLevel(logrus.TraceLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	logrus.Info("Logger initialized.")
}

func main() {
	initLogger()

	models.Init()
	defer models.Cleanup()

	routes.Run()
}
