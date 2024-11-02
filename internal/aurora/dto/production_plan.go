package dto

import (
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/commondto"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

type ProductionPlanFilter struct {
	IDs  []string `json:"ids"`
	Name string   `json:"name"`
	//CustomerID  string                      `json:"customerID"`
	ProductName string                      `json:"productName"`
	ProductCode string                      `json:"productCode"`
	Statuses    []enum.ProductionPlanStatus `json:"statuses"`
	Stage       enum.ProductionPlanStage    `json:"stage"`
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
	Name         string         `json:"name" binding:"required"`
	ProductName  string         `json:"productName" binding:"required"`
	ProductCode  string         `json:"productCode" binding:"required"`
	QtyPaper     int64          `json:"qtyPaper,omitempty"`
	QtyFinished  int64          `json:"qtyFinished,omitempty"`
	QtyDelivered int64          `json:"qtyDelivered,omitempty"`
	Thumbnail    string         `json:"thumbnail,omitempty"`
	Note         string         `json:"note,omitempty"`
	Workflow     any            `json:"workflow"`
	CustomField  []*CustomField `json:"customField" binding:"required"`
}

type CreateProductionPlanResponse struct {
	ID string `json:"id"`
}

type EditProductionPlanRequest struct {
	ID           string                    `json:"id,omitempty"`
	Name         string                    `json:"name" binding:"required"`
	ProductName  string                    `json:"productName" binding:"required"`
	ProductCode  string                    `json:"productCode" binding:"required"`
	QtyPaper     int64                     `json:"qtyPaper,omitempty"`
	QtyFinished  int64                     `json:"qtyFinished,omitempty"`
	QtyDelivered int64                     `json:"qtyDelivered,omitempty"`
	Thumbnail    string                    `json:"thumbnail,omitempty"`
	Status       enum.ProductionPlanStatus `json:"status,omitempty"`
	Note         string                    `json:"note,omitempty"`
	CustomField  []*CustomField            `json:"customField,omitempty"`
	Workflow     any                       `json:"workflow"`
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
	ID                string                    `json:"id,omitempty"`
	ProductionOrderID string                    `json:"productionOrderID,omitempty"`
	CustomerData      *Customer                 `json:"customerData,omitempty"`
	ProductName       string                    `json:"productName,omitempty"`
	ProductCode       string                    `json:"productCode,omitempty"`
	QtyPaper          int64                     `json:"qtyPaper,omitempty"`
	QtyFinished       int64                     `json:"qtyFinished,omitempty"`
	QtyDelivered      int64                     `json:"qtyDelivered,omitempty"`
	Thumbnail         string                    `json:"thumbnail,omitempty"`
	Status            enum.ProductionPlanStatus `json:"status"`
	CurrentStage      enum.ProductionPlanStage  `json:"currentStage"`
	Note              string                    `json:"note,omitempty"`
	CreatedBy         string                    `json:"createdBy"`
	CreatedAt         time.Time                 `json:"createdAt"`
	UpdatedBy         string                    `json:"updatedBy"`
	UpdatedAt         time.Time                 `json:"updatedAt"`
	CreatedByName     string                    `json:"createdByName"`
	UpdatedByName     string                    `json:"updatedByName,omitempty"`
	Name              string                    `json:"name,omitempty"`
	CustomData        map[string]string         `json:"customData,omitempty"`
	UserFields        map[string][]*UserField   `json:"userFields,omitempty"`
	Workflow          any                       `json:"workflow"`
}

type UpdateCustomFieldPLValuesRequest struct {
	ProductionPlanID string         `json:"productionPlanID" binding:"required"`
	CustomField      []*CustomField `json:"customField" binding:"required"`
}

type UpdateCustomFieldPLValuesResponse struct{}

type UpdateCurrentStageRequest struct {
	ProductionPlanID string                   `json:"productionPlanID" binding:"required"`
	CurrentStage     enum.ProductionPlanStage `json:"currentStage" binding:"required"`
}

type UpdateCurrentStageResponse struct{}

type SummaryProductionPlanRequest struct {
	StartDate time.Time `json:"startDate,omitempty"`
	EndDate   time.Time `json:"endDate,omitempty"`
}

type SummaryProductionPlanResponse struct {
	Items []*SummaryProductionPlanItem `json:"items,omitempty"`
	Total int64                        `json:"total,omitempty"`
}

type SummaryProductionPlanItem struct {
	Stage  enum.ProductionPlanStage  `json:"stage,omitempty"`
	Status enum.ProductionPlanStatus `json:"status,omitempty"`
	Count  int64                     `json:"count,omitempty"`
}
