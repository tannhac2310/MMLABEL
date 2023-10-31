package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/interceptor"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/transportutil"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (w *webSocketController) HandlerWS(c *gin.Context) {
	userID := interceptor.UserIDFromCtx(c)
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	logger := ctxzap.Extract(c.Request.Context())
	wc := w.webSocketService.NewWebConn(conn, userID, logger)
	w.webSocketService.HubRegister(wc)

	wc.Pump()
}
