package model

import (
	"database/sql"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

type UserLinked string

const (
	UserLinkedUserNamePassword = "email_password"
)

const (
	UserFieldID          = "id"
	UserFieldName        = "name"
	UserFieldAvatar      = "avatar"
	UserFieldPhoneNumber = "phone_number"
	UserFieldEmail       = "email"
	UserFieldLinked      = "linked"
	UserFieldStatus      = "status"
	UserFieldAddress     = "address"
	UserFieldType        = "type"
	UserFieldCreatedAt   = "created_at"
	UserFieldUpdatedAt   = "updated_at"
	UserFieldDeletedAt   = "deleted_at"
)

type User struct {
	ID          string          `db:"id"`
	Name        string          `db:"name"`
	Avatar      string          `db:"avatar"`
	Address     string          `db:"address"`
	PhoneNumber string          `db:"phone_number"`
	Email       string          `db:"email"`
	Linked      []string        `db:"linked"`
	Status      enum.UserStatus `db:"status"`
	Type        enum.UserType   `db:"type"`
	CreatedAt   time.Time       `db:"created_at"`
	UpdatedAt   time.Time       `db:"updated_at"`
	DeletedAt   sql.NullTime    `db:"deleted_at"`
}

func (*User) TableName() string {
	return "users"
}

func (u *User) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		UserFieldID,
		UserFieldName,
		UserFieldAvatar,
		UserFieldPhoneNumber,
		UserFieldEmail,
		UserFieldLinked,
		UserFieldStatus,
		UserFieldAddress,
		UserFieldType,
		UserFieldCreatedAt,
		UserFieldUpdatedAt,
		UserFieldDeletedAt,
	}

	values = []interface{}{
		&u.ID,
		&u.Name,
		&u.Avatar,
		&u.PhoneNumber,
		&u.Email,
		&u.Linked,
		&u.Status,
		&u.Address,
		&u.Type,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.DeletedAt,
	}

	return
}
