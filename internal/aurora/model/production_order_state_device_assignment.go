package model

import (
	"database/sql"
	"time"
)

const (
	ProductionOrderStateDeviceAssignmentFieldID                     = "id"
	ProductionOrderStateDeviceAssignmentFieldProductionOrderStageID = "production_order_stage_id"
	ProductionOrderStateDeviceAssignmentFieldDeviceID               = "device_id"
	ProductionOrderStateDeviceAssignmentFieldQuantity               = "quantity"
	ProductionOrderStateDeviceAssignmentFieldProcessStatus          = "process_status"
	ProductionOrderStateDeviceAssignmentFieldStatus                 = "status"
	ProductionOrderStateDeviceAssignmentFieldSettings               = "settings"
	ProductionOrderStateDeviceAssignmentFieldNote                   = "note"
	ProductionOrderStateDeviceAssignmentFieldCreatedAt              = "created_at"
	ProductionOrderStateDeviceAssignmentFieldUpdatedAt              = "updated_at"
	ProductionOrderStateDeviceAssignmentFieldDeletedAt              = "deleted_at"
)

type ProductionOrderStateDeviceAssignment struct {
	ID                     string                 `db:"id"`
	ProductionOrderStageID string                 `db:"production_order_stage_id"`
	DeviceID               string                 `db:"device_id"`
	Quantity               int64                  `db:"quantity"`
	ProcessStatus          int16                  `db:"process_status"`
	Status                 int16                  `db:"status"`
	Settings               map[string]interface{} `db:"settings"`
	Note                   sql.NullString         `db:"note"`
	CreatedAt              time.Time              `db:"created_at"`
	UpdatedAt              time.Time              `db:"updated_at"`
	DeletedAt              sql.NullTime           `db:"deleted_at"`
}

func (rcv *ProductionOrderStateDeviceAssignment) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		ProductionOrderStateDeviceAssignmentFieldID,
		ProductionOrderStateDeviceAssignmentFieldProductionOrderStageID,
		ProductionOrderStateDeviceAssignmentFieldDeviceID,
		ProductionOrderStateDeviceAssignmentFieldQuantity,
		ProductionOrderStateDeviceAssignmentFieldProcessStatus,
		ProductionOrderStateDeviceAssignmentFieldStatus,
		ProductionOrderStateDeviceAssignmentFieldSettings,
		ProductionOrderStateDeviceAssignmentFieldNote,
		ProductionOrderStateDeviceAssignmentFieldCreatedAt,
		ProductionOrderStateDeviceAssignmentFieldUpdatedAt,
		ProductionOrderStateDeviceAssignmentFieldDeletedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.ProductionOrderStageID,
		&rcv.DeviceID,
		&rcv.Quantity,
		&rcv.ProcessStatus,
		&rcv.Status,
		&rcv.Settings,
		&rcv.Note,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
		&rcv.DeletedAt,
	}

	return
}

func (*ProductionOrderStateDeviceAssignment) TableName() string {
	return "production_order_state_device_assignments"
}
