package model

import (
	"database/sql"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"time"
)

const (
	UserFieldID          = "id"
	UserFieldName        = "name"
	UserFieldCode        = "code"
	UserFieldAvatar      = "avatar"
	UserFieldPhoneNumber = "phone_number"
	UserFieldEmail       = "email"
	UserFieldLinked      = "linked"
	UserFieldDepartments = "departments"
	UserFieldStatus      = "status"
	UserFieldType        = "type"
	UserFieldLanguageID  = "language_id"
	UserFieldBlocked     = "blocked"
	UserFieldAddress     = "address"
	UserFieldCreatedAt   = "created_at"
	UserFieldUpdatedAt   = "updated_at"
	UserFieldDeletedAt   = "deleted_at"
)

type User struct {
	ID          string          `db:"id"`
	Name        string          `db:"name"`
	Code        string          `db:"code"`
	Avatar      string          `db:"avatar"`
	PhoneNumber string          `db:"phone_number"`
	Email       string          `db:"email"`
	Linked      sql.NullString  `db:"linked"`
	Departments sql.NullString  `db:"departments"`
	Status      enum.UserStatus `db:"status"`
	Type        enum.UserType   `db:"type"`
	LanguageID  int16           `db:"language_id"`
	Blocked     bool            `db:"blocked"`
	Address     string          `db:"address"`
	CreatedAt   time.Time       `db:"created_at"`
	UpdatedAt   time.Time       `db:"updated_at"`
	DeletedAt   sql.NullTime    `db:"deleted_at"`
}

func (rcv *User) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		UserFieldID,
		UserFieldName,
		UserFieldCode,
		UserFieldAvatar,
		UserFieldPhoneNumber,
		UserFieldEmail,
		UserFieldLinked,
		UserFieldDepartments,
		UserFieldStatus,
		UserFieldType,
		UserFieldLanguageID,
		UserFieldBlocked,
		UserFieldAddress,
		UserFieldCreatedAt,
		UserFieldUpdatedAt,
		UserFieldDeletedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.Name,
		&rcv.Code,
		&rcv.Avatar,
		&rcv.PhoneNumber,
		&rcv.Email,
		&rcv.Linked,
		&rcv.Departments,
		&rcv.Status,
		&rcv.Type,
		&rcv.LanguageID,
		&rcv.Blocked,
		&rcv.Address,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
		&rcv.DeletedAt,
	}

	return
}

func (*User) TableName() string {
	return "users"
}
