package dto

import (
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/commondto"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
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
	ID                  string                          `json:"id"`
	ProductionOrderID   string                          `json:"productionOrderID"`
	StageID             string                          `json:"stageID"`
	EstimatedStartAt    time.Time                       `json:"estimatedStartAt"`
	EstimatedCompleteAt time.Time                       `json:"estimatedCompleteAt"`
	StartedAt           time.Time                       `json:"startedAt"`
	CompletedAt         time.Time                       `json:"completedAt"`
	Status              enum.ProductionOrderStageStatus `json:"status"`
	SoLuong             int64                           `json:"soLuong"`
	Condition           string                          `json:"condition"`
	Note                string                          `json:"note"`
	Data                map[string]interface{}          `json:"data"`
	CreatedAt           time.Time                       `json:"createdAt"`
	UpdatedAt           time.Time                       `json:"updatedAt"`
}

type CreateProductionOrderStageRequest struct {
	ProductionOrderID   string                          `json:"productionOrderID"  binding:"required"`
	StageID             string                          `json:"stageID" binding:"required"`
	EstimatedStartAt    time.Time                       `json:"estimatedStartAt"`
	EstimatedCompleteAt time.Time                       `json:"estimatedCompleteAt"`
	StartedAt           time.Time                       `json:"startedAt"`
	CompletedAt         time.Time                       `json:"completedAt"`
	Status              enum.ProductionOrderStageStatus `json:"status"`
	Condition           string                          `json:"condition"`
	Note                string                          `json:"note"`
	Data                map[string]interface{}          `json:"data"`
}

type CreateProductionOrderStageResponse struct {
	ID string `json:"id"`
}

type EditProductionOrderStageRequest struct {
	ID                  string                          `json:"id" binding:"required"`
	StageID             string                          `json:"stageID" binding:"required"`
	EstimatedStartAt    time.Time                       `json:"estimatedStartAt"`
	EstimatedCompleteAt time.Time                       `json:"estimatedCompleteAt"`
	StartedAt           time.Time                       `json:"startedAt"`
	CompletedAt         time.Time                       `json:"completedAt"`
	Status              enum.ProductionOrderStageStatus `json:"status" binding:"required"`
	SoLuong             int64                           `json:"soLuong"`
	Condition           string                          `json:"condition"`
	Note                string                          `json:"note"`
	Data                map[string]interface{}          `json:"data"`
	Sorting             int16                           `json:"sorting"`
}

type EditProductionOrderStageResponse struct {
}

type DeleteProductionOrderStageRequest struct {
	ID string `json:"id"`
}

type DeleteProductionOrderStageResponse struct {
}
