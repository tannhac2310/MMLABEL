package model

import (
	"database/sql"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

const (
	ProductionOrderStageDeviceFieldID                     = "id"
	ProductionOrderStageDeviceFieldProductionOrderStageID = "production_order_stage_id"
	ProductionOrderStageDeviceFieldDeviceID               = "device_id"
	ProductionOrderStageDeviceFieldQuantity               = "quantity"
	ProductionOrderStageDeviceFieldProcessStatus          = "process_status"
	ProductionOrderStageDeviceFieldStatus                 = "status"
	ProductionOrderStageDeviceFieldSettings               = "settings"
	ProductionOrderStageDeviceFieldNote                   = "note"
	ProductionOrderStageDeviceFieldCreatedAt              = "created_at"
	ProductionOrderStageDeviceFieldUpdatedAt              = "updated_at"
	ProductionOrderStageDeviceFieldDeletedAt              = "deleted_at"
	ProductionOrderStageDeviceFieldResponsible            = "responsible"
	ProductionOrderStageDeviceFieldEstimatedCompleteAt    = "estimated_complete_at"
	ProductionOrderStageDeviceFieldAssignedQuantity       = "assigned_quantity"
	ProductionOrderStageDeviceFieldEstimatedStartAt       = "estimated_start_at"
	ProductionOrderStageDeviceFieldColor                  = "color"
	ProductionOrderStageDeviceFieldStartAt                = "start_at"
	ProductionOrderStageDeviceFieldCompleteAt             = "complete_at"
)

type ProductionOrderStageDevice struct {
	ID                     string                                `db:"id"`
	ProductionOrderStageID string                                `db:"production_order_stage_id"`
	DeviceID               string                                `db:"device_id"`
	Quantity               int64                                 `db:"quantity"`
	ProcessStatus          enum.ProductionOrderStageDeviceStatus `db:"process_status"`
	Status                 enum.CommonStatus                     `db:"status"`
	Settings               map[string]interface{}                `db:"settings"`
	Note                   sql.NullString                        `db:"note"`
	CreatedAt              time.Time                             `db:"created_at"`
	UpdatedAt              time.Time                             `db:"updated_at"`
	DeletedAt              sql.NullTime                          `db:"deleted_at"`
	Responsible            []string                              `db:"responsible"`
	EstimatedCompleteAt    sql.NullTime                          `db:"estimated_complete_at"`
	AssignedQuantity       int64                                 `db:"assigned_quantity"`
	EstimatedStartAt       sql.NullTime                          `db:"estimated_start_at"`
	Color                  sql.NullString                        `db:"color"`
	StartAt                sql.NullTime                          `db:"start_at"`
	CompleteAt             sql.NullTime                          `db:"complete_at"`
}

func (rcv *ProductionOrderStageDevice) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		ProductionOrderStageDeviceFieldID,
		ProductionOrderStageDeviceFieldProductionOrderStageID,
		ProductionOrderStageDeviceFieldDeviceID,
		ProductionOrderStageDeviceFieldQuantity,
		ProductionOrderStageDeviceFieldProcessStatus,
		ProductionOrderStageDeviceFieldStatus,
		ProductionOrderStageDeviceFieldSettings,
		ProductionOrderStageDeviceFieldNote,
		ProductionOrderStageDeviceFieldCreatedAt,
		ProductionOrderStageDeviceFieldUpdatedAt,
		ProductionOrderStageDeviceFieldDeletedAt,
		ProductionOrderStageDeviceFieldResponsible,
		ProductionOrderStageDeviceFieldEstimatedCompleteAt,
		ProductionOrderStageDeviceFieldAssignedQuantity,
		ProductionOrderStageDeviceFieldEstimatedStartAt,
		ProductionOrderStageDeviceFieldColor,
		ProductionOrderStageDeviceFieldStartAt,
		ProductionOrderStageDeviceFieldCompleteAt,
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
		&rcv.Responsible,
		&rcv.EstimatedCompleteAt,
		&rcv.AssignedQuantity,
		&rcv.EstimatedStartAt,
		&rcv.Color,
		&rcv.StartAt,
		&rcv.CompleteAt,
	}

	return
}

func (*ProductionOrderStageDevice) TableName() string {
	return "production_order_stage_devices"
}
