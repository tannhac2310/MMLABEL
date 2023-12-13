package enum

import (
	"bytes"
	"fmt"
)

type PermissionEntityType uint8

const (
	PermissionEntityTypeStage PermissionEntityType = iota + 1
	PermissionEntityTypeDevice
	PermissionEntityTypeScreen
)

var PermissionEntityTypeName = map[PermissionEntityType]string{
	PermissionEntityTypeStage:  "stage",
	PermissionEntityTypeDevice: "device",
	PermissionEntityTypeScreen: "screen",
}

var PermissionEntityTypeValue = func() map[string]PermissionEntityType {
	value := map[string]PermissionEntityType{}
	for k, v := range PermissionEntityTypeName {
		value[v] = k
		value[fmt.Sprintf("%v", k)] = k
	}

	return value
}()

func (e PermissionEntityType) MarshalJSON() ([]byte, error) {
	v, ok := PermissionEntityTypeName[e]
	if !ok {
		return nil, fmt.Errorf("invalid values of PermissionEntityType")
	}

	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(v)
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (e *PermissionEntityType) UnmarshalJSON(data []byte) error {
	data = bytes.Trim(data, "\"")
	v, ok := PermissionEntityTypeValue[string(data)]
	if !ok {
		return fmt.Errorf("enum '%s' is not register, must be one of: %v", data, e.EnumDescriptions())
	}

	*e = v

	return nil
}

func (*PermissionEntityType) EnumDescriptions() []string {
	vals := []string{}

	for _, name := range PermissionEntityTypeName {
		vals = append(vals, name)
	}

	return vals
}
