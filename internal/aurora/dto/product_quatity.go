package dto

import (
	"mmlabel.gitlab.com/mm-printing-backend/pkg/commondto"
	"time"
)

type ProductQualityFilter struct {
	ProductionOrderID string    `json:"productionOrderID"`
	DefectType        string    `json:"defectType"`
	DefectCode        string    `json:"defectCode"`
	CreatedAtFrom     time.Time `json:"createdAtFrom"`
	CreatedAtTo       time.Time `json:"createdAtTo"`
}

type FindProductQualityRequest struct {
	Filter *ProductQualityFilter `json:"filter" binding:"required"`
	Paging *commondto.Paging     `json:"paging" binding:"required"`
}

type FindProductQualityResponse struct {
	ProductQuality []*ProductQuality         `json:"productQuality"`
	Total          int64                     `json:"total"`
	Analysis       []*ProductQualityAnalysis `json:"analysis"`
}
type ProductQualityAnalysis struct {
	DefectType string `json:"defectType"`
	Count      int64  `json:"count"`
}
type ProductQuality struct {
	ID                  string    `json:"id"`
	ProductionOrderID   string    `json:"productionOrderID"`
	ProductionOrderName string    `json:"productionOrderName"`
	ProductID           string    `json:"productID"`
	DefectType          string    `json:"defectType"`
	DefectCode          string    `json:"defectCode"`
	DefectLevel         int16     `json:"defectLevel"`
	ProductionStageID   string    `json:"productionStageID"`
	DefectiveQuantity   int64     `json:"defectiveQuantity"`
	GoodQuantity        int64     `json:"goodQuantity"`
	Description         string    `json:"description"`
	CreatedBy           string    `json:"createdBy"`
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`
}

type CreateProductQualityRequest struct {
	ProductionOrderID string `json:"productionOrderID"  binding:"required"`
	ProductID         string `json:"productID"`
	DefectType        string `json:"defectType" binding:"required"`
	DefectCode        string `json:"defectCode" binding:"required"`
	DefectLevel       int16  `json:"defectLevel"`
	ProductionStageID string `json:"productionStageID"`
	DefectiveQuantity int64  `json:"defectiveQuantity"`
	GoodQuantity      int64  `json:"goodQuantity"`
	Description       string `json:"description"`
}

type CreateProductQualityResponse struct {
	ID string `json:"id"`
}

type EditProductQualityRequest struct {
	ID                string `json:"id" binding:"required"`
	DefectType        string `json:"defectType"`
	DefectCode        string `json:"defectCode"`
	DefectLevel       string `json:"defectLevel"`
	ProductionStageID string `json:"productionStageID"`
	DefectiveQuantity int64  `json:"defectiveQuantity"`
	GoodQuantity      int64  `json:"goodQuantity"`
	Description       string `json:"description"`
}

type EditProductQualityResponse struct {
}

type DeleteProductQualityRequest struct {
	ID string `json:"id"`
}

type DeleteProductQualityResponse struct {
}
