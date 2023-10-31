package ws

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/gorilla/websocket"
)

const (
	SendQueueSize    = 256
	SendDeadlockWarn = (SendQueueSize * 95) / 100
	WriteWait        = 30 * time.Second
	PongWait         = 100 * time.Second
	PingPeriod       = (PongWait * 6) / 10

	SocketMaxMessageSizeKb = 8 * 1024 // 8KB
)

type WebConn struct {
	ctx          context.Context
	cancelFunc   context.CancelFunc
	App          *wsApp
	WebSocket    *websocket.Conn
	Send         chan WebSocketMessage
	UserID       string
	logger       *zap.Logger
	Sequence     int64
	closeOnce    sync.Once
	endWritePump chan struct{}
	pumpFinished chan struct{}
}

func (c *WebConn) Close() {
	c.WebSocket.Close()
	c.closeOnce.Do(func() {
		c.cancelFunc()
		close(c.endWritePump)
	})
	<-c.pumpFinished
}

func (c *WebConn) Pump() {
	ch := make(chan struct{})
	go func() {
		c.writePump()
		close(ch)
	}()
	c.readPump()
	c.closeOnce.Do(func() {
		c.cancelFunc()
		close(c.endWritePump)
	})
	<-ch
	c.App.HubUnregister(c)
	close(c.pumpFinished)
}

func (c *WebConn) readPump() {
	defer func() {
		_ = c.WebSocket.Close()
	}()
	c.WebSocket.SetReadLimit(SocketMaxMessageSizeKb)
	_ = c.WebSocket.SetReadDeadline(time.Now().Add(PongWait))
	c.WebSocket.SetPongHandler(func(string) error {
		_ = c.WebSocket.SetReadDeadline(time.Now().Add(PongWait))
		return nil
	})

	for {
		var req WebSocketRequest
		if err := c.WebSocket.ReadJSON(&req); err != nil {
			// browsers will appear as CloseNoStatusReceived
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseNoStatusReceived) {
				c.logger.Debug("websocket.read: client side closed socket")
			} else {
				c.logger.Debug("websocket.read: closing websocket", zap.Error(err))
			}
			return
		}
		c.App.serveWebSocket(c, &req)
	}
}

func (c *WebConn) writePump() {
	ticker := time.NewTicker(PingPeriod)

	defer func() {
		ticker.Stop()
		_ = c.WebSocket.Close()
	}()

	for {
		select {
		case msg, ok := <-c.Send:
			if !ok {
				_ = c.WebSocket.SetWriteDeadline(time.Now().Add(WriteWait))
				_ = c.WebSocket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			evt, evtOk := msg.(*WebSocketEvent)

			var msgBytes []byte
			if evtOk {
				cpyEvt := evt.SetSequence(c.Sequence)
				msgBytes = []byte(cpyEvt.ToJSON())
				c.Sequence++
			} else {
				msgBytes = []byte(msg.ToJSON())
			}

			if len(c.Send) >= SendDeadlockWarn {
				c.logger.Warn(
					"websocket.full",
					zap.String("type", msg.EventType()),
					zap.Int("size", len(msg.ToJSON())),
				)
			}

			_ = c.WebSocket.SetWriteDeadline(time.Now().Add(WriteWait))
			if err := c.WebSocket.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
				// browsers will appear as CloseNoStatusReceived
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseNoStatusReceived) {
					c.logger.Debug("websocket.send: client side closed socket", zap.String("userID", c.UserID))
				} else {
					c.logger.Debug("websocket.send: closing websocket", zap.Error(err))
				}
				return
			}

		case <-ticker.C:
			_ = c.WebSocket.SetWriteDeadline(time.Now().Add(WriteWait))
			if err := c.WebSocket.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				// browsers will appear as CloseNoStatusReceived
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseNoStatusReceived) {
					c.logger.Debug("websocket.ticker: client side closed socket", zap.String("userID", c.UserID))
				} else {
					c.logger.Debug("websocket.ticker: closing websocket", zap.String("user_id", c.UserID), zap.Error(err))
				}
				return
			}
		case <-c.ctx.Done():
			return
		case <-c.endWritePump:
			return
		}
	}
}
