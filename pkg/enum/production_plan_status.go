package enum

import (
	"bytes"
	"fmt"
)

type ProductionPlanStatus uint8

const (
	ProductionPlanStatusWaiting ProductionPlanStatus = iota + 1
	ProductionPlanStatusDoing
	ProductionPlanStatusPause
	ProductionPlanStatusComplete
	ProductionPlanStatusCancel
)

var ProductionPlanStatusName = map[ProductionPlanStatus]string{
	ProductionPlanStatusWaiting:  "waiting",
	ProductionPlanStatusDoing:    "doing",
	ProductionPlanStatusPause:    "pause",
	ProductionPlanStatusComplete: "complete",
	ProductionPlanStatusCancel:   "cancel",
}

var ProductionPlanStatusValue = func() map[string]ProductionPlanStatus {
	value := map[string]ProductionPlanStatus{}
	for k, v := range ProductionPlanStatusName {
		value[v] = k
		value[fmt.Sprintf("%v", k)] = k
	}

	return value
}()

func (e ProductionPlanStatus) MarshalJSON() ([]byte, error) {
	v, ok := ProductionPlanStatusName[e]
	if !ok {
		return nil, fmt.Errorf("invalid values of ProductionPlanStatus")
	}

	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(v)
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (e *ProductionPlanStatus) UnmarshalJSON(data []byte) error {
	data = bytes.Trim(data, "\"")
	v, ok := ProductionPlanStatusValue[string(data)]
	if !ok {
		return fmt.Errorf("enum '%s' is not register, must be one of: %v", data, e.EnumDescriptions())
	}

	*e = v

	return nil
}

func (*ProductionPlanStatus) EnumDescriptions() []string {
	vals := []string{}

	for _, name := range ProductionPlanStatusName {
		vals = append(vals, name)
	}

	return vals
}
