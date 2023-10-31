package ws

import (
	"encoding/json"
	"fmt"
	"io"
)

type WebSocketMessage interface {
	ToJSON() string
	EventType() string
}

type precomputedWebSocketEventJSON struct {
	Event json.RawMessage
	Data  json.RawMessage
}

// webSocketEventJSON mirrors WebSocketEvent to make some of its unexported fields serializable
type webSocketEventJSON struct {
	Event    string      `json:"event"`
	Data     interface{} `json:"data"`
	Sequence int64       `json:"seq"`
}

type WebSocketEvent struct {
	event           string
	data            interface{}
	sequence        int64
	precomputedJSON *precomputedWebSocketEventJSON
}

// PrecomputeJSON precomputes and stores the serialized JSON for all fields other than Sequence.
// This makes ToJSON much more efficient when sending the same event to multiple connections.
func (ev *WebSocketEvent) PrecomputeJSON() *WebSocketEvent {
	c := ev.copy()
	event, _ := json.Marshal(c.event)
	data, _ := json.Marshal(c.data)

	c.precomputedJSON = &precomputedWebSocketEventJSON{
		Event: json.RawMessage(event),
		Data:  json.RawMessage(data),
	}
	return c
}

func NewWebSocketEvent(event string, data interface{}) *WebSocketEvent {
	return &WebSocketEvent{
		event: event,
		data:  data,
	}
}

func (ev *WebSocketEvent) copy() *WebSocketEvent {
	c := &WebSocketEvent{
		event:           ev.event,
		data:            ev.data,
		sequence:        ev.sequence,
		precomputedJSON: ev.precomputedJSON,
	}
	return c
}

func (ev *WebSocketEvent) SetEvent(event string) *WebSocketEvent {
	c := ev.copy()
	c.event = event
	return c
}

func (ev *WebSocketEvent) SetData(data map[string]interface{}) *WebSocketEvent {
	c := ev.copy()
	c.data = data
	return c
}

func (ev *WebSocketEvent) SetSequence(seq int64) *WebSocketEvent {
	c := ev.copy()
	c.sequence = seq
	return c
}

func (ev *WebSocketEvent) IsValid() bool {
	return ev.event != ""
}

func (ev *WebSocketEvent) EventType() string {
	return ev.event
}

func (ev *WebSocketEvent) ToJSON() string {
	if ev.precomputedJSON != nil {
		return fmt.Sprintf(`{"event": %s, "data": %s, "seq": %d}`, ev.precomputedJSON.Event, ev.precomputedJSON.Data, ev.sequence)
	}
	b, _ := json.Marshal(webSocketEventJSON{
		ev.event,
		ev.data,
		ev.sequence,
	})
	return string(b)
}

func WebSocketEventFromJSON(data io.Reader) *WebSocketEvent {
	var ev WebSocketEvent
	var o webSocketEventJSON
	if err := json.NewDecoder(data).Decode(&o); err != nil {
		return nil
	}

	ev.event = o.Event
	ev.data = o.Data
	ev.sequence = o.Sequence
	return &ev
}

type WebSocketResponse struct {
	Event    string      `json:"event"`
	Action   string      `json:"action"`
	Status   string      `json:"status"`
	SeqReply int64       `json:"seq_reply,omitempty"`
	Data     interface{} `json:"data,omitempty"`
	Error    string      `json:"error,omitempty"`
}

func NewWebSocketResponse(action, status string, seqReply int64, data interface{}) *WebSocketResponse {
	return &WebSocketResponse{Action: action, Status: status, SeqReply: seqReply, Data: data}
}

func NewWebSocketError(action string, seqReply int64, err error) *WebSocketResponse {
	return &WebSocketResponse{Action: action, Status: StatusFail, SeqReply: seqReply, Error: err.Error()}
}

func (o *WebSocketResponse) EventType() string {
	return WebsocketEventResponse
}

func (o *WebSocketResponse) ToJSON() string {
	o.Event = WebsocketEventResponse
	b, _ := json.Marshal(o)
	return string(b)
}
