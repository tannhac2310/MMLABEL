package model

import (
	"database/sql"
	"time"
)

const (
	UserRoleFieldID        = "id"
	UserRoleFieldUserID    = "user_id"
	UserRoleFieldRoleID    = "role_id"
	UserRoleFieldCreatedAt = "created_at"
	UserRoleFieldCreatedBy = "created_by"
	UserRoleFieldUpdatedAt = "updated_at"
	UserRoleFieldDeletedAt = "deleted_at"
)

type UserRole struct {
	ID        string         `db:"id"`
	UserID    string         `db:"user_id"`
	RoleID    string         `db:"role_id"`
	CreatedBy sql.NullString `db:"created_by"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
	DeletedAt sql.NullTime   `db:"deleted_at"`
}

func (*UserRole) TableName() string {
	return "user_role"
}

func (u *UserRole) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		UserRoleFieldID,
		UserRoleFieldUserID,
		UserRoleFieldRoleID,
		UserRoleFieldCreatedAt,
		UserRoleFieldCreatedBy,
		UserRoleFieldUpdatedAt,
		UserRoleFieldDeletedAt,
	}

	values = []interface{}{
		&u.ID,
		&u.UserID,
		&u.RoleID,
		&u.CreatedAt,
		&u.CreatedBy,
		&u.UpdatedAt,
		&u.DeletedAt,
	}

	return
}
