package model

import (
	"database/sql"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
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
	InkExportFieldCreatedBy         = "created_by"
	InkExportFieldUpdatedBy         = "updated_by"
	InkExportFieldCreatedAt         = "created_at"
	InkExportFieldUpdatedAt         = "updated_at"
	InkExportFieldDeletedAt         = "deleted_at"
)

type InkExport struct {
	ID                string                     `db:"id"`
	Code              string                     `db:"code"`
	ProductionOrderID string                     `db:"production_order_id"`
	ExportDate        time.Time                  `db:"export_date"`
	ExportUser        string                     `db:"export_user"`
	ExportWarehouse   string                     `db:"export_warehouse"`
	Description       sql.NullString             `db:"description"`
	Status            enum.InventoryCommonStatus `db:"status"`
	Data              map[string]interface{}     `db:"data"`
	CreatedBy         string                     `db:"created_by"`
	UpdatedBy         string                     `db:"updated_by"`
	CreatedAt         time.Time                  `db:"created_at"`
	UpdatedAt         time.Time                  `db:"updated_at"`
	DeletedAt         sql.NullTime               `db:"deleted_at"`
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
		InkExportFieldCreatedBy,
		InkExportFieldUpdatedBy,
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
		&rcv.CreatedBy,
		&rcv.UpdatedBy,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
		&rcv.DeletedAt,
	}

	return
}

func (*InkExport) TableName() string {
	return "ink_export"
}
