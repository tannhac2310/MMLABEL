package model

import (
	"database/sql"
	"time"
)

const (
	UserGroupFieldID        = "id"
	UserGroupFieldUserID    = "user_id"
	UserGroupFieldGroupID   = "group_id"
	UserGroupFieldCreatedAt = "created_at"
	UserGroupFieldUpdatedAt = "updated_at"
	UserGroupFieldDeletedAt = "deleted_at"
)

type UserGroup struct {
	ID        string       `db:"id"`
	UserID    string       `db:"user_id"`
	GroupID   string       `db:"group_id"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}

func (*UserGroup) TableName() string {
	return "user_group"
}

func (u *UserGroup) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		UserGroupFieldID,
		UserGroupFieldUserID,
		UserGroupFieldGroupID,
		UserGroupFieldCreatedAt,
		UserGroupFieldUpdatedAt,
		UserGroupFieldDeletedAt,
	}

	values = []interface{}{
		&u.ID,
		&u.UserID,
		&u.GroupID,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.DeletedAt,
	}

	return
}
