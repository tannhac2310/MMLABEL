package model

import (
	"database/sql"
	"time"
)

const (
	OptionFieldID        = "id"
	OptionFieldEntity    = "entity"
	OptionFieldCode      = "code"
	OptionFieldName      = "name"
	OptionFieldData      = "data"
	OptionFieldStatus    = "status"
	OptionFieldCreatedBy = "created_by"
	OptionFieldCreatedAt = "created_at"
	OptionFieldUpdatedAt = "updated_at"
	OptionFieldDeletedAt = "deleted_at"
)

type Option struct {
	ID        string                 `db:"id"`
	Entity    string                 `db:"entity"`
	Code      string                 `db:"code"`
	Name      string                 `db:"name"`
	Data      map[string]interface{} `db:"data"`
	Status    int16                  `db:"status"`
	CreatedBy string                 `db:"created_by"`
	CreatedAt time.Time              `db:"created_at"`
	UpdatedAt time.Time              `db:"updated_at"`
	DeletedAt sql.NullTime           `db:"deleted_at"`
}

func (rcv *Option) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		OptionFieldID,
		OptionFieldEntity,
		OptionFieldCode,
		OptionFieldName,
		OptionFieldData,
		OptionFieldStatus,
		OptionFieldCreatedBy,
		OptionFieldCreatedAt,
		OptionFieldUpdatedAt,
		OptionFieldDeletedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.Entity,
		&rcv.Code,
		&rcv.Name,
		&rcv.Data,
		&rcv.Status,
		&rcv.CreatedBy,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
		&rcv.DeletedAt,
	}

	return
}

func (*Option) TableName() string {
	return "options"
}
