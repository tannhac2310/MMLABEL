package model

import (
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"time"
)

const (
	RolePermissionFieldID         = "id"
	RolePermissionFieldRoleID     = "role_id"
	RolePermissionFieldEntityType = "entity_type"
	RolePermissionFieldEntityID   = "entity_id"
	RolePermissionFieldCreatedAt  = "created_at"
)

type RolePermission struct {
	ID         string                    `db:"id"`
	RoleID     string                    `db:"role_id"`
	EntityType enum.PermissionEntityType `db:"entity_type"`
	EntityID   string                    `db:"entity_id"`
	CreatedAt  time.Time                 `db:"created_at"`
}

func (rcv *RolePermission) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		RolePermissionFieldID,
		RolePermissionFieldRoleID,
		RolePermissionFieldEntityType,
		RolePermissionFieldEntityID,
		RolePermissionFieldCreatedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.RoleID,
		&rcv.EntityType,
		&rcv.EntityID,
		&rcv.CreatedAt,
	}

	return
}

func (*RolePermission) TableName() string {
	return "role_permissions"
}
