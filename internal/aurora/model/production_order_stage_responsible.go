package model

import (
	"database/sql"
	"time"
)

const (
	ProductionOrderStageResponsibleFieldID                     = "id"
	ProductionOrderStageResponsibleFieldProductionOrderStageID = "production_order_stage_id"
	ProductionOrderStageResponsibleFieldUserID                 = "user_id"
	ProductionOrderStageResponsibleFieldCreatedAt              = "created_at"
	ProductionOrderStageResponsibleFieldUpdatedAt              = "updated_at"
	ProductionOrderStageResponsibleFieldDeletedAt              = "deleted_at"
)

type ProductionOrderStageResponsible struct {
	ID                     string       `db:"id"`
	ProductionOrderStageID string       `db:"production_order_stage_id"`
	UserID                 string       `db:"user_id"`
	CreatedAt              time.Time    `db:"created_at"`
	UpdatedAt              time.Time    `db:"updated_at"`
	DeletedAt              sql.NullTime `db:"deleted_at"`
}

func (rcv *ProductionOrderStageResponsible) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		ProductionOrderStageResponsibleFieldID,
		ProductionOrderStageResponsibleFieldProductionOrderStageID,
		ProductionOrderStageResponsibleFieldUserID,
		ProductionOrderStageResponsibleFieldCreatedAt,
		ProductionOrderStageResponsibleFieldUpdatedAt,
		ProductionOrderStageResponsibleFieldDeletedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.ProductionOrderStageID,
		&rcv.UserID,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
		&rcv.DeletedAt,
	}

	return
}

func (*ProductionOrderStageResponsible) TableName() string {
	return "production_order_stage_responsible"
}
