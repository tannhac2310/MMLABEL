package enum

import (
	"bytes"
	"fmt"
)

type ProductionOrderStageDeviceStatus uint8

const (
	ProductionOrderStageDeviceStatusNone ProductionOrderStageDeviceStatus = iota + 1
	ProductionOrderStageDeviceStatusStart
	ProductionOrderStageDeviceStatusPause
	ProductionOrderStageDeviceStatusComplete
	ProductionOrderStageDeviceStatusFailed
	ProductionOrderStageDeviceStatusTestProduce
	ProductionOrderStageDeviceStatusCompleteTestProduce
)

var ProductionOrderStageDeviceStatusName = map[ProductionOrderStageDeviceStatus]string{
	ProductionOrderStageDeviceStatusNone:                "",                      // chua bat dau
	ProductionOrderStageDeviceStatusStart:               "start",                 // bat dau sx
	ProductionOrderStageDeviceStatusPause:               "pause",                 // tam dung sx
	ProductionOrderStageDeviceStatusComplete:            "completed",             // hoan thanh sx
	ProductionOrderStageDeviceStatusFailed:              "failed",                // loi
	ProductionOrderStageDeviceStatusTestProduce:         "test_produce",          // test_produce
	ProductionOrderStageDeviceStatusCompleteTestProduce: "complete_test_produce", // complete_test_produce
}

var ProductionOrderStageDeviceStatusValue = func() map[string]ProductionOrderStageDeviceStatus {
	value := map[string]ProductionOrderStageDeviceStatus{}
	for k, v := range ProductionOrderStageDeviceStatusName {
		value[v] = k
		value[fmt.Sprintf("%v", k)] = k
	}

	return value
}()

func (e ProductionOrderStageDeviceStatus) MarshalJSON() ([]byte, error) {
	v, ok := ProductionOrderStageDeviceStatusName[e]
	if !ok {
		return nil, fmt.Errorf("invalid values of ProductionOrderStageDeviceStatus")
	}

	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(v)
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (e *ProductionOrderStageDeviceStatus) UnmarshalJSON(data []byte) error {
	data = bytes.Trim(data, "\"")
	v, ok := ProductionOrderStageDeviceStatusValue[string(data)]
	if !ok {
		return fmt.Errorf("enum '%s' is not register, must be one of: %v", data, e.EnumDescriptions())
	}

	*e = v

	return nil
}

func (*ProductionOrderStageDeviceStatus) EnumDescriptions() []string {
	vals := []string{}

	for _, name := range ProductionOrderStageDeviceStatusName {
		vals = append(vals, name)
	}

	return vals
}
