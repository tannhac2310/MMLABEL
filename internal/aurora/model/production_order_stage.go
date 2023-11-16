package model

import (
	"database/sql"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"time"
)

const (
	ProductionOrderStageFieldID                     = "id"
	ProductionOrderStageFieldProductionOrderID      = "production_order_id"
	ProductionOrderStageFieldStageID                = "stage_id"
	ProductionOrderStageFieldStartedAt              = "started_at"
	ProductionOrderStageFieldCompletedAt            = "completed_at"
	ProductionOrderStageFieldStatus                 = "status"
	ProductionOrderStageFieldCondition              = "condition"
	ProductionOrderStageFieldNote                   = "note"
	ProductionOrderStageFieldData                   = "data"
	ProductionOrderStageFieldCreatedAt              = "created_at"
	ProductionOrderStageFieldUpdatedAt              = "updated_at"
	ProductionOrderStageFieldWaitingAt              = "waiting_at"
	ProductionOrderStageFieldReceptionAt            = "reception_at"
	ProductionOrderStageFieldProductionStartAt      = "production_start_at"
	ProductionOrderStageFieldProductionCompletionAt = "production_completion_at"
	ProductionOrderStageFieldProductDeliveryAt      = "product_delivery_at"
	ProductionOrderStageFieldDeletedAt              = "deleted_at"
	ProductionOrderStageFieldEstimatedStartAt       = "estimated_start_at"
	ProductionOrderStageFieldEstimatedCompleteAt    = "estimated_complete_at"
	ProductionOrderStageFieldSorting                = "sorting"
)

type ProductionOrderStage struct {
	ID                     string                          `db:"id"`
	ProductionOrderID      string                          `db:"production_order_id"`
	StageID                string                          `db:"stage_id"`
	StartedAt              sql.NullTime                    `db:"started_at"`
	CompletedAt            sql.NullTime                    `db:"completed_at"`
	Status                 enum.ProductionOrderStageStatus `db:"status"`
	Condition              sql.NullString                  `db:"condition"`
	Note                   sql.NullString                  `db:"note"`
	Data                   map[string]interface{}          `db:"data"`
	CreatedAt              time.Time                       `db:"created_at"`
	UpdatedAt              time.Time                       `db:"updated_at"`
	WaitingAt              sql.NullTime                    `db:"waiting_at"`
	ReceptionAt            sql.NullTime                    `db:"reception_at"`
	ProductionStartAt      sql.NullTime                    `db:"production_start_at"`
	ProductionCompletionAt sql.NullTime                    `db:"production_completion_at"`
	ProductDeliveryAt      sql.NullTime                    `db:"product_delivery_at"`
	DeletedAt              sql.NullTime                    `db:"deleted_at"`
	EstimatedStartAt       sql.NullTime                    `db:"estimated_start_at"`
	EstimatedCompleteAt    sql.NullTime                    `db:"estimated_complete_at"`
	Sorting                int16                           `db:"sorting"`
}

func (rcv *ProductionOrderStage) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		ProductionOrderStageFieldID,
		ProductionOrderStageFieldProductionOrderID,
		ProductionOrderStageFieldStageID,
		ProductionOrderStageFieldStartedAt,
		ProductionOrderStageFieldCompletedAt,
		ProductionOrderStageFieldStatus,
		ProductionOrderStageFieldCondition,
		ProductionOrderStageFieldNote,
		ProductionOrderStageFieldData,
		ProductionOrderStageFieldCreatedAt,
		ProductionOrderStageFieldUpdatedAt,
		ProductionOrderStageFieldWaitingAt,
		ProductionOrderStageFieldReceptionAt,
		ProductionOrderStageFieldProductionStartAt,
		ProductionOrderStageFieldProductionCompletionAt,
		ProductionOrderStageFieldProductDeliveryAt,
		ProductionOrderStageFieldDeletedAt,
		ProductionOrderStageFieldEstimatedStartAt,
		ProductionOrderStageFieldEstimatedCompleteAt,
		ProductionOrderStageFieldSorting,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.ProductionOrderID,
		&rcv.StageID,
		&rcv.StartedAt,
		&rcv.CompletedAt,
		&rcv.Status,
		&rcv.Condition,
		&rcv.Note,
		&rcv.Data,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
		&rcv.WaitingAt,
		&rcv.ReceptionAt,
		&rcv.ProductionStartAt,
		&rcv.ProductionCompletionAt,
		&rcv.ProductDeliveryAt,
		&rcv.DeletedAt,
		&rcv.EstimatedStartAt,
		&rcv.EstimatedCompleteAt,
		&rcv.Sorting,
	}

	return
}

func (*ProductionOrderStage) TableName() string {
	return "production_order_stages"
}
