package enum

import (
	"bytes"
	"fmt"
)

type CustomFieldType uint8

const (
	CustomFieldTypeProductionOrder CustomFieldType = iota + 1
	CustomFieldTypeProductionPlan
	CustomFieldTypeProduct
	CustomFieldTypeCustomer
)

var CustomFieldTypeName = map[CustomFieldType]string{
	CustomFieldTypeProductionOrder: "production_order",
	CustomFieldTypeProductionPlan:  "production_plan",
	CustomFieldTypeProduct:         "product",
	CustomFieldTypeCustomer:        "customer",
}

var CustomFieldTypeValue = func() map[string]CustomFieldType {
	value := map[string]CustomFieldType{}
	for k, v := range CustomFieldTypeName {
		value[v] = k
		value[fmt.Sprintf("%v", k)] = k
	}

	return value
}()

func (e CustomFieldType) MarshalJSON() ([]byte, error) {
	v, ok := CustomFieldTypeName[e]
	if !ok {
		return nil, fmt.Errorf("invalid values of CustomFieldType")
	}

	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(v)
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (e *CustomFieldType) UnmarshalJSON(data []byte) error {
	data = bytes.Trim(data, "\"")
	v, ok := CustomFieldTypeValue[string(data)]
	if !ok {
		return fmt.Errorf("enum '%s' is not register, must be one of: %v", data, e.EnumDescriptions())
	}

	*e = v

	return nil
}

func (*CustomFieldType) EnumDescriptions() []string {
	vals := []string{}

	for _, name := range CustomFieldTypeName {
		vals = append(vals, name)
	}

	return vals
}
