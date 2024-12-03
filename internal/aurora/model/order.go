package model

import (
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

const (
	OrderFieldID                  = "id"
	OrderFieldTitle               = "title"
	OrderFieldCode                = "code"
	OrderFieldSaleName            = "sale_name"
	OrderFieldSaleAdminName       = "sale_admin_name"
	OrderFieldProductCode         = "product_code"
	OrderFieldProductName         = "product_name"
	OrderFieldCustomerID          = "customer_id"
	OrderFieldCustomerProductCode = "customer_product_code"
	OrderFieldCustomerProductName = "customer_product_name"
	OrderFieldStatus              = "status"
	OrderFieldCreatedBy           = "created_by"
	OrderFieldUpdatedBy           = "updated_by"
	OrderFieldCreatedAt           = "created_at"
	OrderFieldUpdatedAt           = "updated_at"
)

type Order struct {
	ID                  string           `db:"id"`
	Title               string           `db:"title"`
	Code                string           `db:"code"`
	SaleName            string           `db:"sale_name"`
	SaleAdminName       string           `db:"sale_admin_name"`
	ProductCode         string           `db:"product_code"`
	ProductName         string           `db:"product_name"`
	CustomerID          string           `db:"customer_id"`
	CustomerProductCode string           `db:"customer_product_code"`
	CustomerProductName string           `db:"customer_product_name"`
	Status              enum.OrderStatus `db:"status"`
	CreatedBy           string           `db:"created_by"`
	UpdatedBy           string           `db:"updated_by"`
	CreatedAt           time.Time        `db:"created_at"`
	UpdatedAt           time.Time        `db:"updated_at"`
}

func (rcv *Order) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		OrderFieldID,
		OrderFieldTitle,
		OrderFieldCode,
		OrderFieldSaleName,
		OrderFieldSaleAdminName,
		OrderFieldProductCode,
		OrderFieldProductName,
		OrderFieldCustomerID,
		OrderFieldCustomerProductCode,
		OrderFieldCustomerProductName,
		OrderFieldStatus,
		OrderFieldCreatedBy,
		OrderFieldUpdatedBy,
		OrderFieldCreatedAt,
		OrderFieldUpdatedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.Title,
		&rcv.Code,
		&rcv.SaleName,
		&rcv.SaleAdminName,
		&rcv.ProductCode,
		&rcv.ProductName,
		&rcv.CustomerID,
		&rcv.CustomerProductCode,
		&rcv.CustomerProductName,
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
