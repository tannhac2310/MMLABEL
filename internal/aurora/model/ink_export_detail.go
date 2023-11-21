package model

import (
	"database/sql"
	"time"
)

const (
	InkExportDetailFieldID          = "id"
	InkExportDetailFieldInkExportID = "ink_export_id"
	InkExportDetailFieldQuantity    = "quantity"
	InkExportDetailFieldColorDetail = "color_detail"
	InkExportDetailFieldDescription = "description"
	InkExportDetailFieldData        = "data"
	InkExportDetailFieldCreatedAt   = "created_at"
	InkExportDetailFieldUpdatedAt   = "updated_at"
	InkExportDetailFieldDeletedAt   = "deleted_at"
)

type InkExportDetail struct {
	ID          string                 `db:"id"`
	InkExportID string                 `db:"ink_export_id"`
	Quantity    float64                `db:"quantity"`
	ColorDetail map[string]interface{} `db:"color_detail"`
	Description sql.NullString         `db:"description"`
	Data        map[string]interface{} `db:"data"`
	CreatedAt   time.Time              `db:"created_at"`
	UpdatedAt   time.Time              `db:"updated_at"`
	DeletedAt   sql.NullTime           `db:"deleted_at"`
}

func (rcv *InkExportDetail) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		InkExportDetailFieldID,
		InkExportDetailFieldInkExportID,
		InkExportDetailFieldQuantity,
		InkExportDetailFieldColorDetail,
		InkExportDetailFieldDescription,
		InkExportDetailFieldData,
		InkExportDetailFieldCreatedAt,
		InkExportDetailFieldUpdatedAt,
		InkExportDetailFieldDeletedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.InkExportID,
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

func (*InkExportDetail) TableName() string {
	return "ink_export_detail"
}
