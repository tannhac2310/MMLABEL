package enum

import (
	"bytes"
	"fmt"
)

type StageStatus uint8

const (
	StageStatusPending StageStatus = iota + 1
	StageStatusStart
	StageStatusDoing
	StageStatusPausing
	StageStatusComplete
)

var StageStatusName = map[StageStatus]string{
	StageStatusPending:  "pending",
	StageStatusStart:    "start",
	StageStatusDoing:    "doing",
	StageStatusPausing:  "pausing",
	StageStatusComplete: "complete",
}

var StageStatusValue = func() map[string]StageStatus {
	value := map[string]StageStatus{}
	for k, v := range StageStatusName {
		value[v] = k
		value[fmt.Sprintf("%v", k)] = k
	}

	return value
}()

func (e StageStatus) MarshalJSON() ([]byte, error) {
	v, ok := StageStatusName[e]
	if !ok {
		return nil, fmt.Errorf("invalid values of StageStatus")
	}

	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(v)
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (e *StageStatus) UnmarshalJSON(data []byte) error {
	data = bytes.Trim(data, "\"")
	v, ok := StageStatusValue[string(data)]
	if !ok {
		return fmt.Errorf("enum '%s' is not register, must be one of: %v", data, e.EnumDescriptions())
	}

	*e = v

	return nil
}

func (*StageStatus) EnumDescriptions() []string {
	vals := []string{}

	for _, name := range StageStatusName {
		vals = append(vals, name)
	}

	return vals
}
