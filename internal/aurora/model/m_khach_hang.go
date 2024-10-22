package model

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/genproto/googleapis/type/decimal"
)

const (
	MKhachHangFieldID        = "id"
	MKhachHangFieldData      = "data"
	MKhachHangFieldCreatedBy = "created_by"
	MKhachHangFieldCreatedAt = "created_at"
	MKhachHangFieldUpdatedBy = "updated_by"
	MKhachHangFieldUpdatedAt = "updated_at"
	MKhachHangFieldDeletedAt = "deleted_at"
)

type MKhachHang struct {
	ID        string         `db:"id"`
	Data      MKhachHangData `db:"data"`
	CreatedBy string         `db:"created_by"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedBy string         `db:"updated_by"`
	UpdatedAt time.Time      `db:"updated_at"`
	DeletedAt sql.NullTime   `db:"deleted_at"`
}

func (rcv *MKhachHang) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		MKhachHangFieldID,
		MKhachHangFieldData,
		MKhachHangFieldCreatedBy,
		MKhachHangFieldCreatedAt,
		MKhachHangFieldUpdatedBy,
		MKhachHangFieldUpdatedAt,
		MKhachHangFieldDeletedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.Data,
		&rcv.CreatedBy,
		&rcv.CreatedAt,
		&rcv.UpdatedBy,
		&rcv.UpdatedAt,
		&rcv.DeletedAt,
	}

	return
}

func (*MKhachHang) TableName() string {
	return "m_khach_hang"
}

type MKhachHangData struct {
	MaSo             string                     `json:"ma_so,omitempty"`
	Ten              string                     `json:"ten,omitempty"`
	NganhNghe        string                     `json:"nganh_nghe,omitempty"`
	MaSoThue         string                     `json:"ma_so_thue,omitempty"`
	QuocGia          string                     `json:"quoc_gia,omitempty"`
	TinhThanh        string                     `json:"tinh_thanh,omitempty"`
	DiaChi           string                     `json:"dia_chi,omitempty"`
	SdtCongTy        string                     `json:"sdt_cong_ty,omitempty"`
	Fax              string                     `json:"fax,omitempty"`
	Website          string                     `json:"website,omitempty"`
	EmailCongTy      string                     `json:"email_cong_ty,omitempty"`
	TongDoanhThu     *decimal.Decimal           `json:"tong_doanh_thu,omitempty"`
	GhiChu           string                     `json:"ghi_chu,omitempty"`
	DanhSachLienHe   []*MKhachHangData_LienHe   `json:"danh_sach_lien_he,omitempty"`
	DanhSachGiaoHang []*MKhachHangData_GiaoHang `json:"danh_sach_giao_hang,omitempty"`
	DanhSachFile     []string                   `json:"danh_sach_file,omitempty"`
}

type MKhachHangData_LienHe struct {
	NguoiLienHe string `json:"nguoi_lien_he,omitempty"`
	ChucVu      string `json:"chuc_vu,omitempty"`
	Sdt         string `json:"sdt,omitempty"`
	Email       string `json:"email,omitempty"`
}

type MKhachHangData_GiaoHang struct {
	DiaChi string `json:"dia_chi,omitempty"`
	Sdt    string `json:"sdt,omitempty"`
}

func (e MKhachHangData) Value() (driver.Value, error) {
	return json.Marshal(e)
}

func (e *MKhachHangData) Scan(v any) error {
	switch vv := v.(type) {
	case []byte:
		return json.Unmarshal(vv, e)
	case string:
		return json.Unmarshal([]byte(vv), e)
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
}
