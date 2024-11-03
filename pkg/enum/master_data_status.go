package enum

import (
	"bytes"
	"fmt"
)

type MasterDataStatus uint8

const (
	MasterDataStatusActive MasterDataStatus = iota + 1
	MasterDataStatusInactive
	MasterDataStatusNew
)

var MasterDataStatusName = map[MasterDataStatus]string{
	MasterDataStatusActive:   "active",
	MasterDataStatusInactive: "inactive",
	MasterDataStatusNew:      "new",
}

var MasterDataStatusValue = func() map[string]MasterDataStatus {
	value := map[string]MasterDataStatus{}
	for k, v := range MasterDataStatusName {
		value[v] = k
		value[fmt.Sprintf("%v", k)] = k
	}

	return value
}()

func (e MasterDataStatus) MarshalJSON() ([]byte, error) {
	v, ok := MasterDataStatusName[e]
	if !ok {
		return nil, fmt.Errorf("invalid values of MasterDataStatus")
	}

	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(v)
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (e *MasterDataStatus) UnmarshalJSON(data []byte) error {
	data = bytes.Trim(data, "\"")
	v, ok := MasterDataStatusValue[string(data)]
	if !ok {
		return fmt.Errorf("enum '%s' is not register, must be one of: %v", data, e.EnumDescriptions())
	}

	*e = v

	return nil
}

func (*MasterDataStatus) EnumDescriptions() []string {
	vals := []string{}

	for _, name := range MasterDataStatusName {
		vals = append(vals, name)
	}

	return vals
}
