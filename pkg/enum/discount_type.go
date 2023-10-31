package enum

import (
	"bytes"
	"fmt"
)

type DiscountType uint8

const (
	DiscountTypeNone DiscountType = iota
	DiscountTypePercentage
	DiscountTypeAmount
)

var DiscountTypeName = map[DiscountType]string{
	DiscountTypeNone:       "none",
	DiscountTypePercentage: "percentage",
	DiscountTypeAmount:     "amount",
}

var DiscountTypeValue = func() map[string]DiscountType {
	value := map[string]DiscountType{}
	for k, v := range DiscountTypeName {
		value[v] = k
		value[fmt.Sprintf("%v", k)] = k
	}

	return value
}()

func (e DiscountType) MarshalJSON() ([]byte, error) {
	v, ok := DiscountTypeName[e]
	if !ok {
		return nil, fmt.Errorf("invalid values of DiscountType")
	}

	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(v)
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (e *DiscountType) UnmarshalJSON(data []byte) error {
	data = bytes.Trim(data, "\"")
	v, ok := DiscountTypeValue[string(data)]
	if !ok {
		return fmt.Errorf("enum '%s' is not register, must be one of: %v", data, e.EnumDescriptions())
	}

	*e = v

	return nil
}

func (*DiscountType) EnumDescriptions() []string {
	vals := []string{}

	for _, name := range DiscountTypeName {
		vals = append(vals, name)
	}

	return vals
}
