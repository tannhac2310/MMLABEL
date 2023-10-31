package enum

import (
	"bytes"
	"fmt"
)

type OrderStatus uint8

const (
	OrderStatusWaiting OrderStatus = iota + 1
	OrderStatusAccept
	OrderStatusReject
	OrderStatusFinish
	OrderStatusCancel
)

var OrderStatusName = map[OrderStatus]string{
	OrderStatusWaiting: "waiting",
	OrderStatusAccept:  "accept",
	OrderStatusReject:  "reject",
	OrderStatusFinish:  "finish",
	OrderStatusCancel:  "cancel",
}

var OrderStatusValue = func() map[string]OrderStatus {
	value := map[string]OrderStatus{}
	for k, v := range OrderStatusName {
		value[v] = k
		value[fmt.Sprintf("%v", k)] = k
	}

	return value
}()

func (e OrderStatus) MarshalJSON() ([]byte, error) {
	v, ok := OrderStatusName[e]
	if !ok {
		return nil, fmt.Errorf("invalid values of OrderStatus")
	}

	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(v)
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (e *OrderStatus) UnmarshalJSON(data []byte) error {
	data = bytes.Trim(data, "\"")
	v, ok := OrderStatusValue[string(data)]
	if !ok {
		return fmt.Errorf("enum '%s' is not register, must be one of: %v", data, e.EnumDescriptions())
	}

	*e = v

	return nil
}

func (*OrderStatus) EnumDescriptions() []string {
	vals := []string{}

	for _, name := range OrderStatusName {
		vals = append(vals, name)
	}

	return vals
}
