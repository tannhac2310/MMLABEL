package enum

import (
	"bytes"
	"fmt"
)

type InventoryCommonStatus uint8

const (
	InventoryCommonStatusStatusNew InventoryCommonStatus = iota + 1
	InventoryCommonStatusStatusCompleted
)

var InventoryCommonStatusStatusName = map[InventoryCommonStatus]string{
	InventoryCommonStatusStatusNew:       "new",
	InventoryCommonStatusStatusCompleted: "completed",
}

var InventoryCommonStatusStatusValue = func() map[string]InventoryCommonStatus {
	value := map[string]InventoryCommonStatus{}
	for k, v := range InventoryCommonStatusStatusName {
		value[v] = k
		value[fmt.Sprintf("%v", k)] = k
	}

	return value
}()

func (e InventoryCommonStatus) MarshalJSON() ([]byte, error) {
	v, ok := InventoryCommonStatusStatusName[e]
	if !ok {
		return nil, fmt.Errorf("invalid values of InventoryCommonStatus")
	}

	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(v)
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (e *InventoryCommonStatus) UnmarshalJSON(data []byte) error {
	data = bytes.Trim(data, "\"")
	v, ok := InventoryCommonStatusStatusValue[string(data)]
	if !ok {
		return fmt.Errorf("enum '%s' is not register, must be one of: %v", data, e.EnumDescriptions())
	}

	*e = v

	return nil
}

func (*InventoryCommonStatus) EnumDescriptions() []string {
	vals := []string{}

	for _, name := range InventoryCommonStatusStatusName {
		vals = append(vals, name)
	}

	return vals
}
