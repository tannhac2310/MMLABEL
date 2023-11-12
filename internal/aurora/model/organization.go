package model

import (
	"database/sql"
	"time"
)

const (
	OrganizationFieldID          = "id"
	OrganizationFieldName        = "name"
	OrganizationFieldAvatar      = "avatar"
	OrganizationFieldPhoneNumber = "phone_number"
	OrganizationFieldEmail       = "email"
	OrganizationFieldStatus      = "status"
	OrganizationFieldType        = "type"
	OrganizationFieldAddress     = "address"
	OrganizationFieldCreatedAt   = "created_at"
	OrganizationFieldUpdatedAt   = "updated_at"
	OrganizationFieldDeletedAt   = "deleted_at"
)

type Organization struct {
	ID          string         `db:"id"`
	Name        string         `db:"name"`
	Avatar      sql.NullString `db:"avatar"`
	PhoneNumber sql.NullString `db:"phone_number"`
	Email       sql.NullString `db:"email"`
	Status      int16          `db:"status"`
	Type        int16          `db:"type"`
	Address     sql.NullString `db:"address"`
	CreatedAt   time.Time      `db:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at"`
	DeletedAt   sql.NullTime   `db:"deleted_at"`
}

func (rcv *Organization) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		OrganizationFieldID,
		OrganizationFieldName,
		OrganizationFieldAvatar,
		OrganizationFieldPhoneNumber,
		OrganizationFieldEmail,
		OrganizationFieldStatus,
		OrganizationFieldType,
		OrganizationFieldAddress,
		OrganizationFieldCreatedAt,
		OrganizationFieldUpdatedAt,
		OrganizationFieldDeletedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.Name,
		&rcv.Avatar,
		&rcv.PhoneNumber,
		&rcv.Email,
		&rcv.Status,
		&rcv.Type,
		&rcv.Address,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
		&rcv.DeletedAt,
	}

	return
}

func (*Organization) TableName() string {
	return "organizations"
}
