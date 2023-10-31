package model

import (
	"database/sql"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

const (
	PermissionFieldID        = "id"
	PermissionFieldUserID    = "user_id"
	PermissionFieldEntity    = "entity"
	PermissionFieldElementID = "element_id"
	PermissionFieldCreatedBy = "created_by"
	PermissionFieldUpdatedBy = "updated_by"
	PermissionFieldCreatedAt = "created_at"
	PermissionFieldUpdatedAt = "updated_at"
	PermissionFieldDeletedAt = "deleted_at"
)

type Permission struct {
	ID        string                `db:"id"`
	UserID    string                `db:"user_id"`
	Entity    enum.PermissionEntity `db:"entity"`
	ElementID string                `db:"element_id"`
	CreatedBy string                `db:"created_by"`
	UpdatedBy sql.NullString        `db:"updated_by"`
	CreatedAt time.Time             `db:"created_at"`
	UpdatedAt time.Time             `db:"updated_at"`
	DeletedAt sql.NullTime          `db:"deleted_at"`
}

func (rcv *Permission) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		PermissionFieldID,
		PermissionFieldUserID,
		PermissionFieldEntity,
		PermissionFieldElementID,
		PermissionFieldCreatedBy,
		PermissionFieldUpdatedBy,
		PermissionFieldCreatedAt,
		PermissionFieldUpdatedAt,
		PermissionFieldDeletedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.UserID,
		&rcv.Entity,
		&rcv.ElementID,
		&rcv.CreatedBy,
		&rcv.UpdatedBy,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
		&rcv.DeletedAt,
	}

	return
}

func (*Permission) TableName() string {
	return "permissions"
}
