package enum

import (
	"bytes"
	"fmt"
)

type StageCalendar uint8

const (
	StageCalendarMonday StageCalendar = iota + 1
	StageCalendarTuesday
	StageCalendarWednesday
	StageCalendarThursday
	StageCalendarFriday
	StageCalendarSaturday
	StageCalendarSunday
)

var StageCalendarName = map[StageCalendar]string{
	StageCalendarMonday:    "monday",
	StageCalendarTuesday:   "tuesday",
	StageCalendarWednesday: "wednesday",
	StageCalendarThursday:  "thursday",
	StageCalendarFriday:    "friday",
	StageCalendarSaturday:  "saturday",
	StageCalendarSunday:    "sunday",
}

var StageCalendarValue = func() map[string]StageCalendar {
	value := map[string]StageCalendar{}
	for k, v := range StageCalendarName {
		value[v] = k
		value[fmt.Sprintf("%v", k)] = k
	}

	return value
}()

func (e StageCalendar) MarshalJSON() ([]byte, error) {
	v, ok := StageCalendarName[e]
	if !ok {
		return nil, fmt.Errorf("invalid values of StageCalendar")
	}

	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(v)
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (e *StageCalendar) UnmarshalJSON(data []byte) error {
	data = bytes.Trim(data, "\"")
	v, ok := StageCalendarValue[string(data)]
	if !ok {
		return fmt.Errorf("enum '%s' is not register, must be one of: %v", data, e.EnumDescriptions())
	}

	*e = v

	return nil
}

func (*StageCalendar) EnumDescriptions() []string {
	vals := []string{}

	for _, name := range StageCalendarName {
		vals = append(vals, name)
	}

	return vals
}
