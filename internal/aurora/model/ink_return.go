package model

import (
	"database/sql"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"time"
)

const (
	InkReturnFieldID              = "id"
	InkReturnFieldName            = "name"
	InkReturnFieldCode            = "code"
	InkReturnFieldInkExportID     = "ink_export_id"
	InkReturnFieldReturnDate      = "return_date"
	InkReturnFieldReturnWarehouse = "return_warehouse"
	InkReturnFieldDescription     = "description"
	InkReturnFieldStatus          = "status"
	InkReturnFieldData            = "data"
	InkReturnFieldCreatedBy       = "created_by"
	InkReturnFieldUpdatedBy       = "updated_by"
	InkReturnFieldCreatedAt       = "created_at"
	InkReturnFieldUpdatedAt       = "updated_at"
	InkReturnFieldDeletedAt       = "deleted_at"
)

type InkReturn struct {
	ID              string                     `db:"id"`
	Name            string                     `db:"name"`
	Code            string                     `db:"code"`
	InkExportID     string                     `db:"ink_export_id"`
	ReturnDate      sql.NullTime               `db:"return_date"`
	ReturnWarehouse sql.NullString             `db:"return_warehouse"`
	Description     sql.NullString             `db:"description"`
	Status          enum.InventoryCommonStatus `db:"status"`
	Data            map[string]interface{}     `db:"data"`
	CreatedBy       string                     `db:"created_by"`
	UpdatedBy       string                     `db:"updated_by"`
	CreatedAt       time.Time                  `db:"created_at"`
	UpdatedAt       time.Time                  `db:"updated_at"`
	DeletedAt       sql.NullTime               `db:"deleted_at"`
}

func (rcv *InkReturn) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		InkReturnFieldID,
		InkReturnFieldName,
		InkReturnFieldCode,
		InkReturnFieldInkExportID,
		InkReturnFieldReturnDate,
		InkReturnFieldReturnWarehouse,
		InkReturnFieldDescription,
		InkReturnFieldStatus,
		InkReturnFieldData,
		InkReturnFieldCreatedBy,
		InkReturnFieldUpdatedBy,
		InkReturnFieldCreatedAt,
		InkReturnFieldUpdatedAt,
		InkReturnFieldDeletedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.Name,
		&rcv.Code,
		&rcv.InkExportID,
		&rcv.ReturnDate,
		&rcv.ReturnWarehouse,
		&rcv.Description,
		&rcv.Status,
		&rcv.Data,
		&rcv.CreatedBy,
		&rcv.UpdatedBy,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
		&rcv.DeletedAt,
	}

	return
}

func (*InkReturn) TableName() string {
	return "ink_return"
}
