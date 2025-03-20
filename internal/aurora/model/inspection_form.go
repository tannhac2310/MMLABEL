package model

import (
	"time"
)

const (
	InspectionFormFieldID                  = "id"
	InspectionFormFieldProductionOrderID   = "production_order_id"
	InspectionFormFieldInspectionDate      = "inspection_date"
	InspectionFormFieldInspectorName       = "inspector_name"
	InspectionFormFieldQuantity            = "quantity"
	InspectionFormFieldMaSanPham           = "ma_san_pham"
	InspectionFormFieldTenSanPham          = "ten_san_pham"
	InspectionFormFieldSoLuongHopDong      = "so_luong_hop_dong"
	InspectionFormFieldSoLuongIn           = "so_luong_in"
	InspectionFormFieldMaDonDatHang        = "ma_don_dat_hang"
	InspectionFormFieldNguoiKiemTra        = "nguoi_kiem_tra"
	InspectionFormFieldNguoiPheDuyet       = "nguoi_phe_duyet"
	InspectionFormFieldSoLuongThanhPhamDat = "so_luong_thanh_pham_dat"
	InspectionFormFieldNote                = "note"
	InspectionFormFieldCreatedBy           = "created_by"
	InspectionFormFieldUpdatedBy           = "updated_by"
	InspectionFormFieldCreatedAt           = "created_at"
	InspectionFormFieldUpdatedAt           = "updated_at"
)

type InspectionForm struct {
	ID                  string    `db:"id"`
	ProductionOrderID   string    `db:"production_order_id"`
	InspectionDate      time.Time `db:"inspection_date"`
	InspectorName       string    `db:"inspector_name"`
	Quantity            int64     `db:"quantity"`
	MaSanPham           string    `db:"ma_san_pham"`
	TenSanPham          string    `db:"ten_san_pham"`
	SoLuongHopDong      int64     `db:"so_luong_hop_dong"`
	SoLuongIn           int64     `db:"so_luong_in"`
	MaDonDatHang        string    `db:"ma_don_dat_hang"`
	NguoiKiemTra        string    `db:"nguoi_kiem_tra"`
	NguoiPheDuyet       string    `db:"nguoi_phe_duyet"`
	SoLuongThanhPhamDat int64     `db:"so_luong_thanh_pham_dat"`
	Note                string    `db:"note"`
	CreatedBy           string    `db:"created_by"`
	UpdatedBy           string    `db:"updated_by"`
	CreatedAt           time.Time `db:"created_at"`
	UpdatedAt           time.Time `db:"updated_at"`
}

func (rcv *InspectionForm) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		InspectionFormFieldID,
		InspectionFormFieldProductionOrderID,
		InspectionFormFieldInspectionDate,
		InspectionFormFieldInspectorName,
		InspectionFormFieldQuantity,
		InspectionFormFieldMaSanPham,
		InspectionFormFieldTenSanPham,
		InspectionFormFieldSoLuongHopDong,
		InspectionFormFieldSoLuongIn,
		InspectionFormFieldMaDonDatHang,
		InspectionFormFieldNguoiKiemTra,
		InspectionFormFieldNguoiPheDuyet,
		InspectionFormFieldSoLuongThanhPhamDat,
		InspectionFormFieldNote,
		InspectionFormFieldCreatedBy,
		InspectionFormFieldUpdatedBy,
		InspectionFormFieldCreatedAt,
		InspectionFormFieldUpdatedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.ProductionOrderID,
		&rcv.InspectionDate,
		&rcv.InspectorName,
		&rcv.Quantity,
		&rcv.MaSanPham,
		&rcv.TenSanPham,
		&rcv.SoLuongHopDong,
		&rcv.SoLuongIn,
		&rcv.MaDonDatHang,
		&rcv.NguoiKiemTra,
		&rcv.NguoiPheDuyet,
		&rcv.SoLuongThanhPhamDat,
		&rcv.Note,
		&rcv.CreatedBy,
		&rcv.UpdatedBy,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
	}

	return
}

func (*InspectionForm) TableName() string {
	return "inspection_forms"
}
