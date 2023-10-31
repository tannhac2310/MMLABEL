package enum

import (
	"bytes"
	"fmt"
)

type StageStudentStatus uint8

const (
	StageStudentStatusLearning StageStudentStatus = iota + 1
	StageStudentStatusFinish
	StageStudentStatusReserve
)

var StageStudentStatusName = map[StageStudentStatus]string{
	StageStudentStatusLearning: "active",
	StageStudentStatusFinish:   "disable",
	StageStudentStatusReserve:  "reserve",
}

var StageStudentStatusValue = func() map[string]StageStudentStatus {
	value := map[string]StageStudentStatus{}
	for k, v := range StageStudentStatusName {
		value[v] = k
		value[fmt.Sprintf("%v", k)] = k
	}

	return value
}()

func (e StageStudentStatus) MarshalJSON() ([]byte, error) {
	v, ok := StageStudentStatusName[e]
	if !ok {
		return nil, fmt.Errorf("invalid values of StageStudentStatus")
	}

	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(v)
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (e *StageStudentStatus) UnmarshalJSON(data []byte) error {
	data = bytes.Trim(data, "\"")
	v, ok := StageStudentStatusValue[string(data)]
	if !ok {
		return fmt.Errorf("enum '%s' is not register, must be one of: %v", data, e.EnumDescriptions())
	}

	*e = v

	return nil
}

func (*StageStudentStatus) EnumDescriptions() []string {
	vals := []string{}

	for _, name := range StageStudentStatusName {
		vals = append(vals, name)
	}

	return vals
}
