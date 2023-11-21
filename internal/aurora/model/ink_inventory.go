package model

import (
	"database/sql"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"time"
)

const (
	InkInventoryFieldID                 = "id"
	InkInventoryFieldName               = "name"
	InkInventoryFieldCode               = "code"
	InkInventoryFieldInventoryDate      = "inventory_date"
	InkInventoryFieldInventoryUser      = "inventory_user"
	InkInventoryFieldInventoryWarehouse = "inventory_warehouse"
	InkInventoryFieldDescription        = "description"
	InkInventoryFieldStatus             = "status"
	InkInventoryFieldData               = "data"
	InkInventoryFieldCreatedBy          = "created_by"
	InkInventoryFieldUpdatedBy          = "updated_by"
	InkInventoryFieldCreatedAt          = "created_at"
	InkInventoryFieldUpdatedAt          = "updated_at"
	InkInventoryFieldDeletedAt          = "deleted_at"
)

type InkInventory struct {
	ID                 string                     `db:"id"`
	Name               string                     `db:"name"`
	Code               string                     `db:"code"`
	InventoryDate      sql.NullTime               `db:"inventory_date"`
	InventoryUser      string                     `db:"inventory_user"`
	InventoryWarehouse sql.NullString             `db:"inventory_warehouse"`
	Description        sql.NullString             `db:"description"`
	Status             enum.InventoryCommonStatus `db:"status"`
	Data               map[string]interface{}     `db:"data"`
	CreatedBy          string                     `db:"created_by"`
	UpdatedBy          string                     `db:"updated_by"`
	CreatedAt          time.Time                  `db:"created_at"`
	UpdatedAt          time.Time                  `db:"updated_at"`
	DeletedAt          sql.NullTime               `db:"deleted_at"`
}

func (rcv *InkInventory) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		InkInventoryFieldID,
		InkInventoryFieldName,
		InkInventoryFieldCode,
		InkInventoryFieldInventoryDate,
		InkInventoryFieldInventoryUser,
		InkInventoryFieldInventoryWarehouse,
		InkInventoryFieldDescription,
		InkInventoryFieldStatus,
		InkInventoryFieldData,
		InkInventoryFieldCreatedBy,
		InkInventoryFieldUpdatedBy,
		InkInventoryFieldCreatedAt,
		InkInventoryFieldUpdatedAt,
		InkInventoryFieldDeletedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.Name,
		&rcv.Code,
		&rcv.InventoryDate,
		&rcv.InventoryUser,
		&rcv.InventoryWarehouse,
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

func (*InkInventory) TableName() string {
	return "ink_inventory"
}
