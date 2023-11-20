package model

import (
	"database/sql"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"time"
)

const (
	InkReturnFieldID              = "id"
	InkReturnFieldCode            = "code"
	InkReturnFieldReturnDate      = "return_date"
	InkReturnFieldReturnUser      = "return_user"
	InkReturnFieldReturnWarehouse = "return_warehouse"
	InkReturnFieldDescription     = "description"
	InkReturnFieldStatus          = "status"
	InkReturnFieldData            = "data"
	InkReturnFieldCreatedAt       = "created_at"
	InkReturnFieldUpdatedAt       = "updated_at"
	InkReturnFieldDeletedAt       = "deleted_at"
)

type InkReturn struct {
	ID              string                           `db:"id"`
	Code            string                           `db:"code"`
	ReturnDate      time.Time                        `db:"return_date"`
	ReturnUser      string                           `db:"return_user"`
	ReturnWarehouse string                           `db:"return_warehouse"`
	Description     sql.NullString                   `db:"description"`
	Status          enum.InventoryCommonStatusStatus `db:"status"`
	Data            map[string]interface{}           `db:"data"`
	CreatedAt       time.Time                        `db:"created_at"`
	UpdatedAt       time.Time                        `db:"updated_at"`
	DeletedAt       sql.NullTime                     `db:"deleted_at"`
}

func (rcv *InkReturn) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		InkReturnFieldID,
		InkReturnFieldCode,
		InkReturnFieldReturnDate,
		InkReturnFieldReturnUser,
		InkReturnFieldReturnWarehouse,
		InkReturnFieldDescription,
		InkReturnFieldStatus,
		InkReturnFieldData,
		InkReturnFieldCreatedAt,
		InkReturnFieldUpdatedAt,
		InkReturnFieldDeletedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.Code,
		&rcv.ReturnDate,
		&rcv.ReturnUser,
		&rcv.ReturnWarehouse,
		&rcv.Description,
		&rcv.Status,
		&rcv.Data,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
		&rcv.DeletedAt,
	}

	return
}

func (*InkReturn) TableName() string {
	return "ink_return"
}
