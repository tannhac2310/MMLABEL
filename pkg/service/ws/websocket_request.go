package ws

import (
	"encoding/json"
)

type WebSocketRequest struct {
	// Client-provided fields
	Seq    int64           `json:"seq"`
	Action string          `json:"action"`
	Data   json.RawMessage `json:"data"`
}

func (o *WebSocketRequest) ToJSON() string {
	b, _ := json.Marshal(o)
	return string(b)
}
