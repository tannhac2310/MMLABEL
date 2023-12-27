package model

import (
	"database/sql"
	"time"
)

const (
	DeviceWorkingHistoryFieldID                           = "id"
	DeviceWorkingHistoryFieldProductionOrderStageDeviceID = "production_order_stage_device_id"
	DeviceWorkingHistoryFieldDeviceID                     = "device_id"
	DeviceWorkingHistoryFieldDate                         = "date"
	DeviceWorkingHistoryFieldQuantity                     = "quantity"
	DeviceWorkingHistoryFieldWorkingTime                  = "working_time"
	DeviceWorkingHistoryFieldPoQuantity                   = "po_quantity"
	DeviceWorkingHistoryFieldPoWorkingTime                = "po_working_time"
	DeviceWorkingHistoryFieldUpdatedAt                    = "updated_at"
	DeviceWorkingHistoryFieldCreatedAt                    = "created_at"
)

type DeviceWorkingHistory struct {
	ID                           string       `db:"id"`
	ProductionOrderStageDeviceID string       `db:"production_order_stage_device_id"`
	DeviceID                     string       `db:"device_id"`
	Date                         string       `db:"date"`
	Quantity                     int64        `db:"quantity"`
	WorkingTime                  int64        `db:"working_time"`
	PoQuantity                   int64        `db:"po_quantity"`
	PoWorkingTime                int64        `db:"po_working_time"`
	UpdatedAt                    sql.NullTime `db:"updated_at"`
	CreatedAt                    time.Time    `db:"created_at"`
}

func (rcv *DeviceWorkingHistory) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		DeviceWorkingHistoryFieldID,
		DeviceWorkingHistoryFieldProductionOrderStageDeviceID,
		DeviceWorkingHistoryFieldDeviceID,
		DeviceWorkingHistoryFieldDate,
		DeviceWorkingHistoryFieldQuantity,
		DeviceWorkingHistoryFieldWorkingTime,
		DeviceWorkingHistoryFieldPoQuantity,
		DeviceWorkingHistoryFieldPoWorkingTime,
		DeviceWorkingHistoryFieldUpdatedAt,
		DeviceWorkingHistoryFieldCreatedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.ProductionOrderStageDeviceID,
		&rcv.DeviceID,
		&rcv.Date,
		&rcv.Quantity,
		&rcv.WorkingTime,
		&rcv.PoQuantity,
		&rcv.PoWorkingTime,
		&rcv.UpdatedAt,
		&rcv.CreatedAt,
	}

	return
}

func (*DeviceWorkingHistory) TableName() string {
	return "device_working_history"
}
