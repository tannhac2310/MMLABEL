package model

import (
	"database/sql"
	"time"
)

const (
	InkMixingDetailFieldID          = "id"
	InkMixingDetailFieldInkMixingID = "ink_mixing_id"
	InkMixingDetailFieldInkID       = "ink_id"
	InkMixingDetailFieldQuantity    = "quantity"
	InkMixingDetailFieldDescription = "description"
	InkMixingDetailFieldData        = "data"
	InkMixingDetailFieldCreatedAt   = "created_at"
	InkMixingDetailFieldUpdatedAt   = "updated_at"
	InkMixingDetailFieldDeletedAt   = "deleted_at"
)

type InkMixingDetail struct {
	ID          string                 `db:"id"`
	InkMixingID string                 `db:"ink_mixing_id"`
	InkID       string                 `db:"ink_id"`
	Quantity    float64                `db:"quantity"`
	Description string                 `db:"description"`
	Data        map[string]interface{} `db:"data"`
	CreatedAt   time.Time              `db:"created_at"`
	UpdatedAt   time.Time              `db:"updated_at"`
	DeletedAt   sql.NullTime           `db:"deleted_at"`
}

func (rcv *InkMixingDetail) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		InkMixingDetailFieldID,
		InkMixingDetailFieldInkMixingID,
		InkMixingDetailFieldInkID,
		InkMixingDetailFieldQuantity,
		InkMixingDetailFieldDescription,
		InkMixingDetailFieldData,
		InkMixingDetailFieldCreatedAt,
		InkMixingDetailFieldUpdatedAt,
		InkMixingDetailFieldDeletedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.InkMixingID,
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

func (*InkMixingDetail) TableName() string {
	return "ink_mixing_detail"
}
