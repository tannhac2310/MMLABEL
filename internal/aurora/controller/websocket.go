package controller

import (
	"github.com/gin-gonic/gin"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/configs"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/ws"
)

type WebSocketController interface {
	HandlerWS(c *gin.Context)
}

type webSocketController struct {
	webSocketService ws.WebSocketService
}

func RegisterWebSocketController(
	cfg *configs.Config,
	r *gin.RouterGroup,
	webSocketService ws.WebSocketService,
) {
	g := r.Group("ws")

	var c WebSocketController = &webSocketController{
		webSocketService: webSocketService,
	}

	g.GET("", c.HandlerWS)

	if cfg.Debug {
		r.Static("ws-docs", cfg.WebSocketDocsDir)
	}
}
