package model

import (
	"database/sql"
	"time"
)

const (
	UserRoleUser = "user"
	UserRoleRoot = "root"
)

const (
	RoleFieldID        = "id"
	RoleFieldName      = "name"
	RoleFieldPriority  = "priority"
	RoleFieldCreatedAt = "created_at"
	RoleFieldUpdatedAt = "updated_at"
	RoleFieldDeletedAt = "deleted_at"
)

type Role struct {
	ID        string       `db:"id"`
	Name      string       `db:"name"`
	Priority  int          `db:"priority"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}

func (*Role) TableName() string {
	return "roles"
}

func (r *Role) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		RoleFieldID,
		RoleFieldName,
		RoleFieldPriority,
		RoleFieldCreatedAt,
		RoleFieldUpdatedAt,
		RoleFieldDeletedAt,
	}

	values = []interface{}{
		&r.ID,
		&r.Name,
		&r.Priority,
		&r.CreatedAt,
		&r.UpdatedAt,
		&r.DeletedAt,
	}

	return
}
