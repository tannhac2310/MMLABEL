package model

import (
	"database/sql"
	"time"
)

const (
	ProductionOrderDeviceConfigFieldID                = "id"
	ProductionOrderDeviceConfigFieldProductionOrderID = "production_order_id"
	ProductionOrderDeviceConfigFieldDeviceID          = "device_id"
	ProductionOrderDeviceConfigFieldColor             = "color"
	ProductionOrderDeviceConfigFieldDescription       = "description"
	ProductionOrderDeviceConfigFieldSearch            = "search"
	ProductionOrderDeviceConfigFieldDeviceConfig      = "device_config"
	ProductionOrderDeviceConfigFieldCreatedBy         = "created_by"
	ProductionOrderDeviceConfigFieldCreatedAt         = "created_at"
	ProductionOrderDeviceConfigFieldUpdatedBy         = "updated_by"
	ProductionOrderDeviceConfigFieldUpdatedAt         = "updated_at"
	ProductionOrderDeviceConfigFieldDeletedAt         = "deleted_at"
)

type ProductionOrderDeviceConfig struct {
	ID                string                 `db:"id"`
	ProductionOrderID string                 `db:"production_order_id"`
	DeviceID          sql.NullString         `db:"device_id"`
	Color             sql.NullString         `db:"color"`
	Description       sql.NullString         `db:"description"`
	Search            sql.NullString         `db:"search"`
	DeviceConfig      map[string]interface{} `db:"device_config"`
	CreatedBy         string                 `db:"created_by"`
	CreatedAt         time.Time              `db:"created_at"`
	UpdatedBy         string                 `db:"updated_by"`
	UpdatedAt         time.Time              `db:"updated_at"`
	DeletedAt         sql.NullTime           `db:"deleted_at"`
}

func (rcv *ProductionOrderDeviceConfig) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		ProductionOrderDeviceConfigFieldID,
		ProductionOrderDeviceConfigFieldProductionOrderID,
		ProductionOrderDeviceConfigFieldDeviceID,
		ProductionOrderDeviceConfigFieldColor,
		ProductionOrderDeviceConfigFieldDescription,
		ProductionOrderDeviceConfigFieldSearch,
		ProductionOrderDeviceConfigFieldDeviceConfig,
		ProductionOrderDeviceConfigFieldCreatedBy,
		ProductionOrderDeviceConfigFieldCreatedAt,
		ProductionOrderDeviceConfigFieldUpdatedBy,
		ProductionOrderDeviceConfigFieldUpdatedAt,
		ProductionOrderDeviceConfigFieldDeletedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.ProductionOrderID,
		&rcv.DeviceID,
		&rcv.Color,
		&rcv.Description,
		&rcv.Search,
		&rcv.DeviceConfig,
		&rcv.CreatedBy,
		&rcv.CreatedAt,
		&rcv.UpdatedBy,
		&rcv.UpdatedAt,
		&rcv.DeletedAt,
	}

	return
}

func (*ProductionOrderDeviceConfig) TableName() string {
	return "production_order_device_config"
}
