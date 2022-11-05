package routes

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"kifuan.me/hello-board-server/models"
)

func addUnsubscribeRoutes(rg *gin.RouterGroup) {
	g := rg.Group("/unsubscribe")
	g.GET("", func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Query("id"))
		if err != nil {
			ctx.String(http.StatusBadRequest, fmt.Sprintf("failed to parse id as int: %s", err))
			return
		}
		key := ctx.Query("key")

		message, err := models.GetFullMessage(id)
		if err != nil {
			logrus.Warnf("failed to find id %d in database", id)
			ctx.String(http.StatusInternalServerError, "Invalid id")
			return
		}

		if message.GenerateUnsubscribeKey() != key {
			logrus.Warnf("attempted to unsubscribe %d with wrong key %s, client ip: %s", id, key, ctx.ClientIP())
			ctx.String(http.StatusBadRequest, "Invalid key")
			return
		}

		if err := models.UnsubscribeMailNotice(message.ID); err != nil {
			logrus.Warnf("failed to update database: %s", err)
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		logrus.Infof("unsubscribed id %d", id)
		ctx.String(http.StatusOK, "Unsubscribed successfully.")
	})
}
