package model

import (
	"database/sql"
	"time"
)

const (
	InkReturnDetailFieldID          = "id"
	InkReturnDetailFieldCode        = "code"
	InkReturnDetailFieldInkCode     = "ink_code"
	InkReturnDetailFieldQuantity    = "quantity"
	InkReturnDetailFieldColorDetail = "color_detail"
	InkReturnDetailFieldDescription = "description"
	InkReturnDetailFieldData        = "data"
	InkReturnDetailFieldCreatedAt   = "created_at"
	InkReturnDetailFieldUpdatedAt   = "updated_at"
	InkReturnDetailFieldDeletedAt   = "deleted_at"
)

type InkReturnDetail struct {
	ID          string                 `db:"id"`
	Code        string                 `db:"code"`
	InkCode     string                 `db:"ink_code"`
	Quantity    float64                `db:"quantity"`
	ColorDetail map[string]interface{} `db:"color_detail"`
	Description sql.NullString         `db:"description"`
	Data        map[string]interface{} `db:"data"`
	CreatedAt   time.Time              `db:"created_at"`
	UpdatedAt   time.Time              `db:"updated_at"`
	DeletedAt   sql.NullTime           `db:"deleted_at"`
}

func (rcv *InkReturnDetail) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		InkReturnDetailFieldID,
		InkReturnDetailFieldCode,
		InkReturnDetailFieldInkCode,
		InkReturnDetailFieldQuantity,
		InkReturnDetailFieldColorDetail,
		InkReturnDetailFieldDescription,
		InkReturnDetailFieldData,
		InkReturnDetailFieldCreatedAt,
		InkReturnDetailFieldUpdatedAt,
		InkReturnDetailFieldDeletedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.Code,
		&rcv.InkCode,
		&rcv.Quantity,
		&rcv.ColorDetail,
		&rcv.Description,
		&rcv.Data,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
		&rcv.DeletedAt,
	}

	return
}

func (*InkReturnDetail) TableName() string {
	return "ink_return_detail"
}
