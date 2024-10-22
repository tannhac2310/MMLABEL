package enum

import (
	"bytes"
	"fmt"
)

type MasterDataType uint8

const (
	MasterDataType_KhachHang MasterDataType = iota
	MasterDataType_KhungIn
	MasterDataType_KhuonBe
	MasterDataType_KhuonDap
	MasterDataType_NguyenVatLieu
	MasterDataType_Phim
	MasterDataType_ThongSoMayIn
	MasterDataType_ThongSoMayKhac
)

var MasterDataTypeName = map[MasterDataType]string{
	MasterDataType_KhachHang:      "m_khach_hang",
	MasterDataType_KhungIn:        "m_khung_in",
	MasterDataType_KhuonBe:        "m_khuon_be",
	MasterDataType_KhuonDap:       "m_khuong_dap",
	MasterDataType_NguyenVatLieu:  "m_nguyen_vat_lieu",
	MasterDataType_Phim:           "m_phim",
	MasterDataType_ThongSoMayIn:   "m_thong_so_may_in",
	MasterDataType_ThongSoMayKhac: "m_thong_so_may_khac",
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
