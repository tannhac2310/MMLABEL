package model

import (
	"database/sql"
	"time"
)

const (
	InkInventoryDetailFieldID             = "id"
	InkInventoryDetailFieldInkInventoryID = "ink_inventory_id"
	InkInventoryDetailFieldQuantity       = "quantity"
	InkInventoryDetailFieldColorDetail    = "color_detail"
	InkInventoryDetailFieldDescription    = "description"
	InkInventoryDetailFieldData           = "data"
	InkInventoryDetailFieldCreatedAt      = "created_at"
	InkInventoryDetailFieldUpdatedAt      = "updated_at"
	InkInventoryDetailFieldDeletedAt      = "deleted_at"
)

type InkInventoryDetail struct {
	ID             string                 `db:"id"`
	InkInventoryID string                 `db:"ink_inventory_id"`
	Quantity       float64                `db:"quantity"`
	ColorDetail    map[string]interface{} `db:"color_detail"`
	Description    sql.NullString         `db:"description"`
	Data           map[string]interface{} `db:"data"`
	CreatedAt      time.Time              `db:"created_at"`
	UpdatedAt      time.Time              `db:"updated_at"`
	DeletedAt      sql.NullTime           `db:"deleted_at"`
}

func (rcv *InkInventoryDetail) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		InkInventoryDetailFieldID,
		InkInventoryDetailFieldInkInventoryID,
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
		&rcv.InkInventoryID,
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
