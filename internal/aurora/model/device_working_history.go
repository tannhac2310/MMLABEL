package model

import (
	"database/sql"
	"time"
)

const (
	DeviceWorkingHistoryFieldID        						= "id"
	DeviceWorkingHistoryFieldProductionOrderStageDeviceID 	= "production_order_stage_device_id"
	DeviceWorkingHistoryFieldDeviceID                     	= "device_id"
	DeviceWorkingHistoryFieldQty                     		= "qty"
	DeviceWorkingHistoryFieldWorkingDate	    			= "working_date"
	DeviceWorkingHistoryFieldCreatedAt 						= "created_at"
	DeviceWorkingHistoryFieldUpdatedAt 						= "updated_at"
)

type DeviceWorkingHistory struct {
	ID        						string                 	`db:"id"`
	ProductionOrderStageDeviceID    string                 	`db:"production_order_stage_device_id"`
	DeviceID      					string                	`db:"device_id"`
	Qty      						int           			`db:"qty"`
	WorkingDate 					sql.NullTime            `db:"working_date"`
	CreatedAt 						time.Time              	`db:"created_at"`
	UpdatedAt 						sql.NullTime            `db:"updated_at"`
}

func (rcv *DeviceWorkingHistory) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		DeviceFieldID,
		DeviceWorkingHistoryFieldProductionOrderStageDeviceID,
		DeviceWorkingHistoryFieldDeviceID,
		DeviceWorkingHistoryFieldQty,
		DeviceWorkingHistoryFieldWorkingDate,
		DeviceWorkingHistoryFieldCreatedAt,
		DeviceWorkingHistoryFieldUpdatedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.ProductionOrderStageDeviceID,
		&rcv.DeviceID,
		&rcv.Qty,
		&rcv.WorkingDate,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
	}

	return
}

func (*DeviceWorkingHistory) TableName() string {
	return "device_working_history"
}
