package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"kifuan.me/hello-board-server/models"
)

func addMessageRoutes(rg *gin.RouterGroup) {
	g := rg.Group("messages")
	g.GET("", func(ctx *gin.Context) {
		messages, err := models.GetAllMessages()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorJSON(err))
			return
		}
		ctx.JSON(http.StatusOK, successJSON(messages))
	})

	g.POST("", func(ctx *gin.Context) {
		var err error
		var msg models.Message
		if err = ctx.ShouldBindJSON(&msg); err != nil {
			ctx.JSON(http.StatusBadRequest, errorJSON(fmt.Errorf("failed to parse request body: %w", err)))
			return
		}
		if msg, err = models.InsertMessage(msg); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorJSON(err))
			return
		}
		ctx.JSON(http.StatusOK, successJSON(nil))
	})
}
