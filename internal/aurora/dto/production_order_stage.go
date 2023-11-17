package dto

import (
	"mmlabel.gitlab.com/mm-printing-backend/pkg/commondto"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"time"
)

type ProductionOrderStageFilter struct {
	ProductOrderID string                          `json:"productOrderID"`
	Status         enum.ProductionOrderStageStatus `json:"status"`
}

type FindProductionOrderStagesRequest struct {
	Filter *ProductionOrderStageFilter `json:"filter" binding:"required"`
	Paging *commondto.Paging           `json:"paging" binding:"required"`
}

type FindProductionOrderStagesResponse struct {
	ProductionOrderStages []*ProductionOrderStage `json:"production_order_stages"`
	Total                 int64                   `json:"total"`
}
type ProductionOrderStage struct {
	ID                  string                          `json:"ID"`
	ProductionOrderID   string                          `json:"ProductionOrderID"`
	StageID             string                          `json:"StageID"`
	EstimatedStartAt    time.Time                       `json:"EstimatedStartAt"`
	EstimatedCompleteAt time.Time                       `json:"EstimatedCompleteAt"`
	StartedAt           time.Time                       `json:"StartedAt"`
	CompletedAt         time.Time                       `json:"CompletedAt"`
	Status              enum.ProductionOrderStageStatus `json:"Status"`
	Condition           string                          `json:"Condition"`
	Note                string                          `json:"Note"`
	Data                map[string]interface{}          `json:"Data"`
	CreatedAt           time.Time                       `json:"CreatedAt"`
	UpdatedAt           time.Time                       `json:"UpdatedAt"`
}

type CreateProductionOrderStageRequest struct {
	ProductionOrderID   string                          `json:"ProductionOrderID"  binding:"required"`
	StageID             string                          `json:"StageID"  binding:"required"`
	EstimatedStartAt    time.Time                       `json:"EstimatedStartAt"`
	EstimatedCompleteAt time.Time                       `json:"EstimatedCompleteAt"`
	StartedAt           time.Time                       `json:"StartedAt"`
	CompletedAt         time.Time                       `json:"CompletedAt"`
	Status              enum.ProductionOrderStageStatus `json:"Status"  binding:"required"`
	Condition           string                          `json:"Condition"`
	Note                string                          `json:"Note"`
	Data                map[string]interface{}          `json:"Data"`
	CreatedAt           time.Time                       `json:"CreatedAt"`
	UpdatedAt           time.Time                       `json:"UpdatedAt"`
}

type CreateProductionOrderStageResponse struct {
	ID string `json:"id"`
}

type EditProductionOrderStageRequest struct {
	ID                  string                          `json:"id" binding:"required"`
	EstimatedStartAt    time.Time                       `json:"EstimatedStartAt"`
	EstimatedCompleteAt time.Time                       `json:"EstimatedCompleteAt"`
	StartedAt           time.Time                       `json:"StartedAt"`
	CompletedAt         time.Time                       `json:"CompletedAt"`
	Status              enum.ProductionOrderStageStatus `json:"Status"  binding:"required"`
	Condition           string                          `json:"Condition"`
	Note                string                          `json:"Note"`
	Data                map[string]interface{}          `json:"Data"`
	Sorting             int16                           `json:"sorting"`
	CreatedAt           time.Time                       `json:"CreatedAt"`
	UpdatedAt           time.Time                       `json:"UpdatedAt"`
}

type EditProductionOrderStageResponse struct {
}

type DeleteProductionOrderStageRequest struct {
	ID string `json:"id"`
}

type DeleteProductionOrderStageResponse struct {
}
