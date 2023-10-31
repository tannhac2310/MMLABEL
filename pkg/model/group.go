package model

import (
	"database/sql"
	"time"
)

const (
	GroupFieldID        = "id"
	GroupFieldName      = "name"
	GroupFieldCreatedAt = "created_at"
	GroupFieldUpdatedAt = "updated_at"
	GroupFieldDeletedAt = "deleted_at"
)

type Group struct {
	ID        string       `db:"id"`
	Name      string       `db:"name"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}

func (*Group) TableName() string {
	return "groups"
}

func (g *Group) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		GroupFieldID,
		GroupFieldName,
		GroupFieldCreatedAt,
		GroupFieldUpdatedAt,
		GroupFieldDeletedAt,
	}

	values = []interface{}{
		&g.ID,
		&g.Name,
		&g.CreatedAt,
		&g.UpdatedAt,
		&g.DeletedAt,
	}

	return
}
