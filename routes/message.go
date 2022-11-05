package routes

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"kifuan.me/hello-board-server/models"
)

func addMessageRoutes(rg *gin.RouterGroup) {
	g := rg.Group("messages")
	g.GET("", func(ctx *gin.Context) {
		page, err := strconv.Atoi(ctx.Query("page"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorJSON(err))
		}

		messages, err := models.GetMessages(page)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorJSON(err))
			return
		}
		ctx.JSON(http.StatusOK, messages)
	})

	g.POST("", func(ctx *gin.Context) {
		var err error
		var msg models.Message
		if err = ctx.ShouldBindJSON(&msg); err != nil {
			ctx.JSON(http.StatusBadRequest, errorJSON(fmt.Errorf("failed to parse request body: %w", err)))
			return
		}
		if err = models.InsertMessage(&msg); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorJSON(err))
			return
		}
		ctx.JSON(http.StatusOK, msg)
	})
}
