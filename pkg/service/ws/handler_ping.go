package ws

import (
	"encoding/json"
	"fmt"
)

type (
	PingRequest struct {
		Message string `json:"message"`
	}

	PingResponse struct {
		Message string `json:"message"`
	}
)

func (a *wsApp) pingHandler(conn *WebConn, r *WebSocketRequest) error {
	req := &PingRequest{}
	err := json.Unmarshal(r.Data, req)
	if err != nil {
		return fmt.Errorf("json.Unmarshal: %w", err)
	}

	if req.Message == "" {
		return fmt.Errorf("message is required")
	}

	// send response
	conn.Send <- NewWebSocketResponse(r.Action, StatusOk, r.Seq, &PingResponse{
		Message: req.Message + " pong",
	})

	return nil
}
