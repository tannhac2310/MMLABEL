package dto

import (
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/commondto"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

type ProductionPlanFilter struct {
	IDs         []string                    `json:"ids"`
	Name        string                      `json:"name"`
	CustomerID  string                      `json:"customerID"`
	ProductName string                      `json:"productName"`
	ProductCode string                      `json:"productCode"`
	Statuses    []enum.ProductionPlanStatus `json:"statuses"`
	Stage       int                         `json:"stage"`
}

type FindProductionPlansRequest struct {
	Filter *ProductionPlanFilter `json:"filter" binding:"required"`
	Paging *commondto.Paging     `json:"paging" binding:"required"`
	Sort   *commondto.Sort       `json:"sort"`
}

type FindProductionPlansResponse struct {
	ProductionPlans []*ProductionPlan `json:"productionPlans"`
	Total           int64             `json:"total"`
}

type CreateProductionPlanRequest struct {
	Name         string                    `json:"name" binding:"required"`
	CustomerID   string                    `json:"customerID" binding:"required"`
	SalesID      string                    `json:"salesID" binding:"required"`
	ProductName  string                    `json:"productName" binding:"required"`
	ProductCode  string                    `json:"productCode" binding:"required"`
	QtyPaper     int64                     `json:"qtyPaper,omitempty"`
	QtyFinished  int64                     `json:"qtyFinished,omitempty"`
	QtyDelivered int64                     `json:"qtyDelivered,omitempty"`
	Thumbnail    string                    `json:"thumbnail,omitempty"`
	Status       enum.ProductionPlanStatus `json:"status" binding:"required"`
	Note         string                    `json:"note,omitempty"`
	CustomField  []*CustomField            `json:"customField" binding:"required"`
}

type CreateProductionPlanResponse struct {
	ID string `json:"id"`
}

type EditProductionPlanRequest struct {
	ID           string                    `json:"id,omitempty"`
	Name         string                    `json:"name" binding:"required"`
	CustomerID   string                    `json:"customerID" binding:"required"`
	SalesID      string                    `json:"salesID" binding:"required"`
	ProductName  string                    `json:"productName" binding:"required"`
	ProductCode  string                    `json:"productCode" binding:"required"`
	QtyPaper     int64                     `json:"qtyPaper,omitempty"`
	QtyFinished  int64                     `json:"qtyFinished,omitempty"`
	QtyDelivered int64                     `json:"qtyDelivered,omitempty"`
	Thumbnail    string                    `json:"thumbnail,omitempty"`
	Status       enum.ProductionPlanStatus `json:"status,omitempty"`
	Note         string                    `json:"note,omitempty"`
	CustomField  []*CustomField            `json:"customField,omitempty"`
	CreatedBy    string                    `json:"createdBy,omitempty"`
}

type EditProductionPlanResponse struct{}

type DeleteProductionPlanRequest struct {
	ID string `json:"id"`
}

type DeleteProductionPlanResponse struct{}

type ProcessProductionOrderRequest struct {
	ID                  string             `json:"id"`
	Stages              []CreateOrderStage `json:"productionOrderStages"`
	EstimatedStartAt    time.Time          `json:"estimatedStartAt"`
	EstimatedCompleteAt time.Time          `json:"estimatedCompleteAt"`
}

type ProcessProductionOrderResponse struct {
	ID string `json:"id"` // production order id
}

type ProductionPlan struct {
	ID           string                    `json:"id,omitempty"`
	CustomerID   string                    `json:"customerID,omitempty"`
	CustomerData *Customer                 `json:"customerData,omitempty"`
	SalesID      string                    `json:"salesID,omitempty"`
	ProductName  string                    `json:"productName,omitempty"`
	ProductCode  string                    `json:"productCode,omitempty"`
	QtyPaper     int64                     `json:"qtyPaper,omitempty"`
	QtyFinished  int64                     `json:"qtyFinished,omitempty"`
	QtyDelivered int64                     `json:"qtyDelivered,omitempty"`
	Thumbnail    string                    `json:"thumbnail,omitempty"`
	Status       enum.ProductionPlanStatus `json:"status,omitempty"`
	Note         string                    `json:"note,omitempty"`
	CreatedBy    string                    `json:"createdBy,omitempty"`
	CreatedAt    time.Time                 `json:"createdAt,omitempty"`
	UpdatedBy    string                    `json:"updatedBy,omitempty"`
	UpdatedAt    time.Time                 `json:"updatedAt,omitempty"`
	Name         string                    `json:"name,omitempty"`
	CustomData   map[string]string         `json:"customData,omitempty"`
}
