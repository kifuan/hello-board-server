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
	router.Run(PORT)
}

func responseJSON(data interface{}, message string, success bool) gin.H {
	return gin.H{
		"data":    data,
		"success": success,
		"message": message,
	}
}

func errorJSON(err error) gin.H {
	return responseJSON(nil, err.Error(), false)
}

func successJSON(data interface{}) gin.H {
	return responseJSON(data, "", true)
}
