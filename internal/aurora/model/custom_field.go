package model

import (
	"database/sql"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

const (
	CustomFieldFieldID          = "id"
	CustomFieldFieldEntityID    = "entity_id"
	CustomFieldFieldEntityType  = "entity_type"
	CustomFieldFieldField       = "field"
	CustomFieldFieldValue       = "value"
	CustomFieldFieldDescription = "description"
)

type CustomField struct {
	ID          string               `db:"id"`
	EntityID    string               `db:"entity_id"`
	EntityType  enum.CustomFieldType `db:"entity_type"`
	Field       string               `db:"field"`
	Value       string               `db:"value"`
	Description sql.NullString       `db:"description"`
}

func (rcv *CustomField) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		CustomFieldFieldID,
		CustomFieldFieldEntityID,
		CustomFieldFieldEntityType,
		CustomFieldFieldField,
		CustomFieldFieldValue,
		CustomFieldFieldDescription,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.EntityID,
		&rcv.EntityType,
		&rcv.Field,
		&rcv.Value,
		&rcv.Description,
	}

	return
}

func (*CustomField) TableName() string {
	return "custom_fields"
}
