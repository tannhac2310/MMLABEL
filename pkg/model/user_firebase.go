package model

import (
	"database/sql"
	"time"
)

const (
	UserFirebaseFieldID        = "id"
	UserFirebaseFieldUserID    = "user_id"
	UserFirebaseFieldCreatedAt = "created_at"
	UserFirebaseFieldUpdatedAt = "updated_at"
	UserFirebaseFieldDeletedAt = "deleted_at"
)

type UserFirebase struct {
	ID        string       `db:"id"`
	UserID    string       `db:"user_id"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}

func (*UserFirebase) TableName() string {
	return "user_firebases"
}

func (u *UserFirebase) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		UserFirebaseFieldID,
		UserFirebaseFieldUserID,
		UserFirebaseFieldCreatedAt,
		UserFirebaseFieldUpdatedAt,
		UserFirebaseFieldDeletedAt,
	}

	values = []interface{}{
		&u.ID,
		&u.UserID,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.DeletedAt,
	}

	return
}
