package enum

import (
	"bytes"
	"fmt"
)

type SMSProvider uint8

const (
	SMSProviderBrandSMS SMSProvider = iota + 1
)

var SMSProviderName = map[SMSProvider]string{
	SMSProviderBrandSMS: "brandsms",
}

var SMSProviderValue = func() map[string]SMSProvider {
	value := map[string]SMSProvider{}
	for k, v := range SMSProviderName {
		value[v] = k
		value[fmt.Sprintf("%v", k)] = k
	}

	return value
}()

func (e SMSProvider) MarshalJSON() ([]byte, error) {
	v, ok := SMSProviderName[e]
	if !ok {
		return nil, fmt.Errorf("invalid values of ShrimpDiaryStatus")
	}

	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(v)
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (e *SMSProvider) UnmarshalJSON(data []byte) error {
	data = bytes.Trim(data, "\"")
	v, ok := SMSProviderValue[string(data)]
	if !ok {
		return fmt.Errorf("enum '%s' is not register, must be one of: %v", data, e.EnumDescriptions())
	}

	*e = v

	return nil
}

func (*SMSProvider) EnumDescriptions() []string {
	vals := []string{}

	for _, name := range SMSProviderName {
		vals = append(vals, name)
	}

	return vals
}
