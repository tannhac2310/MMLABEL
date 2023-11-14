package model

import (
	"database/sql"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"time"
)

const (
	DeviceFieldID        = "id"
	DeviceFieldName      = "name"
	DeviceFieldCode      = "code"
	DeviceFieldOptionID  = "option_id"
	DeviceFieldData      = "data"
	DeviceFieldStatus    = "status"
	DeviceFieldCreatedBy = "created_by"
	DeviceFieldCreatedAt = "created_at"
	DeviceFieldUpdatedAt = "updated_at"
	DeviceFieldDeletedAt = "deleted_at"
)

type Device struct {
	ID        string                 `db:"id"`
	Name      string                 `db:"name"`
	Code      string                 `db:"code"`
	OptionID  string                 `db:"option_id"`
	Data      map[string]interface{} `db:"data"`
	Status    enum.CommonStatus      `db:"status"`
	CreatedBy string                 `db:"created_by"`
	CreatedAt time.Time              `db:"created_at"`
	UpdatedAt time.Time              `db:"updated_at"`
	DeletedAt sql.NullTime           `db:"deleted_at"`
}

func (rcv *Device) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		DeviceFieldID,
		DeviceFieldName,
		DeviceFieldCode,
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
