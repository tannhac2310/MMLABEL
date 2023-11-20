package model

import (
	"database/sql"
	"time"
)

const (
	InkExportFieldID                = "id"
	InkExportFieldCode              = "code"
	InkExportFieldProductionOrderID = "production_order_id"
	InkExportFieldExportDate        = "export_date"
	InkExportFieldExportUser        = "export_user"
	InkExportFieldExportWarehouse   = "export_warehouse"
	InkExportFieldDescription       = "description"
	InkExportFieldStatus            = "status"
	InkExportFieldData              = "data"
	InkExportFieldCreatedAt         = "created_at"
	InkExportFieldUpdatedAt         = "updated_at"
	InkExportFieldDeletedAt         = "deleted_at"
)

type InkExport struct {
	ID                string                 `db:"id"`
	Code              string                 `db:"code"`
	ProductionOrderID string                 `db:"production_order_id"`
	ExportDate        time.Time              `db:"export_date"`
	ExportUser        string                 `db:"export_user"`
	ExportWarehouse   string                 `db:"export_warehouse"`
	Description       sql.NullString         `db:"description"`
	Status            int16                  `db:"status"`
	Data              map[string]interface{} `db:"data"`
	CreatedAt         time.Time              `db:"created_at"`
	UpdatedAt         time.Time              `db:"updated_at"`
	DeletedAt         sql.NullTime           `db:"deleted_at"`
}

func (rcv *InkExport) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		InkExportFieldID,
		InkExportFieldCode,
		InkExportFieldProductionOrderID,
		InkExportFieldExportDate,
		InkExportFieldExportUser,
		InkExportFieldExportWarehouse,
		InkExportFieldDescription,
		InkExportFieldStatus,
		InkExportFieldData,
		InkExportFieldCreatedAt,
		InkExportFieldUpdatedAt,
		InkExportFieldDeletedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.Code,
		&rcv.ProductionOrderID,
		&rcv.ExportDate,
		&rcv.ExportUser,
		&rcv.ExportWarehouse,
		&rcv.Description,
		&rcv.Status,
		&rcv.Data,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
		&rcv.DeletedAt,
	}

	return
}

func (*InkExport) TableName() string {
	return "ink_export"
}
