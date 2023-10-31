package enum

import (
	"bytes"
	"fmt"
)

type MessageStatus uint8

const (
	MessageStatusReceived MessageStatus = iota + 1
	MessageStatusDelivered
)

var MessageStatusName = map[MessageStatus]string{
	MessageStatusReceived:  "received",
	MessageStatusDelivered: "delivered",
}

var MessageStatusValue = func() map[string]MessageStatus {
	value := map[string]MessageStatus{}
	for k, v := range MessageStatusName {
		value[v] = k
		value[fmt.Sprintf("%v", k)] = k
	}

	return value
}()

func (e MessageStatus) MarshalJSON() ([]byte, error) {
	v, ok := MessageStatusName[e]
	if !ok {
		return nil, fmt.Errorf("invalid values of MessageStatus")
	}

	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(v)
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (e *MessageStatus) UnmarshalJSON(data []byte) error {
	data = bytes.Trim(data, "\"")
	v, ok := MessageStatusValue[string(data)]
	if !ok {
		return fmt.Errorf("enum '%s' is not register, must be one of: %v", data, e.EnumDescriptions())
	}

	*e = v

	return nil
}

func (*MessageStatus) EnumDescriptions() []string {
	vals := []string{}

	for _, name := range MessageStatusName {
		vals = append(vals, name)
	}

	return vals
}
