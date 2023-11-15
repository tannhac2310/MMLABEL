package enum

import (
	"bytes"
	"fmt"
)

type ProductionOrderStatus uint8

const (
	ProductionOrderStatusWaiting ProductionOrderStatus = iota + 1
	ProductionOrderStatusDoing
	ProductionOrderStatusPause
	ProductionOrderStatusComplete
	ProductionOrderStatusCancel
)

var ProductionOrderStatusName = map[ProductionOrderStatus]string{
	ProductionOrderStatusWaiting:  "waiting",
	ProductionOrderStatusDoing:    "doing",
	ProductionOrderStatusPause:    "pause",
	ProductionOrderStatusComplete: "complete",
	ProductionOrderStatusCancel:   "cancel",
}

var ProductionOrderStatusValue = func() map[string]ProductionOrderStatus {
	value := map[string]ProductionOrderStatus{}
	for k, v := range ProductionOrderStatusName {
		value[v] = k
		value[fmt.Sprintf("%v", k)] = k
	}

	return value
}()

func (e ProductionOrderStatus) MarshalJSON() ([]byte, error) {
	v, ok := ProductionOrderStatusName[e]
	if !ok {
		return nil, fmt.Errorf("invalid values of ProductionOrderStatus")
	}

	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(v)
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (e *ProductionOrderStatus) UnmarshalJSON(data []byte) error {
	data = bytes.Trim(data, "\"")
	v, ok := ProductionOrderStatusValue[string(data)]
	if !ok {
		return fmt.Errorf("enum '%s' is not register, must be one of: %v", data, e.EnumDescriptions())
	}

	*e = v

	return nil
}

func (*ProductionOrderStatus) EnumDescriptions() []string {
	vals := []string{}

	for _, name := range ProductionOrderStatusName {
		vals = append(vals, name)
	}

	return vals
}
