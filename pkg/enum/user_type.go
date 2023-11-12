package enum

import (
	"bytes"
	"fmt"
)

type UserType uint8

const (
	UserTypeCustomer UserType = iota + 1
	UserTypeEmployee
)

var UserTypeName = map[UserType]string{
	UserTypeCustomer: "customer",
	UserTypeEmployee: "employee",
}

var UserTypeValue = func() map[string]UserType {
	value := map[string]UserType{}
	for k, v := range UserTypeName {
		value[v] = k
		value[fmt.Sprintf("%v", k)] = k
	}

	return value
}()

func (e UserType) MarshalJSON() ([]byte, error) {
	v, ok := UserTypeName[e]
	if !ok {
		return nil, fmt.Errorf("invalid values of UserType")
	}

	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(v)
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (e *UserType) UnmarshalJSON(data []byte) error {
	data = bytes.Trim(data, "\"")
	v, ok := UserTypeValue[string(data)]
	if !ok {
		return fmt.Errorf("enum '%s' is not register, must be one of: %v", data, e.EnumDescriptions())
	}

	*e = v

	return nil
}

func (*UserType) EnumDescriptions() []string {
	vals := []string{}

	for _, name := range UserTypeName {
		vals = append(vals, name)
	}

	return vals
}
