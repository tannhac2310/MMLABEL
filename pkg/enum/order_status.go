package enum

import (
	"bytes"
	"fmt"
)

type OrderStatus string

// chờ sản xuất/ sản xuất/ hoàn thành sx/ giao hàng
const (
	OrderStatusChoSX     OrderStatus = "cho_san_xuat"
	OrderStatusSanXuat   OrderStatus = "san_xuat"
	OrderStatusHoanThanh OrderStatus = "hoan_thanh"
	OrderStatusGiaoHang  OrderStatus = "giao_hang"
)

var OrderStatusName = map[OrderStatus]string{
	OrderStatusChoSX:     "cho_san_xuat",
	OrderStatusSanXuat:   "san_xuat",
	OrderStatusHoanThanh: "hoan_thanh",
	OrderStatusGiaoHang:  "giao_hang",
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
