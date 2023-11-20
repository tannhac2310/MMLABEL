package model

import (
	"database/sql"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"time"
)

const (
	InkInventoryFieldID                 = "id"
	InkInventoryFieldCode               = "code"
	InkInventoryFieldInventoryDate      = "inventory_date"
	InkInventoryFieldInventoryUser      = "inventory_user"
	InkInventoryFieldInventoryWarehouse = "inventory_warehouse"
	InkInventoryFieldDescription        = "description"
	InkInventoryFieldStatus             = "status"
	InkInventoryFieldData               = "data"
	InkInventoryFieldCreatedAt          = "created_at"
	InkInventoryFieldUpdatedAt          = "updated_at"
	InkInventoryFieldDeletedAt          = "deleted_at"
)

type InkInventory struct {
	ID                 string                           `db:"id"`
	Code               string                           `db:"code"`
	InventoryDate      time.Time                        `db:"inventory_date"`
	InventoryUser      string                           `db:"inventory_user"`
	InventoryWarehouse string                           `db:"inventory_warehouse"`
	Description        sql.NullString                   `db:"description"`
	Status             enum.InventoryCommonStatusStatus `db:"status"`
	Data               map[string]interface{}           `db:"data"`
	CreatedAt          time.Time                        `db:"created_at"`
	UpdatedAt          time.Time                        `db:"updated_at"`
	DeletedAt          sql.NullTime                     `db:"deleted_at"`
}

func (rcv *InkInventory) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		InkInventoryFieldID,
		InkInventoryFieldCode,
		InkInventoryFieldInventoryDate,
		InkInventoryFieldInventoryUser,
		InkInventoryFieldInventoryWarehouse,
		InkInventoryFieldDescription,
		InkInventoryFieldStatus,
		InkInventoryFieldData,
		InkInventoryFieldCreatedAt,
		InkInventoryFieldUpdatedAt,
		InkInventoryFieldDeletedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.Code,
		&rcv.InventoryDate,
		&rcv.InventoryUser,
		&rcv.InventoryWarehouse,
		&rcv.Description,
		&rcv.Status,
		&rcv.Data,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
		&rcv.DeletedAt,
	}

	return
}

func (*InkInventory) TableName() string {
	return "ink_inventory"
}
