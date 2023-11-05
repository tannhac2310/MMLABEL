package model

import (
	"database/sql"
	"time"
)

const (
	CustomerFieldID          = "id"
	CustomerFieldName        = "name"
	CustomerFieldAvatar      = "avatar"
	CustomerFieldPhoneNumber = "phone_number"
	CustomerFieldEmail       = "email"
	CustomerFieldStatus      = "status"
	CustomerFieldType        = "type"
	CustomerFieldAddress     = "address"
	CustomerFieldCreatedBy   = "created_by"
	CustomerFieldCreatedAt   = "created_at"
	CustomerFieldUpdatedAt   = "updated_at"
	CustomerFieldDeletedAt   = "deleted_at"
)

type Customer struct {
	ID          string         `db:"id"`
	Name        string         `db:"name"`
	Avatar      sql.NullString `db:"avatar"`
	PhoneNumber sql.NullString `db:"phone_number"`
	Email       sql.NullString `db:"email"`
	Status      int16          `db:"status"`
	Type        int16          `db:"type"`
	Address     sql.NullString `db:"address"`
	CreatedBy   string         `db:"created_by"`
	CreatedAt   time.Time      `db:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at"`
	DeletedAt   sql.NullTime   `db:"deleted_at"`
}

func (rcv *Customer) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		CustomerFieldID,
		CustomerFieldName,
		CustomerFieldAvatar,
		CustomerFieldPhoneNumber,
		CustomerFieldEmail,
		CustomerFieldStatus,
		CustomerFieldType,
		CustomerFieldAddress,
		CustomerFieldCreatedBy,
		CustomerFieldCreatedAt,
		CustomerFieldUpdatedAt,
		CustomerFieldDeletedAt,
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
		&rcv.CreatedBy,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
		&rcv.DeletedAt,
	}

	return
}

func (*Customer) TableName() string {
	return "customers"
}
