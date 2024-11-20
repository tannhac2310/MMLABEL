package model

import (
	"database/sql"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

const (
	DeviceFieldID        = "id"
	DeviceFieldName      = "name"
	DeviceFieldStep      = "step"
	DeviceFieldCode      = "code"
	DeviceFieldSort      = "sort"
	DeviceFieldOptionID  = "option_id"
	DeviceFieldData      = "data"
	DeviceFieldStatus    = "status"
	DeviceFieldCreatedBy = "created_by"
	DeviceFieldCreatedAt = "created_at"
	DeviceFieldUpdatedAt = "updated_at"
	DeviceFieldDeletedAt = "deleted_at"
)

type SettingsData struct {
	DefectiveError               string `json:"defective_error"`
	DefectiveReason              string `json:"defective_reason"`
	Description                  string `json:"description"`
	ProductionOrderStageID       string `json:"production_order_stage_id"`
	ProductionOrderStageDeviceID string `json:"production_order_stage_device_id"`
}
type Device struct {
	ID        string            `db:"id"`
	Name      string            `db:"name"`
	Code      string            `db:"code"`
	Sort      int               `db:"sort"`
	Step      string            `db:"step"`
	OptionID  sql.NullString    `db:"option_id"`
	Data      SettingsData      `db:"data"`
	Status    enum.CommonStatus `db:"status"`
	CreatedBy string            `db:"created_by"`
	CreatedAt time.Time         `db:"created_at"`
	UpdatedAt time.Time         `db:"updated_at"`
	DeletedAt sql.NullTime      `db:"deleted_at"`
}

func (rcv *Device) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		DeviceFieldID,
		DeviceFieldName,
		DeviceFieldCode,
		DeviceFieldSort,
		DeviceFieldStep,
		DeviceFieldOptionID,
		DeviceFieldData,
		DeviceFieldStatus,
		DeviceFieldCreatedBy,
		DeviceFieldCreatedAt,
		DeviceFieldUpdatedAt,
		DeviceFieldDeletedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.Name,
		&rcv.Code,
		&rcv.Sort,
		&rcv.Step,
		&rcv.OptionID,
		&rcv.Data,
		&rcv.Status,
		&rcv.CreatedBy,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
		&rcv.DeletedAt,
	}

	return
}

func (*Device) TableName() string {
	return "devices"
}
