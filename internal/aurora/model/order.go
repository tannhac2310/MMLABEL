package model

import (
	"database/sql"
	"time"
)

const (
	OrderFieldID                 = "id"
	OrderFieldTitle              = "title"
	OrderFieldMaDatHangMm        = "ma_dat_hang_mm"
	OrderFieldMaHopDongKhachHang = "ma_hop_dong_khach_hang"
	OrderFieldMaHopDong          = "ma_hop_dong"
	OrderFieldSaleName           = "sale_name"
	OrderFieldSaleAdminName      = "sale_admin_name"
	OrderFieldStatus             = "status"
	OrderFieldCreatedBy          = "created_by"
	OrderFieldUpdatedBy          = "updated_by"
	OrderFieldCreatedAt          = "created_at"
	OrderFieldUpdatedAt          = "updated_at"
)

type Order struct {
	ID                 string         `db:"id"`
	Title              string         `db:"title"`
	MaDatHangMm        string         `db:"ma_dat_hang_mm"`
	MaHopDongKhachHang string         `db:"ma_hop_dong_khach_hang"`
	MaHopDong          string         `db:"ma_hop_dong"`
	SaleName           sql.NullString `db:"sale_name"`
	SaleAdminName      sql.NullString `db:"sale_admin_name"`
	Status             string         `db:"status"`
	CreatedBy          string         `db:"created_by"`
	UpdatedBy          string         `db:"updated_by"`
	CreatedAt          time.Time      `db:"created_at"`
	UpdatedAt          time.Time      `db:"updated_at"`
}

func (rcv *Order) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		OrderFieldID,
		OrderFieldTitle,
		OrderFieldMaDatHangMm,
		OrderFieldMaHopDongKhachHang,
		OrderFieldMaHopDong,
		OrderFieldSaleName,
		OrderFieldSaleAdminName,
		OrderFieldStatus,
		OrderFieldCreatedBy,
		OrderFieldUpdatedBy,
		OrderFieldCreatedAt,
		OrderFieldUpdatedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.Title,
		&rcv.MaDatHangMm,
		&rcv.MaHopDongKhachHang,
		&rcv.MaHopDong,
		&rcv.SaleName,
		&rcv.SaleAdminName,
		&rcv.Status,
		&rcv.CreatedBy,
		&rcv.UpdatedBy,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
	}

	return
}

func (*Order) TableName() string {
	return "orders"
}
