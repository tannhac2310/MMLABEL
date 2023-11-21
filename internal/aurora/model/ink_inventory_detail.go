package model

import (
	"database/sql"
	"time"
)

const (
	InkInventoryDetailFieldID             = "id"
	InkInventoryDetailFieldInkInventoryID = "ink_inventory_id"
	InkInventoryDetailFieldInkID          = "ink_id"
	InkInventoryDetailFieldQuantity       = "quantity"
	InkInventoryDetailFieldDescription    = "description"
	InkInventoryDetailFieldData           = "data"
	InkInventoryDetailFieldCreatedAt      = "created_at"
	InkInventoryDetailFieldUpdatedAt      = "updated_at"
	InkInventoryDetailFieldDeletedAt      = "deleted_at"
)

type InkInventoryDetail struct {
	ID             string                 `db:"id"`
	InkInventoryID string                 `db:"ink_inventory_id"`
	InkID          string                 `db:"ink_id"`
	Quantity       float64                `db:"quantity"`
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
		InkInventoryDetailFieldInkID,
		InkInventoryDetailFieldQuantity,
		InkInventoryDetailFieldDescription,
		InkInventoryDetailFieldData,
		InkInventoryDetailFieldCreatedAt,
		InkInventoryDetailFieldUpdatedAt,
		InkInventoryDetailFieldDeletedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.InkInventoryID,
		&rcv.InkID,
		&rcv.Quantity,
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
