package model

import (
	"database/sql"
	"time"
)

const (
	EmployeeFieldID          = "id"
	EmployeeFieldName        = "name"
	EmployeeFieldAvatar      = "avatar"
	EmployeeFieldPhoneNumber = "phone_number"
	EmployeeFieldEmail       = "email"
	EmployeeFieldStatus      = "status"
	EmployeeFieldType        = "type"
	EmployeeFieldAddress     = "address"
	EmployeeFieldCreatedBy   = "created_by"
	EmployeeFieldCreatedAt   = "created_at"
	EmployeeFieldUpdatedAt   = "updated_at"
	EmployeeFieldDeletedAt   = "deleted_at"
)

type Employee struct {
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

func (rcv *Employee) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		EmployeeFieldID,
		EmployeeFieldName,
		EmployeeFieldAvatar,
		EmployeeFieldPhoneNumber,
		EmployeeFieldEmail,
		EmployeeFieldStatus,
		EmployeeFieldType,
		EmployeeFieldAddress,
		EmployeeFieldCreatedBy,
		EmployeeFieldCreatedAt,
		EmployeeFieldUpdatedAt,
		EmployeeFieldDeletedAt,
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

func (*Employee) TableName() string {
	return "employees"
}
