package model

import (
	"database/sql"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

const (
	DeviceProgressStatusHistoryFieldID                           = "id"
	DeviceProgressStatusHistoryFieldProductionOrderStageDeviceID = "production_order_stage_device_id"
	DeviceProgressStatusHistoryFieldDeviceID                     = "device_id"
	DeviceProgressStatusHistoryFieldProcessStatus                = "process_status"
	DeviceProgressStatusHistoryFieldCreatedAt                    = "created_at"
	DeviceProgressStatusHistoryFieldIsResolved                   = "is_resolved"
	DeviceProgressStatusHistoryFieldUpdatedAt                    = "updated_at"
	DeviceProgressStatusHistoryFieldUpdatedBy                    = "updated_by"
	DeviceProgressStatusHistoryFieldCreatedBy                    = "created_by"
	DeviceProgressStatusHistoryFieldErrorCode                    = "error_code"
	DeviceProgressStatusHistoryFieldErrorReason                  = "error_reason"
	DeviceProgressStatusHistoryFieldDescription                  = "description"
)

type DeviceProgressStatusHistory struct {
	ID                           string                                `db:"id"`
	ProductionOrderStageDeviceID string                                `db:"production_order_stage_device_id"`
	DeviceID                     string                                `db:"device_id"`
	ProcessStatus                enum.ProductionOrderStageDeviceStatus `db:"process_status"`
	CreatedAt                    time.Time                             `db:"created_at"`
	IsResolved                   int16                                 `db:"is_resolved"`
	UpdatedAt                    sql.NullTime                          `db:"updated_at"`
	UpdatedBy                    sql.NullString                        `db:"updated_by"`
	CreatedBy                    sql.NullString                        `db:"created_by"`
	ErrorCode                    sql.NullString                        `db:"error_code"`
	ErrorReason                  sql.NullString                        `db:"error_reason"`
	Description                  sql.NullString                        `db:"description"`
}
type DeviceProgressStatusHistoryUpdateIsSolved struct {
	ID                           string                                `db:"id"`
}

func (rcv *DeviceProgressStatusHistory) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		DeviceProgressStatusHistoryFieldID,
		DeviceProgressStatusHistoryFieldProductionOrderStageDeviceID,
		DeviceProgressStatusHistoryFieldDeviceID,
		DeviceProgressStatusHistoryFieldProcessStatus,
		DeviceProgressStatusHistoryFieldCreatedAt,
		DeviceProgressStatusHistoryFieldIsResolved,
		DeviceProgressStatusHistoryFieldUpdatedAt,
		DeviceProgressStatusHistoryFieldUpdatedBy,
		DeviceProgressStatusHistoryFieldCreatedBy,
		DeviceProgressStatusHistoryFieldErrorCode,
		DeviceProgressStatusHistoryFieldErrorReason,
		DeviceProgressStatusHistoryFieldDescription,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.ProductionOrderStageDeviceID,
		&rcv.DeviceID,
		&rcv.ProcessStatus,
		&rcv.CreatedAt,
		&rcv.IsResolved,
		&rcv.UpdatedAt,
		&rcv.UpdatedBy,
		&rcv.CreatedBy,
		&rcv.ErrorCode,
		&rcv.ErrorReason,
		&rcv.Description,
	}

	return
}

func (*DeviceProgressStatusHistory) TableName() string {
	return "device_progress_status_history"
}
