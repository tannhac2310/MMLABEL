package model

import (
	"database/sql"
	"time"
)

const (
	InkImportDetailFieldID          = "id"
	InkImportDetailFieldCode        = "code"
	InkImportDetailFieldInkCode     = "ink_code"
	InkImportDetailFieldQuantity    = "quantity"
	InkImportDetailFieldColorDetail = "color_detail"
	InkImportDetailFieldDescription = "description"
	InkImportDetailFieldData        = "data"
	InkImportDetailFieldCreatedAt   = "created_at"
	InkImportDetailFieldUpdatedAt   = "updated_at"
	InkImportDetailFieldDeletedAt   = "deleted_at"
)

type InkImportDetail struct {
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

func (rcv *InkImportDetail) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		InkImportDetailFieldID,
		InkImportDetailFieldCode,
		InkImportDetailFieldInkCode,
		InkImportDetailFieldQuantity,
		InkImportDetailFieldColorDetail,
		InkImportDetailFieldDescription,
		InkImportDetailFieldData,
		InkImportDetailFieldCreatedAt,
		InkImportDetailFieldUpdatedAt,
		InkImportDetailFieldDeletedAt,
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

func (*InkImportDetail) TableName() string {
	return "ink_import_detail"
}
