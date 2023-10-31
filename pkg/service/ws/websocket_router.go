package ws

import (
	"fmt"

	"go.uber.org/zap"
)

type handlerFunc func(*WebConn, *WebSocketRequest) error

func (a *wsApp) initWebsocketRouter() {
	a.handlers = make(map[string]handlerFunc)

	//TODO: add more action router
	a.handle(WebsocketActionPing, a.pingHandler)
	a.handle(WebsocketActionBroadcast, a.broadcastHandler)
}

func (a *wsApp) handle(action string, handler handlerFunc) {
	a.handlers[action] = handler
}

func (a *wsApp) serveWebSocket(conn *WebConn, r *WebSocketRequest) {
	if r.Action == "" {
		err := fmt.Errorf("action is required")
		returnWebSocketError(conn, r, err)
		return
	}

	if r.Seq <= 0 {
		err := fmt.Errorf("seq is required")
		returnWebSocketError(conn, r, err)
		return
	}

	conn.logger.Debug("received message", zap.Any("request-data", r))
	handler, ok := a.handlers[r.Action]
	if !ok {
		err := fmt.Errorf("action %s not impl", r.Action)
		returnWebSocketError(conn, r, err)
		return
	}

	err := handler(conn, r)
	if err != nil {
		returnWebSocketError(conn, r, err)
		return
	}
}

func returnWebSocketError(conn *WebConn, r *WebSocketRequest, err error) {
	conn.logger.Error(
		"websocket routing error.",
		zap.Int64("seq", r.Seq),
		zap.Error(err),
	)

	errorResp := NewWebSocketError(r.Action, r.Seq, err)

	conn.Send <- errorResp
}
