package enum

import (
	"bytes"
	"fmt"
)

type NotificationType uint8

const (
	NotificationTypeAlert NotificationType = iota + 1
	NotificationTypeAlertIoTDevice
)

var NotificationTypeName = map[NotificationType]string{
	NotificationTypeAlert:          "alert",
	NotificationTypeAlertIoTDevice: "alert_iotdevice",
}

var NotificationTypeValue = func() map[string]NotificationType {
	value := map[string]NotificationType{}
	for k, v := range NotificationTypeName {
		value[v] = k
		value[fmt.Sprintf("%v", k)] = k
	}

	return value
}()

func (e NotificationType) MarshalJSON() ([]byte, error) {
	v, ok := NotificationTypeName[e]
	if !ok {
		return []byte("\"\""), nil
	}

	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(v)
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (e *NotificationType) UnmarshalJSON(data []byte) error {
	data = bytes.Trim(data, "\"")
	v, ok := NotificationTypeValue[string(data)]
	if !ok {
		return fmt.Errorf("enum '%s' is not register, must be one of: %v", data, e.EnumDescriptions())
	}

	*e = v

	return nil
}

func (*NotificationType) EnumDescriptions() []string {
	vals := []string{}

	for _, name := range NotificationTypeName {
		vals = append(vals, name)
	}

	return vals
}
