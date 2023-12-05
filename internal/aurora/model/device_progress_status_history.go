package model

import (
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

const (
	DeviceProgressStatusHistoryFieldID                           = "id"
	DeviceProgressStatusHistoryFieldProductionOrderStageDeviceID = "production_order_stage_device_id"
	DeviceProgressStatusHistoryFieldDeviceID                     = "device_id"
	DeviceProgressStatusHistoryFieldProcessStatus                = "process_status"
	DeviceProgressStatusHistoryFieldCreatedAt                    = "created_at"
	DeviceProgressStatusHistoryFieldUpdatedBy                    = "updtead_by"
	DeviceProgressStatusHistoryFieldSolved                    	 = "solved"
)

type DeviceProgressStatusHistory struct {
	ID                           string                                `db:"id"`
	ProductionOrderStageDeviceID string                                `db:"production_order_stage_device_id"`
	DeviceID                     string                                `db:"device_id"`
	ProcessStatus                enum.ProductionOrderStageDeviceStatus `db:"process_status"`
	CreatedAt                    time.Time                             `db:"created_at"`
}

type DeviceProgressUpdateStatusHistory struct {
	ProcessStatus                enum.ProductionOrderStageDeviceStatus `db:"process_status"`
	UpdatedAt                    time.Time                             `db:"updated_at"`
	Solved						 bool								   `db:"solved"`
	UpdatedBy					 string								   `db:"updated_by"`
}

func (rcv *DeviceProgressStatusHistory) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		DeviceProgressStatusHistoryFieldID,
		DeviceProgressStatusHistoryFieldProductionOrderStageDeviceID,
		DeviceProgressStatusHistoryFieldDeviceID,
		DeviceProgressStatusHistoryFieldProcessStatus,
		DeviceProgressStatusHistoryFieldCreatedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.ProductionOrderStageDeviceID,
		&rcv.DeviceID,
		&rcv.ProcessStatus,
		&rcv.CreatedAt,
	}

	return
}

func (*DeviceProgressStatusHistory) TableName() string {
	return "device_progress_status_history"
}
