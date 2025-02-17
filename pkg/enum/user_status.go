package enum

import (
	"bytes"
	"fmt"
)

type UserStatus uint8

const (
	UserStatusUnknown UserStatus = 0
	UserStatusActive  UserStatus = 1
	UserStatusBan     UserStatus = 2
)

var UserStatusName = map[UserStatus]string{
	UserStatusUnknown: "unknown",
	UserStatusActive:  "active",
	UserStatusBan:     "ban",
}

var UserStatusValue = func() map[string]UserStatus {
	value := map[string]UserStatus{}
	for k, v := range UserStatusName {
		value[v] = k
		value[fmt.Sprintf("%v", k)] = k
	}

	return value
}()

func (e UserStatus) MarshalJSON() ([]byte, error) {
	v, ok := UserStatusName[e]
	if !ok {
		return nil, fmt.Errorf("invalid values of UserStatus")
	}

	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(v)
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (e *UserStatus) UnmarshalJSON(data []byte) error {
	data = bytes.Trim(data, "\"")
	v, ok := UserStatusValue[string(data)]
	if !ok {
		return fmt.Errorf("enum '%s' is not register, must be one of: %v", data, e.EnumDescriptions())
	}

	*e = v

	return nil
}

func (*UserStatus) EnumDescriptions() []string {
	vals := []string{}

	for _, name := range UserStatusName {
		vals = append(vals, name)
	}

	return vals
}
