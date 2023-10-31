package enum

import (
	"bytes"
	"fmt"
)

type AttendanceDetailStatus uint8

const (
	AttendanceDetailStatusOnTime AttendanceDetailStatus = iota + 1
	AttendanceDetailStatusAbsenceWithPermission
	AttendanceDetailStatusAbsenceWithoutPermission
	AttendanceDetailStatusLateToClass
	AttendanceDetailStatusEarlyToHome
	AttendanceDetailStatusLateToClassEarlyToHome
)

var AttendanceDetailStatusName = map[AttendanceDetailStatus]string{
	AttendanceDetailStatusOnTime:                   "on_time",
	AttendanceDetailStatusAbsenceWithPermission:    "absence_with_permission",
	AttendanceDetailStatusAbsenceWithoutPermission: "absence_without_permission",
	AttendanceDetailStatusLateToClass:              "late_to_class",
	AttendanceDetailStatusEarlyToHome:              "early_to_home",
	AttendanceDetailStatusLateToClassEarlyToHome:   "late_to_class_early_to_home",
}

var AttendanceDetailStatusValue = func() map[string]AttendanceDetailStatus {
	value := map[string]AttendanceDetailStatus{}
	for k, v := range AttendanceDetailStatusName {
		value[v] = k
		value[fmt.Sprintf("%v", k)] = k
	}

	return value
}()

func (e AttendanceDetailStatus) MarshalJSON() ([]byte, error) {
	v, ok := AttendanceDetailStatusName[e]
	if !ok {
		return nil, fmt.Errorf("invalid values of AttendanceDetailStatus")
	}

	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(v)
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (e *AttendanceDetailStatus) UnmarshalJSON(data []byte) error {
	data = bytes.Trim(data, "\"")
	v, ok := AttendanceDetailStatusValue[string(data)]
	if !ok {
		return fmt.Errorf("enum '%s' is not register, must be one of: %v", data, e.EnumDescriptions())
	}

	*e = v

	return nil
}

func (*AttendanceDetailStatus) EnumDescriptions() []string {
	vals := []string{}

	for _, name := range AttendanceDetailStatusName {
		vals = append(vals, name)
	}

	return vals
}
