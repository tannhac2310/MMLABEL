package enum

import (
	"bytes"
	"fmt"
)

type PaymentMethod uint8

const (
	PaymentMethodAtHome PaymentMethod = iota + 1
	PaymentMethodBankTransfer
)

var PaymentMethodName = map[PaymentMethod]string{
	PaymentMethodAtHome:       "at_home",
	PaymentMethodBankTransfer: "transfer",
}

var PaymentMethodValue = func() map[string]PaymentMethod {
	value := map[string]PaymentMethod{}
	for k, v := range PaymentMethodName {
		value[v] = k
		value[fmt.Sprintf("%v", k)] = k
	}

	return value
}()

func (e PaymentMethod) MarshalJSON() ([]byte, error) {
	v, ok := PaymentMethodName[e]
	if !ok {
		return nil, fmt.Errorf("invalid values of PaymentMethod")
	}

	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(v)
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (e *PaymentMethod) UnmarshalJSON(data []byte) error {
	data = bytes.Trim(data, "\"")
	v, ok := PaymentMethodValue[string(data)]
	if !ok {
		return fmt.Errorf("enum '%s' is not register, must be one of: %v", data, e.EnumDescriptions())
	}

	*e = v

	return nil
}

func (*PaymentMethod) EnumDescriptions() []string {
	vals := []string{}

	for _, name := range PaymentMethodName {
		vals = append(vals, name)
	}

	return vals
}
