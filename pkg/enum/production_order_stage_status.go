package enum

import (
	"bytes"
	"fmt"
)

type ProductionOrderStageStatus uint8

const (
	ProductionOrderStageStatusWaiting ProductionOrderStageStatus = iota + 1
	ProductionOrderStageStatusReception
	ProductionOrderStageStatusProductionStart
	ProductionOrderStageStatusProductionCompletion
	ProductionOrderStageStatusProductDelivery
)

var ProductionOrderStageStatusName = map[ProductionOrderStageStatus]string{
	ProductionOrderStageStatusWaiting:              "waiting",               // Chờ Tiếp Nhận
	ProductionOrderStageStatusReception:            "reception",             // Tiếp Nhận
	ProductionOrderStageStatusProductionStart:      "production_start",      // Bắt Đầu Sản Xuất
	ProductionOrderStageStatusProductionCompletion: "production_completion", // Hoàn Thành Sản Xuất
	ProductionOrderStageStatusProductDelivery:      "product_delivery",      // Chuyển Giao Bán Thành Phẩm
}

var ProductionOrderStageStatusValue = func() map[string]ProductionOrderStageStatus {
	value := map[string]ProductionOrderStageStatus{}
	for k, v := range ProductionOrderStageStatusName {
		value[v] = k
		value[fmt.Sprintf("%v", k)] = k
	}

	return value
}()

func (e ProductionOrderStageStatus) MarshalJSON() ([]byte, error) {
	v, ok := ProductionOrderStageStatusName[e]
	if !ok {
		return nil, fmt.Errorf("invalid values of ProductionOrderStageStatus")
	}

	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(v)
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (e *ProductionOrderStageStatus) UnmarshalJSON(data []byte) error {
	data = bytes.Trim(data, "\"")
	v, ok := ProductionOrderStageStatusValue[string(data)]
	if !ok {
		return fmt.Errorf("enum '%s' is not register, must be one of: %v", data, e.EnumDescriptions())
	}

	*e = v

	return nil
}

func (*ProductionOrderStageStatus) EnumDescriptions() []string {
	vals := []string{}

	for _, name := range ProductionOrderStageStatusName {
		vals = append(vals, name)
	}

	return vals
}
