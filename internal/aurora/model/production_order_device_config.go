package model

import (
	"database/sql"
	"time"
)

const (
	ProductionOrderDeviceConfigFieldID                = "id"
	ProductionOrderDeviceConfigFieldProductionOrderID = "production_order_id"
	ProductionOrderDeviceConfigFieldDeviceID          = "device_id"
	ProductionOrderDeviceConfigFieldDeviceConfig      = "device_config"
	ProductionOrderDeviceConfigFieldCreatedAt         = "created_at"
	ProductionOrderDeviceConfigFieldUpdatedAt         = "updated_at"
	ProductionOrderDeviceConfigFieldDeletedAt         = "deleted_at"
)

type ProductionOrderDeviceConfig struct {
	ID                string                 `db:"id"`
	ProductionOrderID string                 `db:"production_order_id"`
	DeviceID          string                 `db:"device_id"`
	DeviceConfig      map[string]interface{} `db:"device_config"`
	CreatedAt         time.Time              `db:"created_at"`
	UpdatedAt         time.Time              `db:"updated_at"`
	DeletedAt         sql.NullTime           `db:"deleted_at"`
}

func (rcv *ProductionOrderDeviceConfig) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		ProductionOrderDeviceConfigFieldID,
		ProductionOrderDeviceConfigFieldProductionOrderID,
		ProductionOrderDeviceConfigFieldDeviceID,
		ProductionOrderDeviceConfigFieldDeviceConfig,
		ProductionOrderDeviceConfigFieldCreatedAt,
		ProductionOrderDeviceConfigFieldUpdatedAt,
		ProductionOrderDeviceConfigFieldDeletedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.ProductionOrderID,
		&rcv.DeviceID,
		&rcv.DeviceConfig,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
		&rcv.DeletedAt,
	}

	return
}

func (*ProductionOrderDeviceConfig) TableName() string {
	return "production_order_device_config"
}
