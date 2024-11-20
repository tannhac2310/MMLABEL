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
	ProductionOrderStageDeviceStatusSetup
)

// Màn hình sẽ có 4 nút. 1. "Bắt đầu SX Thử"/"Kết thúc SX Thử", 2. "Bắt đầu SX"/"Hoàn thành SX", 3."Ngưng chuyển PO", 4."Tạm dừng"/"Sản xuất lại"
var ProductionOrderStageDeviceStatusName = map[ProductionOrderStageDeviceStatus]string{
	ProductionOrderStageDeviceStatusNone:                "",                      // none
	ProductionOrderStageDeviceStatusStart:               "start",                 // Bắt đầu sản xuất
	ProductionOrderStageDeviceStatusPause:               "pause",                 // Ngưng chuyển PO
	ProductionOrderStageDeviceStatusComplete:            "completed",             // Hoàn thành SX
	ProductionOrderStageDeviceStatusFailed:              "failed",                // Tạm dừng"/"Sản xuất lại
	ProductionOrderStageDeviceStatusTestProduce:         "test_produce",          // Bắt đầu SX Thử
	ProductionOrderStageDeviceStatusCompleteTestProduce: "complete_test_produce", // Kết thúc SX Thử
	ProductionOrderStageDeviceStatusSetup:               "setup",                 // Setup máy
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
