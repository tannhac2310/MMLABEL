package model

import (
	"database/sql"
	"time"
)

const (
	UserFCMTokenFieldID        = "id"
	UserFCMTokenFieldUserID    = "user_id"
	UserFCMTokenFieldDeviceID  = "device_id"
	UserFCMTokenFieldToken     = "token"
	UserFCMTokenFieldCreatedAt = "created_at"
	UserFCMTokenFieldUpdatedAt = "updated_at"
	UserFCMTokenFieldDeletedAt = "deleted_at"
)

type UserFCMToken struct {
	ID        string       `db:"id"`
	UserID    string       `db:"user_id"`
	DeviceID  string       `db:"device_id"`
	Token     string       `db:"token"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}

func (*UserFCMToken) TableName() string {
	return "user_fcm_tokens"
}

func (u *UserFCMToken) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		UserFCMTokenFieldID,
		UserFCMTokenFieldUserID,
		UserFCMTokenFieldDeviceID,
		UserFCMTokenFieldToken,
		UserFCMTokenFieldCreatedAt,
		UserFCMTokenFieldUpdatedAt,
		UserFCMTokenFieldDeletedAt,
	}

	values = []interface{}{
		&u.ID,
		&u.UserID,
		&u.DeviceID,
		&u.Token,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.DeletedAt,
	}

	return
}
