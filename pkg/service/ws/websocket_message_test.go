package ws

import (
	"encoding/json"
	"log"
	"testing"
)

func TestWebSocketRequest_ToJSON(t *testing.T) {
	data := &PingRequest{
		Message: "ping",
	}

	msg, _ := json.Marshal(data)

	req := &WebSocketRequest{
		Seq:    1,
		Action: WebsocketActionPing,
		Data:   msg,
	}

	d, _ := json.Marshal(req)
	log.Println(string(d))
}
