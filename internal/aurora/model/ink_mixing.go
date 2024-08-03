package model

import (
	"database/sql"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

const (
	InkMixingFieldID          = "id"
	InkMixingFieldName        = "name"
	InkMixingFieldCode        = "code"
	InkMixingFieldInkID       = "ink_id"
	InkMixingFieldDescription = "description"
	InkMixingFieldStatus      = "status"
	InkMixingFieldData        = "data"
	InkMixingFieldCreatedBy   = "created_by"
	InkMixingFieldUpdatedBy   = "updated_by"
	InkMixingFieldCreatedAt   = "created_at"
	InkMixingFieldUpdatedAt   = "updated_at"
	InkMixingFieldDeletedAt   = "deleted_at"
)

type InkMixing struct {
	ID          string                 `db:"id"`
	Name        string                 `db:"name"`
	Code        string                 `db:"code"`
	InkID       string                 `db:"ink_id"`
	Description string                 `db:"description"`
	Status      enum.CommonStatus      `db:"status"`
	Data        map[string]interface{} `db:"data"`
	CreatedBy   string                 `db:"created_by"`
	UpdatedBy   string                 `db:"updated_by"`
	CreatedAt   time.Time              `db:"created_at"`
	UpdatedAt   time.Time              `db:"updated_at"`
	DeletedAt   sql.NullTime           `db:"deleted_at"`
}

func (rcv *InkMixing) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		InkMixingFieldID,
		InkMixingFieldName,
		InkMixingFieldCode,
		InkMixingFieldInkID,
		InkMixingFieldDescription,
		InkMixingFieldStatus,
		InkMixingFieldData,
		InkMixingFieldCreatedBy,
		InkMixingFieldUpdatedBy,
		InkMixingFieldCreatedAt,
		InkMixingFieldUpdatedAt,
		InkMixingFieldDeletedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.Name,
		&rcv.Code,
		&rcv.InkID,
		&rcv.Description,
		&rcv.Status,
		&rcv.Data,
		&rcv.CreatedBy,
		&rcv.UpdatedBy,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
		&rcv.DeletedAt,
	}

	return
}

func (*InkMixing) TableName() string {
	return "ink_mixing"
}
