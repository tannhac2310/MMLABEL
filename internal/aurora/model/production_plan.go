package model

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

const (
	ProductionPlanFieldID                = "id"
	ProductionPlanFieldThumbnail         = "thumbnail"
	ProductionPlanFieldStatus            = "status"
	ProductionPlanFieldNote              = "note"
	ProductionPlanFieldCreatedBy         = "created_by"
	ProductionPlanFieldCreatedAt         = "created_at"
	ProductionPlanFieldUpdatedBy         = "updated_by"
	ProductionPlanFieldUpdatedAt         = "updated_at"
	ProductionPlanFieldDeletedAt         = "deleted_at"
	ProductionPlanFieldName              = "name"
	ProductionPlanFieldPoStages          = "po_stages"
	ProductionPlanFieldCurrentStage      = "current_stage"
	ProductionPlanFieldProductName       = "product_name"
	ProductionPlanFieldProductCode       = "product_code"
	ProductionPlanFieldQtyPaper          = "qty_paper"
	ProductionPlanFieldQtyFinished       = "qty_finished"
	ProductionPlanFieldQtyDelivered      = "qty_delivered"
	ProductionPlanFieldWorkflow          = "workflow"
	ProductionPlanFieldProductionOrderID = "production_order_id"
	ProductionPlanFieldSearchContent     = "search_content"
)

type ProductionPlan struct {
	ID                string                    `db:"id"`
	Thumbnail         sql.NullString            `db:"thumbnail"`
	Status            enum.ProductionPlanStatus `db:"status"`
	Note              sql.NullString            `db:"note"`
	CreatedBy         string                    `db:"created_by"`
	CreatedAt         time.Time                 `db:"created_at"`
	UpdatedBy         string                    `db:"updated_by"`
	UpdatedAt         time.Time                 `db:"updated_at"`
	DeletedAt         sql.NullTime              `db:"deleted_at"`
	Name              string                    `db:"name"`
	PoStages          ProductionStageInfo       `db:"po_stages"`
	CurrentStage      enum.ProductionPlanStage  `db:"current_stage"`
	ProductName       string                    `db:"product_name"`
	ProductCode       string                    `db:"product_code"`
	QtyPaper          int64                     `db:"qty_paper"`
	QtyFinished       int64                     `db:"qty_finished"`
	QtyDelivered      int64                     `db:"qty_delivered"`
	Workflow          any                       `db:"workflow"`
	ProductionOrderID sql.NullString            `db:"production_order_id"`
	SearchContent     string                    `db:"search_content"`
}

func (rcv *ProductionPlan) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		ProductionPlanFieldID,
		ProductionPlanFieldThumbnail,
		ProductionPlanFieldStatus,
		ProductionPlanFieldNote,
		ProductionPlanFieldCreatedBy,
		ProductionPlanFieldCreatedAt,
		ProductionPlanFieldUpdatedBy,
		ProductionPlanFieldUpdatedAt,
		ProductionPlanFieldDeletedAt,
		ProductionPlanFieldName,
		ProductionPlanFieldPoStages,
		ProductionPlanFieldCurrentStage,
		ProductionPlanFieldProductName,
		ProductionPlanFieldProductCode,
		ProductionPlanFieldQtyPaper,
		ProductionPlanFieldQtyFinished,
		ProductionPlanFieldQtyDelivered,
		ProductionPlanFieldWorkflow,
		ProductionPlanFieldProductionOrderID,
		ProductionPlanFieldSearchContent,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.Thumbnail,
		&rcv.Status,
		&rcv.Note,
		&rcv.CreatedBy,
		&rcv.CreatedAt,
		&rcv.UpdatedBy,
		&rcv.UpdatedAt,
		&rcv.DeletedAt,
		&rcv.Name,
		&rcv.PoStages,
		&rcv.CurrentStage,
		&rcv.ProductName,
		&rcv.ProductCode,
		&rcv.QtyPaper,
		&rcv.QtyFinished,
		&rcv.QtyDelivered,
		&rcv.Workflow,
		&rcv.ProductionOrderID,
		&rcv.SearchContent,
	}

	return
}

func (rcv *ProductionPlan) TableName() string {
	return "production_plans"
}

func (rcv *ProductionPlan) Editable() bool {
	return rcv.CurrentStage == enum.ProductionPlanStageSale || rcv.CurrentStage == enum.ProductionPlanStageDesign || rcv.CurrentStage == enum.ProductionPlanStageRandD
}

type ProductionStageInfo struct {
	Items []*ProductionStageItem `json:"items,omitempty"`
}

type ProductionStageItem struct {
	ID                  string    `json:"id,omitempty"`
	ProductionPlanID    string    `json:"productionPlanID,omitempty"`
	StageID             string    `json:"stageID,omitempty"`
	Note                string    `json:"note,omitempty"`
	CreatedAt           time.Time `json:"createdAt,omitempty"`
	UpdatedAt           time.Time `json:"updatedAt,omitempty"`
	EstimatedStartAt    time.Time `json:"estimatedStartAt,omitempty"`
	EstimatedCompleteAt time.Time `json:"estimatedCompleteAt,omitempty"`
	Sorting             int16     `json:"sorting,omitempty"`
}

// implement Scanner for the element type of the slice
func (s *ProductionStageInfo) Scan(src any) error {
	var data []byte
	switch v := src.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	}
	return json.Unmarshal(data, s)
}

// Value implements the [driver.Valuer] interface.
func (s ProductionStageInfo) Value() (driver.Value, error) {
	return json.Marshal(s)
}
