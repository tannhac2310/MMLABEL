package enum

import (
	"bytes"
	"fmt"
)

type CustomerStatus uint8

const (
	CustomerStatusActivate CustomerStatus = iota + 1
	CustomerStatusDeactivate
)

var CustomerStatusName = map[CustomerStatus]string{
	CustomerStatusActivate:   "activate",
	CustomerStatusDeactivate: "deactivate",
}

var CustomerStatusValue = func() map[string]CustomerStatus {
	value := map[string]CustomerStatus{}
	for k, v := range CustomerStatusName {
		value[v] = k
		value[fmt.Sprintf("%v", k)] = k
	}

	return value
}()

func (e CustomerStatus) MarshalJSON() ([]byte, error) {
	v, ok := CustomerStatusName[e]
	if !ok {
		return nil, fmt.Errorf("invalid values of CustomerStatus")
	}

	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(v)
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (e *CustomerStatus) UnmarshalJSON(data []byte) error {
	data = bytes.Trim(data, "\"")
	v, ok := CustomerStatusValue[string(data)]
	if !ok {
		return fmt.Errorf("enum '%s' is not register, must be one of: %v", data, e.EnumDescriptions())
	}

	*e = v

	return nil
}

func (*CustomerStatus) EnumDescriptions() []string {
	vals := []string{}

	for _, name := range CustomerStatusName {
		vals = append(vals, name)
	}

	return vals
}
