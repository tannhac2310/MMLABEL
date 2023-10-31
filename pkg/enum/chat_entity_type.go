package enum

import (
	"bytes"
	"fmt"
)

type ChatEntityType uint8

const (
	ChatEntityTypeInternal ChatEntityType = iota + 1
	ChatEntityTypeOpenChannelZalo
)

var ChatEntityTypeName = map[ChatEntityType]string{
	ChatEntityTypeInternal:        "internal",
	ChatEntityTypeOpenChannelZalo: "zalo",
}

var ChatEntityTypeValue = func() map[string]ChatEntityType {
	value := map[string]ChatEntityType{}
	for k, v := range ChatEntityTypeName {
		value[v] = k
		value[fmt.Sprintf("%v", k)] = k
	}

	return value
}()

func (e ChatEntityType) MarshalJSON() ([]byte, error) {
	v, ok := ChatEntityTypeName[e]
	if !ok {
		return nil, fmt.Errorf("invalid values of ChatEntityType")
	}

	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(v)
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (e *ChatEntityType) UnmarshalJSON(data []byte) error {
	data = bytes.Trim(data, "\"")
	v, ok := ChatEntityTypeValue[string(data)]
	if !ok {
		return fmt.Errorf("enum '%s' is not register, must be one of: %v", data, e.EnumDescriptions())
	}

	*e = v

	return nil
}

func (*ChatEntityType) EnumDescriptions() []string {
	vals := []string{}

	for _, name := range ChatEntityTypeName {
		vals = append(vals, name)
	}

	return vals
}
