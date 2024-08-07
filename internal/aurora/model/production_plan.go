package model

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"math"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

const (
	ProductionPlanFieldID           = "id"
	ProductionPlanFieldCustomerID   = "customer_id"
	ProductionPlanFieldSalesID      = "sales_id"
	ProductionPlanFieldThumbnail    = "thumbnail"
	ProductionPlanFieldStatus       = "status"
	ProductionPlanFieldNote         = "note"
	ProductionPlanFieldCreatedBy    = "created_by"
	ProductionPlanFieldCreatedAt    = "created_at"
	ProductionPlanFieldUpdatedBy    = "updated_by"
	ProductionPlanFieldUpdatedAt    = "updated_at"
	ProductionPlanFieldDeletedAt    = "deleted_at"
	ProductionPlanFieldName         = "name"
	ProductionPlanFieldPoStages     = "po_stages"
	ProductionPlanFieldCurrentStage = "current_stage"
)

type ProductionPlan struct {
	ID           string                    `db:"id"`
	CustomerID   string                    `db:"customer_id"`
	SalesID      string                    `db:"sales_id"`
	Thumbnail    sql.NullString            `db:"thumbnail"`
	Status       enum.ProductionPlanStatus `db:"status"`
	Note         sql.NullString            `db:"note"`
	CreatedBy    string                    `db:"created_by"`
	CreatedAt    time.Time                 `db:"created_at"`
	UpdatedBy    string                    `db:"updated_by"`
	UpdatedAt    time.Time                 `db:"updated_at"`
	DeletedAt    sql.NullTime              `db:"deleted_at"`
	Name         string                    `db:"name"`
	PoStagesInfo ProductionStageInfo       `db:"po_stages"`
	CurrentStage int                       `db:"current_stage"`
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
		ProductionPlanFieldPoStages,
		ProductionPlanFieldCurrentStage,
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
		&p.PoStagesInfo,
		&p.CurrentStage,
	}

	return
}

func (p *ProductionPlan) TableName() string {
	return "production_plans"
}

func (p *ProductionPlan) CanChangeStatusTo(s enum.ProductionPlanStatus) bool {
	if p.Status == s {
		return true
	}

	return math.Abs(float64(p.Status)-float64(s)) <= 1
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
