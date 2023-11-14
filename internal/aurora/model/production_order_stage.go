package model

import (
	"database/sql"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"time"
)

const (
	ProductionOrderStageFieldID                  = "id"
	ProductionOrderStageFieldSorting             = "sorting"
	ProductionOrderStageFieldProductionOrderID   = "production_order_id"
	ProductionOrderStageFieldStageID             = "stage_id"
	ProductionOrderStageFieldEstimatedStartAt    = "estimated_start_at"
	ProductionOrderStageFieldEstimatedCompleteAt = "estimated_complete_at"
	ProductionOrderStageFieldStartedAt           = "started_at"
	ProductionOrderStageFieldCompletedAt         = "completed_at"
	ProductionOrderStageFieldStatus              = "status"
	ProductionOrderStageFieldCondition           = "condition"
	ProductionOrderStageFieldNote                = "note"
	ProductionOrderStageFieldData                = "data"
	ProductionOrderStageFieldCreatedAt           = "created_at"
	ProductionOrderStageFieldUpdatedAt           = "updated_at"
	ProductionOrderStageFieldDeletedAt           = "deleted_at"
)

type ProductionOrderStage struct {
	ID                  string                          `db:"id"`
	Sorting             int16                           `db:"sorting"`
	ProductionOrderID   string                          `db:"production_order_id"`
	StageID             string                          `db:"stage_id"`
	EstimatedStartAt    sql.NullTime                    `db:"estimated_start_at"`
	EstimatedCompleteAt sql.NullTime                    `db:"estimated_complete_at"`
	StartedAt           sql.NullTime                    `db:"started_at"`
	CompletedAt         sql.NullTime                    `db:"completed_at"`
	Status              enum.ProductionOrderStageStatus `db:"status"`
	Condition           sql.NullString                  `db:"condition"`
	Note                sql.NullString                  `db:"note"`
	Data                map[string]interface{}          `db:"data"`
	CreatedAt           time.Time                       `db:"created_at"`
	UpdatedAt           time.Time                       `db:"updated_at"`
	DeletedAt           sql.NullTime                    `db:"deleted_at"`
}

func (rcv *ProductionOrderStage) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		ProductionOrderStageFieldID,
		ProductionOrderStageFieldSorting,
		ProductionOrderStageFieldProductionOrderID,
		ProductionOrderStageFieldStageID,
		ProductionOrderStageFieldEstimatedStartAt,
		ProductionOrderStageFieldEstimatedCompleteAt,
		ProductionOrderStageFieldStartedAt,
		ProductionOrderStageFieldCompletedAt,
		ProductionOrderStageFieldStatus,
		ProductionOrderStageFieldCondition,
		ProductionOrderStageFieldNote,
		ProductionOrderStageFieldData,
		ProductionOrderStageFieldCreatedAt,
		ProductionOrderStageFieldUpdatedAt,
		ProductionOrderStageFieldDeletedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.Sorting,
		&rcv.ProductionOrderID,
		&rcv.StageID,
		&rcv.EstimatedStartAt,
		&rcv.EstimatedCompleteAt,
		&rcv.StartedAt,
		&rcv.CompletedAt,
		&rcv.Status,
		&rcv.Condition,
		&rcv.Note,
		&rcv.Data,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
		&rcv.DeletedAt,
	}

	return
}

func (*ProductionOrderStage) TableName() string {
	return "production_order_stages"
}
