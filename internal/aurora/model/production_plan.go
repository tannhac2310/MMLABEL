package model

import (
	"database/sql"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

const (
	ProductionPlanFieldID         = "id"
	ProductionPlanFieldCustomerID = "customer_id"
	ProductionPlanFieldSalesID    = "sales_id"
	ProductionPlanFieldThumbnail  = "thumbnail"
	ProductionPlanFieldStatus     = "status"
	ProductionPlanFieldNote       = "note"
	ProductionPlanFieldCreatedBy  = "created_by"
	ProductionPlanFieldCreatedAt  = "created_at"
	ProductionPlanFieldUpdatedBy  = "updated_by"
	ProductionPlanFieldUpdatedAt  = "updated_at"
	ProductionPlanFieldDeletedAt  = "deleted_at"
	ProductionPlanFieldName       = "name"
)

const (
	ProductionPlanStageSale    = 1 // 0001
	ProductionPlanStageRAndD   = 2 // 0010
	ProductionPlanStageDesign  = 4 // 0100
	ProductionPlanStageFinance = 8 // 1000
)

type ProductionPlan struct {
	ID         string                    `db:"id"`
	CustomerID string                    `db:"customer_id"`
	SalesID    string                    `db:"sales_id"`
	Thumbnail  sql.NullString            `db:"thumbnail"`
	Status     enum.ProductionPlanStatus `db:"status"`
	Note       sql.NullString            `db:"note"`
	CreatedBy  string                    `db:"created_by"`
	CreatedAt  time.Time                 `db:"created_at"`
	UpdatedBy  string                    `db:"updated_by"`
	UpdatedAt  time.Time                 `db:"updated_at"`
	DeletedAt  sql.NullTime              `db:"deleted_at"`
	Name       string                    `db:"name"`
}

func (p *ProductionPlan) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		ProductionPlanFieldID,
		ProductionPlanFieldCustomerID,
		ProductionPlanFieldSalesID,
		ProductionPlanFieldThumbnail,
		ProductionPlanFieldStatus,
		ProductionPlanFieldNote,
		ProductionPlanFieldCreatedBy,
		ProductionPlanFieldCreatedAt,
		ProductionPlanFieldUpdatedBy,
		ProductionPlanFieldUpdatedAt,
		ProductionPlanFieldDeletedAt,
		ProductionPlanFieldName,
	}

	values = []interface{}{
		&p.ID,
		&p.CustomerID,
		&p.SalesID,
		&p.Thumbnail,
		&p.Status,
		&p.Note,
		&p.CreatedBy,
		&p.CreatedAt,
		&p.UpdatedBy,
		&p.UpdatedAt,
		&p.DeletedAt,
		&p.Name,
	}

	return
}

func (p *ProductionPlan) TableName() string {
	return "production_plans"
}

func (p *ProductionPlan) Editable() bool {
	switch p.Status {
	case enum.ProductionPlanStatusWaiting,
		enum.ProductionPlanStatusDoing,
		enum.ProductionPlanStatusPause,
		enum.ProductionPlanStatusComplete,
		enum.ProductionPlanStatusCancel:
		return true
	default:
		return false
	}
}

func (p *ProductionPlan) CanChangeStatusTo(s enum.ProductionPlanStatus) bool {
	if p.Status == s {
		return true
	}

	switch p.Status {
	case enum.ProductionPlanStatusWaiting:
		return s == enum.ProductionPlanStatusDoing || s == enum.ProductionPlanStatusPause || s == enum.ProductionPlanStatusCancel
	case enum.ProductionPlanStatusDoing:
		return s == enum.ProductionPlanStatusPause || s == enum.ProductionPlanStatusCancel || s == enum.ProductionPlanStatusComplete
	case enum.ProductionPlanStatusPause:
		return s == enum.ProductionPlanStatusDoing || s == enum.ProductionPlanStatusCancel
	case enum.ProductionPlanStatusComplete:
		return false
	case enum.ProductionPlanStatusCancel:
		return false
	default:
		return false
	}
}

type ProductionPlanStage struct {
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
