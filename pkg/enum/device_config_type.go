package enum

import (
	"bytes"
	"fmt"
)

type DeviceConfigType string

const (
	DeviceConfigTypeMayIn   DeviceConfigType = "printer"
	DeviceConfigTypeMayKhac DeviceConfigType = "may_khac"
)

var DeviceConfigTypeName = map[DeviceConfigType]string{
	DeviceConfigTypeMayIn:   "printer",
	DeviceConfigTypeMayKhac: "may_khac",
}

var DeviceConfigTypeValue = func() map[string]DeviceConfigType {
	value := map[string]DeviceConfigType{}
	for k, v := range DeviceConfigTypeName {
		value[v] = k
		value[fmt.Sprintf("%v", k)] = k
	}

	return value
}()

func (e DeviceConfigType) MarshalJSON() ([]byte, error) {
	v, ok := DeviceConfigTypeName[e]
	if !ok {
		return nil, fmt.Errorf("invalid values of DeviceConfigType")
	}

	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(v)
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (e *DeviceConfigType) UnmarshalJSON(data []byte) error {
	data = bytes.Trim(data, "\"")
	v, ok := DeviceConfigTypeValue[string(data)]
	if !ok {
		return fmt.Errorf("enum '%s' is not register, must be one of: %v", data, e.EnumDescriptions())
	}

	*e = v

	return nil
}

func (*DeviceConfigType) EnumDescriptions() []string {
	vals := []string{}

	for _, name := range DeviceConfigTypeName {
		vals = append(vals, name)
	}

	return vals
}
