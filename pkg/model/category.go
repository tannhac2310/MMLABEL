package model

import (
	"database/sql"
	"time"
)

const (
	CategoryFieldID          = "id"
	CategoryFieldName        = "name"
	CategoryFieldDescription = "description"
	CategoryFieldCreatedAt   = "created_at"
	CategoryFieldCreatedBy   = "created_by"
	CategoryFieldUpdatedAt   = "updated_at"
	CategoryFieldUpdatedBy   = "updated_by"
	CategoryFieldDeletedAt   = "deleted_at"
)

type Category struct {
	ID          string       `db:"id"`
	Name        string       `db:"name"`
	Description string       `db:"description"`
	CreatedAt   time.Time    `db:"created_at"`
	CreatedBy   string       `db:"created_by"`
	UpdatedBy   string       `db:"updated_by"`
	UpdatedAt   time.Time    `db:"updated_at"`
	DeletedAt   sql.NullTime `db:"deleted_at"`
}

func (b *Category) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		CategoryFieldID,
		CategoryFieldName,
		CategoryFieldDescription,
		CategoryFieldCreatedAt,
		CategoryFieldUpdatedAt,
		CategoryFieldDeletedAt,
		CategoryFieldCreatedBy,
		CategoryFieldUpdatedBy,
	}

	values = []interface{}{
		&b.ID,
		&b.Name,
		&b.Description,
		&b.CreatedAt,
		&b.UpdatedAt,
		&b.DeletedAt,
		&b.CreatedBy,
		&b.UpdatedBy,
	}
	return
}

func (*Category) TableName() string {
	return "category"
}
