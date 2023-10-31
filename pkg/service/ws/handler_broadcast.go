package ws

import (
	"encoding/json"
	"fmt"
)

type (
	BroadcastRequest struct {
		UserIDs []string `json:"user_ids"`
		Message string   `json:"message"`
	}
	EventBroadcast struct {
		UserIDs []string `json:"user_ids"`
		Message string   `json:"message"`
	}
)

func (a *wsApp) broadcastHandler(conn *WebConn, r *WebSocketRequest) error {
	req := &BroadcastRequest{}
	err := json.Unmarshal(r.Data, req)
	if err != nil {
		return fmt.Errorf("json.Unmarshal: %w", err)
	}

	if len(req.UserIDs) == 0 {
		return fmt.Errorf("userIDs is required")
	}

	if req.Message == "" {
		return fmt.Errorf("message is required")
	}

	a.Publish(req.UserIDs, NewWebSocketEvent(WebsocketEventBroadcast, &EventBroadcast{
		Message: req.Message,
	}), nil)

	return nil
}
