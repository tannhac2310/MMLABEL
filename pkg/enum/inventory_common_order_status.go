package enum

import (
	"bytes"
	"fmt"
)

type InventoryCommonStatusStatus uint8

const (
	InventoryCommonStatusStatusNew InventoryCommonStatusStatus = iota + 1
	InventoryCommonStatusStatusCompleted
)

var InventoryCommonStatusStatusName = map[InventoryCommonStatusStatus]string{
	InventoryCommonStatusStatusNew:       "new",
	InventoryCommonStatusStatusCompleted: "completed",
}

var InventoryCommonStatusStatusValue = func() map[string]InventoryCommonStatusStatus {
	value := map[string]InventoryCommonStatusStatus{}
	for k, v := range InventoryCommonStatusStatusName {
		value[v] = k
		value[fmt.Sprintf("%v", k)] = k
	}

	return value
}()

func (e InventoryCommonStatusStatus) MarshalJSON() ([]byte, error) {
	v, ok := InventoryCommonStatusStatusName[e]
	if !ok {
		return nil, fmt.Errorf("invalid values of InventoryCommonStatusStatus")
	}

	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(v)
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (e *InventoryCommonStatusStatus) UnmarshalJSON(data []byte) error {
	data = bytes.Trim(data, "\"")
	v, ok := InventoryCommonStatusStatusValue[string(data)]
	if !ok {
		return fmt.Errorf("enum '%s' is not register, must be one of: %v", data, e.EnumDescriptions())
	}

	*e = v

	return nil
}

func (*InventoryCommonStatusStatus) EnumDescriptions() []string {
	vals := []string{}

	for _, name := range InventoryCommonStatusStatusName {
		vals = append(vals, name)
	}

	return vals
}
