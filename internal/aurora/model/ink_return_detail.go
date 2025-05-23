package model

import (
	"database/sql"
	"time"
)

const (
	InkReturnDetailFieldID                = "id"
	InkReturnDetailFieldInkReturnID       = "ink_return_id"
	InkReturnDetailFieldInkID             = "ink_id"
	InkReturnDetailFieldInkExportID       = "ink_export_id"
	InkReturnDetailFieldInkExportDetailID = "ink_export_detail_id"
	InkReturnDetailFieldQuantity          = "quantity"
	InkReturnDetailFieldColorDetail       = "color_detail"
	InkReturnDetailFieldDescription       = "description"
	InkReturnDetailFieldData              = "data"
	InkReturnDetailFieldCreatedAt         = "created_at"
	InkReturnDetailFieldUpdatedAt         = "updated_at"
	InkReturnDetailFieldDeletedAt         = "deleted_at"
)

type InkReturnDetail struct {
	ID                string                 `db:"id"`
	InkReturnID       string                 `db:"ink_return_id"`
	InkID             string                 `db:"ink_id"`
	InkExportID       string                 `db:"ink_export_id"` // todo remove this field
	InkExportDetailID string                 `db:"ink_export_detail_id"`
	Quantity          float64                `db:"quantity"`
	ColorDetail       map[string]interface{} `db:"color_detail"`
	Description       sql.NullString         `db:"description"`
	Data              map[string]interface{} `db:"data"`
	CreatedAt         time.Time              `db:"created_at"`
	UpdatedAt         time.Time              `db:"updated_at"`
	DeletedAt         sql.NullTime           `db:"deleted_at"`
}

func (rcv *InkReturnDetail) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		InkReturnDetailFieldID,
		InkReturnDetailFieldInkReturnID,
		InkReturnDetailFieldInkID,
		InkReturnDetailFieldInkExportID,
		InkReturnDetailFieldInkExportDetailID,
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
		&rcv.InkReturnID,
		&rcv.InkID,
		&rcv.InkExportID,
		&rcv.InkExportDetailID,
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
