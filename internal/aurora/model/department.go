package model

import (
	"database/sql"
	"time"
)

const (
	DepartmentFieldID        = "id"
	DepartmentFieldParentID  = "parent_id"
	DepartmentFieldName      = "name"
	DepartmentFieldShortName = "short_name"
	DepartmentFieldCode      = "code"
	DepartmentFieldPriority  = "priority"
	DepartmentFieldCreatedAt = "created_at"
	DepartmentFieldUpdatedAt = "updated_at"
	DepartmentFieldDeletedAt = "deleted_at"
)

type Department struct {
	ID        string         `db:"id"`
	ParentID  sql.NullString `db:"parent_id"`
	Name      string         `db:"name"`
	ShortName string         `db:"short_name"`
	Code      string         `db:"code"`
	Priority  int64          `db:"priority"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
	DeletedAt sql.NullTime   `db:"deleted_at"`
}

func (rcv *Department) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		DepartmentFieldID,
		DepartmentFieldParentID,
		DepartmentFieldName,
		DepartmentFieldShortName,
		DepartmentFieldCode,
		DepartmentFieldPriority,
		DepartmentFieldCreatedAt,
		DepartmentFieldUpdatedAt,
		DepartmentFieldDeletedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.ParentID,
		&rcv.Name,
		&rcv.ShortName,
		&rcv.Code,
		&rcv.Priority,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
		&rcv.DeletedAt,
	}

	return
}

func (*Department) TableName() string {
	return "departments"
}
