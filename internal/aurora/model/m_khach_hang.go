package model

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
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
	MaSo        string `mapstructure:"ma_so,omitempty" json:"ma_so,omitempty"`
	Ten         string `mapstructure:"ten,omitempty" json:"ten,omitempty"`
	NganhNghe   string `mapstructure:"nganh_nghe,omitempty" json:"nganh_nghe,omitempty"`
	MaSoThue    string `mapstructure:"ma_so_thue,omitempty" json:"ma_so_thue,omitempty"`
	QuocGia     string `mapstructure:"quoc_gia,omitempty" json:"quoc_gia,omitempty"`
	TinhThanh   string `mapstructure:"tinh_thanh,omitempty" json:"tinh_thanh,omitempty"`
	DiaChi      string `mapstructure:"dia_chi,omitempty" json:"dia_chi,omitempty"`
	SdtCongTy   string `mapstructure:"sdt_cong_ty,omitempty" json:"sdt_cong_ty,omitempty"`
	Fax         string `mapstructure:"fax,omitempty" json:"fax,omitempty"`
	Website     string `mapstructure:"website,omitempty" json:"website,omitempty"`
	EmailCongTy string `mapstructure:"email_cong_ty,omitempty" json:"email_cong_ty,omitempty"`
	// TongDoanhThu respect floating point, could use decimal.Decimal but need to override mapstructure
	TongDoanhThu     string                     `mapstructure:"tong_doanh_thu,omitempty" json:"tong_doanh_thu,omitempty"`
	GhiChu           string                     `mapstructure:"ghi_chu,omitempty" json:"ghi_chu,omitempty"`
	DanhSachLienHe   []*MKhachHangData_LienHe   `mapstructure:"danh_sach_lien_he,omitempty" json:"danh_sach_lien_he,omitempty"`
	DanhSachGiaoHang []*MKhachHangData_GiaoHang `mapstructure:"danh_sach_giao_hang,omitempty" json:"danh_sach_giao_hang,omitempty"`
	DanhSachFile     []string                   `mapstructure:"danh_sach_file,omitempty" json:"danh_sach_file,omitempty"`
}

type MKhachHangData_LienHe struct {
	NguoiLienHe string `mapstructure:"nguoi_lien_he,omitempty" json:"nguoi_lien_he,omitempty"`
	ChucVu      string `mapstructure:"chuc_vu,omitempty" json:"chuc_vu,omitempty"`
	Sdt         string `mapstructure:"sdt,omitempty" json:"sdt,omitempty"`
	Email       string `mapstructure:"email,omitempty" json:"email,omitempty"`
}

type MKhachHangData_GiaoHang struct {
	DiaChi string `mapstructure:"dia_chi,omitempty" json:"dia_chi,omitempty"`
	Sdt    string `mapstructure:"sdt,omitempty" json:"sdt,omitempty"`
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
