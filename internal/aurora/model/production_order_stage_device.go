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
	ProductionOrderStageDeviceFieldResponsible            = "responsible"
	ProductionOrderStageDeviceFieldSettings               = "settings"
	ProductionOrderStageDeviceFieldNote                   = "note"
	ProductionOrderStageDeviceFieldCreatedAt              = "created_at"
	ProductionOrderStageDeviceFieldUpdatedAt              = "updated_at"
	ProductionOrderStageDeviceFieldDeletedAt              = "deleted_at"
)

type ProductionOrderStageDevice struct {
	ID                     string                                `db:"id"`
	ProductionOrderStageID string                                `db:"production_order_stage_id"`
	DeviceID               string                                `db:"device_id"`
	Quantity               int64                                 `db:"quantity"`
	ProcessStatus          enum.ProductionOrderStageDeviceStatus `db:"process_status"`
	Status                 enum.CommonStatus                     `db:"status"`
	Responsible            []string                              `db:"responsible"`
	Settings               map[string]interface{}                `db:"settings"`
	Note                   sql.NullString                        `db:"note"`
	CreatedAt              time.Time                             `db:"created_at"`
	UpdatedAt              time.Time                             `db:"updated_at"`
	DeletedAt              sql.NullTime                          `db:"deleted_at"`
}

func (rcv *ProductionOrderStageDevice) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		ProductionOrderStageDeviceFieldID,
		ProductionOrderStageDeviceFieldProductionOrderStageID,
		ProductionOrderStageDeviceFieldDeviceID,
		ProductionOrderStageDeviceFieldQuantity,
		ProductionOrderStageDeviceFieldProcessStatus,
		ProductionOrderStageDeviceFieldStatus,
		ProductionOrderStageDeviceFieldResponsible,
		ProductionOrderStageDeviceFieldSettings,
		ProductionOrderStageDeviceFieldNote,
		ProductionOrderStageDeviceFieldCreatedAt,
		ProductionOrderStageDeviceFieldUpdatedAt,
		ProductionOrderStageDeviceFieldDeletedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.ProductionOrderStageID,
		&rcv.DeviceID,
		&rcv.Quantity,
		&rcv.ProcessStatus,
		&rcv.Status,
		&rcv.Responsible,
		&rcv.Settings,
		&rcv.Note,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
		&rcv.DeletedAt,
	}

	return
}

func (*ProductionOrderStageDevice) TableName() string {
	return "production_order_stage_devices"
}
