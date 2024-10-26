package enum

import (
	"bytes"
	"fmt"
)

type MasterDataType string

const (
	MasterDataTypeMBomKhuonIn       MasterDataType = "m-bom-khuon-in"
	MasterDataTypeMBomKhuonPhim     MasterDataType = "m-bom-khuon-phim"
	MasterDataTypeMBomKhuonBe       MasterDataType = "m-bom-khuon-be"
	MasterDataTypeMBomKhuonDap      MasterDataType = "m-bom-khuon-dap"
	MasterDataTypeMBomNguyenVatLieu MasterDataType = "m-bom-nguyen-vat-lieu"
)

var MasterDataTypeName = map[MasterDataType]string{
	MasterDataTypeMBomKhuonIn:       "m-bom-khuon-in",
	MasterDataTypeMBomKhuonPhim:     "m-bom-khuon-phim",
	MasterDataTypeMBomKhuonBe:       "m-bom-khuon-be",
	MasterDataTypeMBomKhuonDap:      "m-bom-khuon-dap",
	MasterDataTypeMBomNguyenVatLieu: "m-bom-nguyen-vat-lieu",
}

var MasterDataTypeValue = func() map[string]MasterDataType {
	value := map[string]MasterDataType{}
	for k, v := range MasterDataTypeName {
		value[v] = k
		value[fmt.Sprintf("%v", k)] = k
	}

	return value
}()

func (e MasterDataType) MarshalJSON() ([]byte, error) {
	v, ok := MasterDataTypeName[e]
	if !ok {
		return nil, fmt.Errorf("invalid values of MasterDataType")
	}

	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(v)
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (e *MasterDataType) UnmarshalJSON(data []byte) error {
	data = bytes.Trim(data, "\"")
	v, ok := MasterDataTypeValue[string(data)]
	if !ok {
		return fmt.Errorf("enum '%s' is not register, must be one of: %v", data, e.EnumDescriptions())
	}

	*e = v

	return nil
}

func (*MasterDataType) EnumDescriptions() []string {
	vals := []string{}

	for _, name := range MasterDataTypeName {
		vals = append(vals, name)
	}

	return vals
}
