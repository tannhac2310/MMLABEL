package enum

import (
	"bytes"
	"fmt"
)

type CourseLevel uint8

const (
	CourseLevelNone  CourseLevel = iota
	CourseLevelBasic             = iota + 1
	CourseLevelAdvance
	CourseLevelExpert
)

var CourseLevelName = map[CourseLevel]string{
	CourseLevelNone:    "none",
	CourseLevelBasic:   "basic",
	CourseLevelAdvance: "advance",
	CourseLevelExpert:  "expert",
}

var CourseLevelValue = func() map[string]CourseLevel {
	value := map[string]CourseLevel{}
	for k, v := range CourseLevelName {
		value[v] = k
		value[fmt.Sprintf("%v", k)] = k
	}

	return value
}()

func (e CourseLevel) MarshalJSON() ([]byte, error) {
	v, ok := CourseLevelName[e]
	if !ok {
		return nil, fmt.Errorf("invalid values of CourseLevel")
	}

	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(v)
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (e *CourseLevel) UnmarshalJSON(data []byte) error {
	data = bytes.Trim(data, "\"")
	v, ok := CourseLevelValue[string(data)]
	if !ok {
		return fmt.Errorf("enum '%s' is not register, must be one of: %v", data, e.EnumDescriptions())
	}

	*e = v

	return nil
}

func (*CourseLevel) EnumDescriptions() []string {
	vals := []string{}

	for _, name := range CourseLevelName {
		vals = append(vals, name)
	}

	return vals
}
