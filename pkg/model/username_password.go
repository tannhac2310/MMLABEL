package model

import (
	"database/sql"
	"time"
)

const (
	UserNamePasswordFieldID          = "id"
	UserNamePasswordFieldUserID      = "user_id"
	UserNamePasswordFieldEmail       = "email"
	UserNamePasswordFieldPhoneNumber = "phone_number"
	UserNamePasswordFieldPassword    = "password"
	UserNamePasswordFieldCreatedAt   = "created_at"
	UserNamePasswordFieldUpdatedAt   = "updated_at"
	UserNamePasswordFieldDeletedAt   = "deleted_at"
)

type UserNamePassword struct {
	ID          string         `db:"id"`
	UserID      string         `db:"user_id"`
	Email       sql.NullString `db:"email,omitempty"`
	PhoneNumber sql.NullString `db:"phone_number"`
	Password    string         `db:"password"`
	CreatedAt   time.Time      `db:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at"`
	DeletedAt   sql.NullTime   `db:"deleted_at"`
}

func (*UserNamePassword) TableName() string {
	return "username_passwords"
}

func (u *UserNamePassword) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		UserNamePasswordFieldID,
		UserNamePasswordFieldUserID,
		UserNamePasswordFieldEmail,
		UserNamePasswordFieldPhoneNumber,
		UserNamePasswordFieldPassword,
		UserNamePasswordFieldCreatedAt,
		UserNamePasswordFieldUpdatedAt,
		UserNamePasswordFieldDeletedAt,
	}

	values = []interface{}{
		&u.ID,
		&u.UserID,
		&u.Email,
		&u.PhoneNumber,
		&u.Password,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.DeletedAt,
	}

	return
}
