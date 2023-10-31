package enum

import (
	"bytes"
	"fmt"
)

type Gender uint8

const (
	GenderFemale Gender = iota
	GenderMale
)

var GenderName = map[Gender]string{
	GenderFemale: "female",
	GenderMale:   "male",
}

var GenderValue = func() map[string]Gender {
	value := map[string]Gender{}
	for k, v := range GenderName {
		value[v] = k
		value[fmt.Sprintf("%v", k)] = k
	}

	return value
}()

func (e Gender) MarshalJSON() ([]byte, error) {
	v, ok := GenderName[e]
	if !ok {
		return nil, fmt.Errorf("invalid values of Gender")
	}

	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(v)
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (e *Gender) UnmarshalJSON(data []byte) error {
	data = bytes.Trim(data, "\"")
	v, ok := GenderValue[string(data)]
	if !ok {
		return fmt.Errorf("enum '%s' is not register, must be one of: %v", data, e.EnumDescriptions())
	}

	*e = v

	return nil
}

func (*Gender) EnumDescriptions() []string {
	vals := []string{}

	for _, name := range GenderName {
		vals = append(vals, name)
	}

	return vals
}
