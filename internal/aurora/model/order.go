package model

import (
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

const (
	OrderFieldID        = "id"
	OrderFieldTitle     = "title"
	OrderFieldCode      = "code"
	OrderFieldStatus    = "status"
	OrderFieldCreatedBy = "created_by"
	OrderFieldUpdatedBy = "updated_by"
	OrderFieldCreatedAt = "created_at"
	OrderFieldUpdatedAt = "updated_at"
)

type Order struct {
	ID        string           `db:"id"`
	Title     string           `db:"title"`
	Code      string           `db:"code"`
	Status    enum.OrderStatus `db:"status"`
	CreatedBy string           `db:"created_by"`
	UpdatedBy string           `db:"updated_by"`
	CreatedAt time.Time        `db:"created_at"`
	UpdatedAt time.Time        `db:"updated_at"`
}

func (rcv *Order) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		OrderFieldID,
		OrderFieldTitle,
		OrderFieldCode,
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
