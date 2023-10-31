package model

import (
	"database/sql"
	"time"
)

const (
	UserDeviceFieldID        = "id"
	UserDeviceFieldUserID    = "user_id"
	UserDeviceFieldDeviceID  = "device_id"
	UserDeviceFieldCreatedAt = "created_at"
	UserDeviceFieldUpdatedAt = "updated_at"
	UserDeviceFieldDeletedAt = "deleted_at"
)

type UserDevice struct {
	ID        string       `db:"id"`
	UserID    string       `db:"user_id"`
	DeviceID  string       `db:"device_id"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}

func (d *UserDevice) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		UserDeviceFieldID,
		UserDeviceFieldUserID,
		UserDeviceFieldDeviceID,
		UserDeviceFieldCreatedAt,
		UserDeviceFieldUpdatedAt,
		UserDeviceFieldDeletedAt,
	}
	values = []interface{}{
		&d.ID,
		&d.UserID,
		&d.DeviceID,
		&d.CreatedAt,
		&d.UpdatedAt,
		&d.DeletedAt,
	}

	return
}

func (*UserDevice) TableName() string {
	return "user_device"
}
