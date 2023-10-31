package enum

import (
	"bytes"
	"fmt"
)

type PermissionEntity uint8

const (
	PermissionEntityOa PermissionEntity = iota + 1
	PermissionEntityCourse
)

var PermissionEntityName = map[PermissionEntity]string{
	PermissionEntityOa:     "official_account",
	PermissionEntityCourse: "course",
}

var PermissionEntityValue = func() map[string]PermissionEntity {
	value := map[string]PermissionEntity{}
	for k, v := range PermissionEntityName {
		value[v] = k
		value[fmt.Sprintf("%v", k)] = k
	}

	return value
}()

func (e PermissionEntity) MarshalJSON() ([]byte, error) {
	v, ok := PermissionEntityName[e]
	if !ok {
		return nil, fmt.Errorf("invalid values of PermissionEntity")
	}

	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(v)
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (e *PermissionEntity) UnmarshalJSON(data []byte) error {
	data = bytes.Trim(data, "\"")
	v, ok := PermissionEntityValue[string(data)]
	if !ok {
		return fmt.Errorf("enum '%s' is not register, must be one of: %v", data, e.EnumDescriptions())
	}

	*e = v

	return nil
}

func (*PermissionEntity) EnumDescriptions() []string {
	vals := []string{}

	for _, name := range PermissionEntityName {
		vals = append(vals, name)
	}

	return vals
}
