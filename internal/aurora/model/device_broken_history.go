package model

import (
	"database/sql"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

const (
	DeviceBrokenHistoryFieldID                           = "id"
	DeviceBrokenHistoryFieldProductionOrderStageDeviceID = "production_order_stage_device_id"
	DeviceBrokenHistoryFieldDeviceID                     = "device_id"
	DeviceBrokenHistoryFieldProcessStatus                = "process_status"
	DeviceBrokenHistoryFieldCreatedAt                    = "created_at"
	DeviceBrokenHistoryFieldIsResolved                   = "is_resolved"
	DeviceBrokenHistoryFieldUpdatedAt                    = "updated_at"
	DeviceBrokenHistoryFieldUpdatedBy                    = "updated_by"
	DeviceBrokenHistoryFieldCreatedBy                    = "created_by"
	DeviceBrokenHistoryFieldErrorCode                    = "error_code"
	DeviceBrokenHistoryFieldErrorReason                  = "error_reason"
	DeviceBrokenHistoryFieldDescription                  = "description"
)

type DeviceBrokenHistory struct {
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
type DeviceBrokenHistoryUpdateIsSolved struct {
	ID string `db:"id"`
}

func (rcv *DeviceBrokenHistory) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		DeviceBrokenHistoryFieldID,
		DeviceBrokenHistoryFieldProductionOrderStageDeviceID,
		DeviceBrokenHistoryFieldDeviceID,
		DeviceBrokenHistoryFieldProcessStatus,
		DeviceBrokenHistoryFieldCreatedAt,
		DeviceBrokenHistoryFieldIsResolved,
		DeviceBrokenHistoryFieldUpdatedAt,
		DeviceBrokenHistoryFieldUpdatedBy,
		DeviceBrokenHistoryFieldCreatedBy,
		DeviceBrokenHistoryFieldErrorCode,
		DeviceBrokenHistoryFieldErrorReason,
		DeviceBrokenHistoryFieldDescription,
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

func (*DeviceBrokenHistory) TableName() string {
	return "device_broken_history"
}
