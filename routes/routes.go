package routes

import (
	"os"

	"github.com/gin-gonic/gin"
)

var PORT = os.Getenv("PORT")

func Run() {
	router := gin.New()
	router.SetTrustedProxies([]string{"127.0.0.1"})
	rg := router.Group("/api")
	addMessageRoutes(rg)
	addUnsubscribeRoutes(rg)
	router.Run(PORT)
}

func errorJSON(err error) gin.H {
	return gin.H{
		"message": err.Error(),
	}
}
