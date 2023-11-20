package model

import (
	"database/sql"
	"time"
)

const (
	InkInventoryDetailFieldID          = "id"
	InkInventoryDetailFieldCode        = "code"
	InkInventoryDetailFieldInkCode     = "ink_code"
	InkInventoryDetailFieldQuantity    = "quantity"
	InkInventoryDetailFieldColorDetail = "color_detail"
	InkInventoryDetailFieldDescription = "description"
	InkInventoryDetailFieldData        = "data"
	InkInventoryDetailFieldCreatedAt   = "created_at"
	InkInventoryDetailFieldUpdatedAt   = "updated_at"
	InkInventoryDetailFieldDeletedAt   = "deleted_at"
)

type InkInventoryDetail struct {
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

func (rcv *InkInventoryDetail) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		InkInventoryDetailFieldID,
		InkInventoryDetailFieldCode,
		InkInventoryDetailFieldInkCode,
		InkInventoryDetailFieldQuantity,
		InkInventoryDetailFieldColorDetail,
		InkInventoryDetailFieldDescription,
		InkInventoryDetailFieldData,
		InkInventoryDetailFieldCreatedAt,
		InkInventoryDetailFieldUpdatedAt,
		InkInventoryDetailFieldDeletedAt,
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

func (*InkInventoryDetail) TableName() string {
	return "ink_inventory_detail"
}
