package enum

import (
	"bytes"
	"fmt"
)

type CourseType uint8

const (
	CourseTypeOnline CourseType = iota + 1
	CourseTypeOffline
)

var CourseTypeName = map[CourseType]string{
	CourseTypeOnline:  "online",
	CourseTypeOffline: "offline",
}

var CourseTypeValue = func() map[string]CourseType {
	value := map[string]CourseType{}
	for k, v := range CourseTypeName {
		value[v] = k
		value[fmt.Sprintf("%v", k)] = k
	}

	return value
}()

func (e CourseType) MarshalJSON() ([]byte, error) {
	v, ok := CourseTypeName[e]
	if !ok {
		return nil, fmt.Errorf("invalid values of CourseType")
	}

	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(v)
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (e *CourseType) UnmarshalJSON(data []byte) error {
	data = bytes.Trim(data, "\"")
	v, ok := CourseTypeValue[string(data)]
	if !ok {
		return fmt.Errorf("enum '%s' is not register, must be one of: %v", data, e.EnumDescriptions())
	}

	*e = v

	return nil
}

func (*CourseType) EnumDescriptions() []string {
	vals := []string{}

	for _, name := range CourseTypeName {
		vals = append(vals, name)
	}

	return vals
}
