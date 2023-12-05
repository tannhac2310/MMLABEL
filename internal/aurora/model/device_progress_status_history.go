package model

import (
	"database/sql"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"time"
)

const (
	DeviceProgressStatusHistoryFieldID                           = "id"
	DeviceProgressStatusHistoryFieldProductionOrderStageDeviceID = "production_order_stage_device_id"
	DeviceProgressStatusHistoryFieldDeviceID                     = "device_id"
	DeviceProgressStatusHistoryFieldProcessStatus                = "process_status"
	DeviceProgressStatusHistoryFieldIsResolved                   = "is_resolved"
	DeviceProgressStatusHistoryFieldUpdatedAt                    = "updated_at"
	DeviceProgressStatusHistoryFieldUpdatedBy                    = "updated_by"
	DeviceProgressStatusHistoryFieldErrorCode                    = "error_code"
	DeviceProgressStatusHistoryFieldErrorReason                  = "error_reason"
	DeviceProgressStatusHistoryFieldDescription                  = "description"
	DeviceProgressStatusHistoryFieldCreatedAt                    = "created_at"
)

type DeviceProgressStatusHistory struct {
	ID                           string                                `db:"id"`
	ProductionOrderStageDeviceID string                                `db:"production_order_stage_device_id"`
	DeviceID                     string                                `db:"device_id"`
	ProcessStatus                enum.ProductionOrderStageDeviceStatus `db:"process_status"`
	IsResolved                   int16                                 `db:"is_resolved"`
	UpdatedAt                    sql.NullTime                          `db:"updated_at"`
	UpdatedBy                    string                                `db:"updated_by"`
	ErrorCode                    sql.NullString                        `db:"error_code"`
	ErrorReason                  sql.NullString                        `db:"error_reason"`
	Description                  sql.NullString                        `db:"description"`
	CreatedAt                    time.Time                             `db:"created_at"`
}

func (rcv *DeviceProgressStatusHistory) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		DeviceProgressStatusHistoryFieldID,
		DeviceProgressStatusHistoryFieldProductionOrderStageDeviceID,
		DeviceProgressStatusHistoryFieldDeviceID,
		DeviceProgressStatusHistoryFieldProcessStatus,
		DeviceProgressStatusHistoryFieldIsResolved,
		DeviceProgressStatusHistoryFieldUpdatedAt,
		DeviceProgressStatusHistoryFieldUpdatedBy,
		DeviceProgressStatusHistoryFieldErrorCode,
		DeviceProgressStatusHistoryFieldErrorReason,
		DeviceProgressStatusHistoryFieldDescription,
		DeviceProgressStatusHistoryFieldCreatedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.ProductionOrderStageDeviceID,
		&rcv.DeviceID,
		&rcv.ProcessStatus,
		&rcv.IsResolved,
		&rcv.UpdatedAt,
		&rcv.UpdatedBy,
		&rcv.ErrorCode,
		&rcv.ErrorReason,
		&rcv.Description,
		&rcv.CreatedAt,
	}

	return
}

func (*DeviceProgressStatusHistory) TableName() string {
	return "device_progress_status_history"
}
