package model

import (
	"database/sql"
	"time"
)

const (
	InspectionFormFieldID                  = "id"
	InspectionFormFieldProductionOrderID   = "production_order_id"
	InspectionFormFieldInspectionDate      = "inspection_date"
	InspectionFormFieldInspectorName       = "inspector_name"
	InspectionFormFieldQuantity            = "quantity"
	InspectionFormFieldProductID           = "product_id"
	InspectionFormFieldSoLuongHopDong      = "so_luong_hop_dong"
	InspectionFormFieldSoLuongIn           = "so_luong_in"
	InspectionFormFieldNguoiKiemTra        = "nguoi_kiem_tra"
	InspectionFormFieldNguoiPheDuyet       = "nguoi_phe_duyet"
	InspectionFormFieldSoLuongThanhPhamDat = "so_luong_thanh_pham_dat"
	InspectionFormFieldNote                = "note"
	InspectionFormFieldCreatedBy           = "created_by"
	InspectionFormFieldUpdatedBy           = "updated_by"
	InspectionFormFieldCreatedAt           = "created_at"
	InspectionFormFieldUpdatedAt           = "updated_at"
	InspectionFormFieldDeletedAt           = "deleted_at"
)

type InspectionForm struct {
	ID                  string       `db:"id"`
	ProductionOrderID   string       `db:"production_order_id"`
	InspectionDate      time.Time    `db:"inspection_date"`
	InspectorName       string       `db:"inspector_name"`
	Quantity            int64        `db:"quantity"`
	ProductID           string       `db:"product_id"`
	SoLuongHopDong      int64        `db:"so_luong_hop_dong"`
	SoLuongIn           int64        `db:"so_luong_in"`
	NguoiKiemTra        string       `db:"nguoi_kiem_tra"`
	NguoiPheDuyet       string       `db:"nguoi_phe_duyet"`
	SoLuongThanhPhamDat int64        `db:"so_luong_thanh_pham_dat"`
	Note                string       `db:"note"`
	CreatedBy           string       `db:"created_by"`
	UpdatedBy           string       `db:"updated_by"`
	CreatedAt           time.Time    `db:"created_at"`
	UpdatedAt           time.Time    `db:"updated_at"`
	DeletedAt           sql.NullTime `db:"deleted_at"`
}

func (rcv *InspectionForm) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		InspectionFormFieldID,
		InspectionFormFieldProductionOrderID,
		InspectionFormFieldInspectionDate,
		InspectionFormFieldInspectorName,
		InspectionFormFieldQuantity,
		InspectionFormFieldProductID,
		InspectionFormFieldSoLuongHopDong,
		InspectionFormFieldSoLuongIn,
		InspectionFormFieldNguoiKiemTra,
		InspectionFormFieldNguoiPheDuyet,
		InspectionFormFieldSoLuongThanhPhamDat,
		InspectionFormFieldNote,
		InspectionFormFieldCreatedBy,
		InspectionFormFieldUpdatedBy,
		InspectionFormFieldCreatedAt,
		InspectionFormFieldUpdatedAt,
		InspectionFormFieldDeletedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.ProductionOrderID,
		&rcv.InspectionDate,
		&rcv.InspectorName,
		&rcv.Quantity,
		&rcv.ProductID,
		&rcv.SoLuongHopDong,
		&rcv.SoLuongIn,
		&rcv.NguoiKiemTra,
		&rcv.NguoiPheDuyet,
		&rcv.SoLuongThanhPhamDat,
		&rcv.Note,
		&rcv.CreatedBy,
		&rcv.UpdatedBy,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
		&rcv.DeletedAt,
	}

	return
}

func (*InspectionForm) TableName() string {
	return "inspection_forms"
}
