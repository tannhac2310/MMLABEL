package model

import (
	"database/sql"
	"time"
)

const (
	ProductionPlanAttributeFieldID             = "id"
	ProductionPlanAttributeFieldKind           = "kind"
	ProductionPlanAttributeFieldDisplayName    = "display_name"
	ProductionPlanAttributeFieldAttributeValue = "attribute_value"
	ProductionPlanAttributeFieldNote           = "note"
	ProductionPlanAttributeFieldData           = "data"
	ProductionPlanAttributeFieldStatus         = "status"
	ProductionPlanAttributeFieldCreatedBy      = "created_by"
	ProductionPlanAttributeFieldCreatedAt      = "created_at"
	ProductionPlanAttributeFieldUpdatedBy      = "updated_by"
	ProductionPlanAttributeFieldUpdatedAt      = "updated_at"
	ProductionPlanAttributeFieldDeletedAt      = "deleted_at"
)

type ProductionPlanAttribute struct {
	ID             string                 `db:"id"`
	Kind           int16                  `db:"kind"`
	DisplayName    string                 `db:"display_name"`
	AttributeValue string                 `db:"attribute_value"`
	Note           sql.NullString         `db:"note"`
	Data           map[string]interface{} `db:"data"`
	Status         int16                  `db:"status"`
	CreatedBy      string                 `db:"created_by"`
	CreatedAt      time.Time              `db:"created_at"`
	UpdatedBy      string                 `db:"updated_by"`
	UpdatedAt      time.Time              `db:"updated_at"`
	DeletedAt      sql.NullTime           `db:"deleted_at"`
}

func (rcv *ProductionPlanAttribute) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		ProductionPlanAttributeFieldID,
		ProductionPlanAttributeFieldKind,
		ProductionPlanAttributeFieldDisplayName,
		ProductionPlanAttributeFieldAttributeValue,
		ProductionPlanAttributeFieldNote,
		ProductionPlanAttributeFieldData,
		ProductionPlanAttributeFieldStatus,
		ProductionPlanAttributeFieldCreatedBy,
		ProductionPlanAttributeFieldCreatedAt,
		ProductionPlanAttributeFieldUpdatedBy,
		ProductionPlanAttributeFieldUpdatedAt,
		ProductionPlanAttributeFieldDeletedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.Kind,
		&rcv.DisplayName,
		&rcv.AttributeValue,
		&rcv.Note,
		&rcv.Data,
		&rcv.Status,
		&rcv.CreatedBy,
		&rcv.CreatedAt,
		&rcv.UpdatedBy,
		&rcv.UpdatedAt,
		&rcv.DeletedAt,
	}

	return
}

func (*ProductionPlanAttribute) TableName() string {
	return "production_plan_attributes"
}
