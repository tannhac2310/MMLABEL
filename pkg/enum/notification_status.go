package enum

import (
	"bytes"
	"fmt"
)

type NotificationStatus uint8

const (
	NotificationStatusNew NotificationStatus = iota + 1
	NotificationStatusSeen
	NotificationStatusRead
)

var NotificationStatusName = map[NotificationStatus]string{
	NotificationStatusNew:  "new",
	NotificationStatusSeen: "seen",
	NotificationStatusRead: "read",
}

var NotificationStatusValue = func() map[string]NotificationStatus {
	value := map[string]NotificationStatus{}
	for k, v := range NotificationStatusName {
		value[v] = k
		value[fmt.Sprintf("%v", k)] = k
	}

	return value
}()

func (e NotificationStatus) MarshalJSON() ([]byte, error) {
	v, ok := NotificationStatusName[e]
	if !ok {
		return []byte("\"\""), nil
	}

	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(v)
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (e *NotificationStatus) UnmarshalJSON(data []byte) error {
	data = bytes.Trim(data, "\"")
	v, ok := NotificationStatusValue[string(data)]
	if !ok {
		return fmt.Errorf("enum '%s' is not register, must be one of: %v", data, e.EnumDescriptions())
	}

	*e = v

	return nil
}

func (*NotificationStatus) EnumDescriptions() []string {
	vals := []string{}

	for _, name := range NotificationStatusName {
		vals = append(vals, name)
	}

	return vals
}
