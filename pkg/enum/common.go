package enum

import (
	"bytes"
	"fmt"
)

type CommonStatus uint8

const (
	CommonStatusActive CommonStatus = iota + 1
	CommonStatusDisable
)

var CommonStatusName = map[CommonStatus]string{
	CommonStatusActive:  "active",
	CommonStatusDisable: "disable",
}

var CommonStatusValue = func() map[string]CommonStatus {
	value := map[string]CommonStatus{}
	for k, v := range CommonStatusName {
		value[v] = k
		value[fmt.Sprintf("%v", k)] = k
	}

	return value
}()

func (e CommonStatus) MarshalJSON() ([]byte, error) {
	v, ok := CommonStatusName[e]
	if !ok {
		return nil, fmt.Errorf("invalid values of CommonStatus")
	}

	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(v)
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (e *CommonStatus) UnmarshalJSON(data []byte) error {
	data = bytes.Trim(data, "\"")
	v, ok := CommonStatusValue[string(data)]
	if !ok {
		return fmt.Errorf("enum '%s' is not register, must be one of: %v", data, e.EnumDescriptions())
	}

	*e = v

	return nil
}

func (*CommonStatus) EnumDescriptions() []string {
	vals := []string{}

	for _, name := range CommonStatusName {
		vals = append(vals, name)
	}

	return vals
}
