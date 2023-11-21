package model

import (
	"database/sql"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"time"
)

const (
	InkExportFieldID                = "id"
	InkExportFieldName              = "name"
	InkExportFieldCode              = "code"
	InkExportFieldProductionOrderID = "production_order_id"
	InkExportFieldExportDate        = "export_date"
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
	Name              string                     `db:"name"`
	Code              string                     `db:"code"`
	ProductionOrderID string                     `db:"production_order_id"`
	ExportDate        sql.NullTime               `db:"export_date"`
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
		InkExportFieldName,
		InkExportFieldCode,
		InkExportFieldProductionOrderID,
		InkExportFieldExportDate,
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
		&rcv.Name,
		&rcv.Code,
		&rcv.ProductionOrderID,
		&rcv.ExportDate,
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
