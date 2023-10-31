package enum

import (
	"bytes"
	"fmt"
)

type CrmContactSource uint8

const (
	CrmContactSourceInternal CrmContactSource = iota + 1
	CrmContactSourceOpenChannelZalo
	CrmContactSourceOpenChannelFacebook
	CrmContactSourceOpenChannelTelegram
)

var CrmContactSourceName = map[CrmContactSource]string{
	CrmContactSourceInternal:            "internal",
	CrmContactSourceOpenChannelZalo:     "zalo",
	CrmContactSourceOpenChannelFacebook: "facebook",
	CrmContactSourceOpenChannelTelegram: "telegram",
}

var CrmContactSourceValue = func() map[string]CrmContactSource {
	value := map[string]CrmContactSource{}
	for k, v := range CrmContactSourceName {
		value[v] = k
		value[fmt.Sprintf("%v", k)] = k
	}

	return value
}()

func (e CrmContactSource) MarshalJSON() ([]byte, error) {
	v, ok := CrmContactSourceName[e]
	if !ok {
		return nil, fmt.Errorf("invalid values of CrmContactSource")
	}

	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(v)
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (e *CrmContactSource) UnmarshalJSON(data []byte) error {
	data = bytes.Trim(data, "\"")
	v, ok := CrmContactSourceValue[string(data)]
	if !ok {
		return fmt.Errorf("enum '%s' is not register, must be one of: %v", data, e.EnumDescriptions())
	}

	*e = v

	return nil
}

func (*CrmContactSource) EnumDescriptions() []string {
	vals := []string{}

	for _, name := range CrmContactSourceName {
		vals = append(vals, name)
	}

	return vals
}
