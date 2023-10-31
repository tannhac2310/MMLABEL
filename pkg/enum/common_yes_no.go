package enum

import (
	"bytes"
	"fmt"
)

type CommonBoolean uint8

const (
	CommonBooleanYes CommonBoolean = iota + 1
	CommonBooleanNo
)

var CommonBooleanName = map[CommonBoolean]string{
	CommonBooleanYes: "yes",
	CommonBooleanNo:  "no",
}

var CommonBooleanValue = func() map[string]CommonBoolean {
	value := map[string]CommonBoolean{}
	for k, v := range CommonBooleanName {
		value[v] = k
		value[fmt.Sprintf("%v", k)] = k
	}

	return value
}()

func (e CommonBoolean) MarshalJSON() ([]byte, error) {
	v, ok := CommonBooleanName[e]
	if !ok {
		return nil, fmt.Errorf("invalid values of CommonBoolean")
	}

	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(v)
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (e *CommonBoolean) UnmarshalJSON(data []byte) error {
	data = bytes.Trim(data, "\"")
	v, ok := CommonBooleanValue[string(data)]
	if !ok {
		return fmt.Errorf("enum '%s' is not register, must be one of: %v", data, e.EnumDescriptions())
	}

	*e = v

	return nil
}

func (*CommonBoolean) EnumDescriptions() []string {
	vals := []string{}

	for _, name := range CommonBooleanName {
		vals = append(vals, name)
	}

	return vals
}
