package model

import (
	"database/sql"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"time"
)

const (
	InkImportFieldID          = "id"
	InkImportFieldName        = "name"
	InkImportFieldCode        = "code"
	InkImportFieldImportDate  = "import_date"
	InkImportFieldDescription = "description"
	InkImportFieldStatus      = "status"
	InkImportFieldData        = "data"
	InkImportFieldCreatedBy   = "created_by"
	InkImportFieldUpdatedBy   = "updated_by"
	InkImportFieldCreatedAt   = "created_at"
	InkImportFieldUpdatedAt   = "updated_at"
	InkImportFieldDeletedAt   = "deleted_at"
)

type InkImport struct {
	ID          string                     `db:"id"`
	Name        string                     `db:"name"`
	Code        string                     `db:"code"`
	ImportDate  sql.NullTime               `db:"import_date"`
	Description sql.NullString             `db:"description"`
	Status      enum.InventoryCommonStatus `db:"status"`
	Data        map[string]interface{}     `db:"data"`
	CreatedBy   string                     `db:"created_by"`
	UpdatedBy   string                     `db:"updated_by"`
	CreatedAt   time.Time                  `db:"created_at"`
	UpdatedAt   time.Time                  `db:"updated_at"`
	DeletedAt   sql.NullTime               `db:"deleted_at"`
}

func (rcv *InkImport) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		InkImportFieldID,
		InkImportFieldName,
		InkImportFieldCode,
		InkImportFieldImportDate,
		InkImportFieldDescription,
		InkImportFieldStatus,
		InkImportFieldData,
		InkImportFieldCreatedBy,
		InkImportFieldUpdatedBy,
		InkImportFieldCreatedAt,
		InkImportFieldUpdatedAt,
		InkImportFieldDeletedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.Name,
		&rcv.Code,
		&rcv.ImportDate,
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

func (*InkImport) TableName() string {
	return "ink_import"
}
