package model

import (
	"database/sql"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

/*
ALTER TABLE orders

	ADD COLUMN payment_method varchar(255) default '',
	ADD COLUMN payment_method_other varchar(255) default '',
	ADD COLUMN customer_id int default '',
	ADD COLUMN customer_address_options text default '',
	ADD COLUMN delivery_address text default '';
*/
const (
	OrderFieldID                     = "id"
	OrderFieldTitle                  = "title"
	OrderFieldMaDatHangMm            = "ma_dat_hang_mm"
	OrderFieldMaHopDongKhachHang     = "ma_hop_dong_khach_hang"
	OrderFieldMaHopDong              = "ma_hop_dong"
	OrderFieldSaleName               = "sale_name"
	OrderFieldSaleAdminName          = "sale_admin_name"
	OrderFieldPaymentMethod          = "payment_method"
	OrderFieldPaymentMethodOther     = "payment_method_other"
	OrderFieldCustomerID             = "customer_id"
	OrderFieldCustomerAddressOptions = "customer_address_options"
	OrderFieldDeliveryAddress        = "delivery_address"
	OrderFieldStatus                 = "status"
	OrderFieldCreatedBy              = "created_by"
	OrderFieldUpdatedBy              = "updated_by"
	OrderFieldCreatedAt              = "created_at"
	OrderFieldUpdatedAt              = "updated_at"
)

type Order struct {
	ID                     string           `db:"id"`
	Title                  string           `db:"title"`
	MaDatHangMm            string           `db:"ma_dat_hang_mm"`
	MaHopDongKhachHang     string           `db:"ma_hop_dong_khach_hang"`
	MaHopDong              string           `db:"ma_hop_dong"`
	SaleName               sql.NullString   `db:"sale_name"`
	SaleAdminName          sql.NullString   `db:"sale_admin_name"`
	Status                 enum.OrderStatus `db:"status"`
	PaymentMethod          string           `db:"payment_method"`
	PaymentMethodOther     string           `db:"payment_method_other"`
	CustomerID             string           `db:"customer_id"`
	CustomerAddressOptions string           `db:"customer_address_options"`
	DeliveryAddress        string           `db:"delivery_address"`
	CreatedBy              string           `db:"created_by"`
	UpdatedBy              string           `db:"updated_by"`
	CreatedAt              time.Time        `db:"created_at"`
	UpdatedAt              time.Time        `db:"updated_at"`
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
		OrderFieldPaymentMethod,
		OrderFieldPaymentMethodOther,
		OrderFieldCustomerID,
		OrderFieldCustomerAddressOptions,
		OrderFieldDeliveryAddress,
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
		&rcv.PaymentMethod,
		&rcv.PaymentMethodOther,
		&rcv.CustomerID,
		&rcv.CustomerAddressOptions,
		&rcv.DeliveryAddress,
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
