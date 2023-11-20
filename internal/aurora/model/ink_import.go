package model

import (
	"database/sql"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"time"
)

const (
	InkImportFieldID              = "id"
	InkImportFieldCode            = "code"
	InkImportFieldImportDate      = "import_date"
	InkImportFieldImportUser      = "import_user"
	InkImportFieldImportWarehouse = "import_warehouse"
	InkImportFieldExportWarehouse = "export_warehouse"
	InkImportFieldDescription     = "description"
	InkImportFieldStatus          = "status"
	InkImportFieldData            = "data"
	InkImportFieldCreatedAt       = "created_at"
	InkImportFieldUpdatedAt       = "updated_at"
	InkImportFieldDeletedAt       = "deleted_at"
)

type InkImport struct {
	ID              string                           `db:"id"`
	Code            string                           `db:"code"`
	ImportDate      time.Time                        `db:"import_date"`
	ImportUser      string                           `db:"import_user"`
	ImportWarehouse string                           `db:"import_warehouse"`
	ExportWarehouse string                           `db:"export_warehouse"`
	Description     sql.NullString                   `db:"description"`
	Status          enum.InventoryCommonStatusStatus `db:"status"`
	Data            map[string]interface{}           `db:"data"`
	CreatedAt       time.Time                        `db:"created_at"`
	UpdatedAt       time.Time                        `db:"updated_at"`
	DeletedAt       sql.NullTime                     `db:"deleted_at"`
}

func (rcv *InkImport) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		InkImportFieldID,
		InkImportFieldCode,
		InkImportFieldImportDate,
		InkImportFieldImportUser,
		InkImportFieldImportWarehouse,
		InkImportFieldExportWarehouse,
		InkImportFieldDescription,
		InkImportFieldStatus,
		InkImportFieldData,
		InkImportFieldCreatedAt,
		InkImportFieldUpdatedAt,
		InkImportFieldDeletedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.Code,
		&rcv.ImportDate,
		&rcv.ImportUser,
		&rcv.ImportWarehouse,
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

func (*InkImport) TableName() string {
	return "ink_import"
}
