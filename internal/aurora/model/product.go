package model

import (
	"database/sql"
	"time"
)

const (
	ProductFieldID          = "id"
	ProductFieldName        = "name"
	ProductFieldCode        = "code"
	ProductFieldCustomerID  = "customer_id"
	ProductFieldSaleID      = "sale_id"
	ProductFieldDescription = "description"
	ProductFieldData        = "data"
	ProductFieldCreatedAt   = "created_at"
	ProductFieldUpdatedAt   = "updated_at"
	ProductFieldCreatedBy   = "created_by"
	ProductFieldUpdatedBy   = "updated_by"
	ProductFieldDeletedAt   = "deleted_at"
)

type Product struct {
	ID          string       `db:"id"`
	Name        string       `db:"name"`
	Code        string       `db:"code"`
	CustomerID  string       `db:"customer_id"`
	SaleID      string       `db:"sale_id"`
	Description string       `db:"description"`
	Data        any          `db:"data"`
	CreatedAt   time.Time    `db:"created_at"`
	UpdatedAt   time.Time    `db:"updated_at"`
	CreatedBy   string       `db:"created_by"`
	UpdatedBy   string       `db:"updated_by"`
	DeletedAt   sql.NullTime `db:"deleted_at"`
}

func (rcv *Product) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		ProductFieldID,
		ProductFieldName,
		ProductFieldCode,
		ProductFieldCustomerID,
		ProductFieldSaleID,
		ProductFieldDescription,
		ProductFieldData,
		ProductFieldCreatedAt,
		ProductFieldUpdatedAt,
		ProductFieldCreatedBy,
		ProductFieldUpdatedBy,
		ProductFieldDeletedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.Name,
		&rcv.Code,
		&rcv.CustomerID,
		&rcv.SaleID,
		&rcv.Description,
		&rcv.Data,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
		&rcv.CreatedBy,
		&rcv.UpdatedBy,
		&rcv.DeletedAt,
	}

	return
}

func (*Product) TableName() string {
	return "products"
}
