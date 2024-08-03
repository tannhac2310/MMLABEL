package dto

import (
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/commondto"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

type ProductionPlanFilter struct {
	IDs        []string                     `json:"ids"`
	Name       string                       `json:"name"`
	CustomerID string                       `json:"customerID"`
	Statuses   []enum.ProductionOrderStatus `json:"statuses"`
}

type FindProductionPlansRequest struct {
	Filter *ProductionPlanFilter `json:"filter" binding:"required"`
	Paging *commondto.Paging     `json:"paging" binding:"required"`
	Sort   *commondto.Sort       `json:"sort"`
}

type FindProductionPlansResponse struct {
	ProductionPlans []*ProductionPlan
	Total           int64
}

type CreateProductionPlanRequest struct {
	Name        string                    `json:"name,omitempty"`
	CustomerID  string                    `json:"customerID,omitempty"`
	SalesID     string                    `json:"salesID,omitempty"`
	Thumbnail   string                    `json:"thumbnail,omitempty"`
	Status      enum.ProductionPlanStatus `json:"status,omitempty"`
	Note        string                    `json:"note,omitempty"`
	CustomField []*CustomField            `json:"customField,omitempty"`
	CreatedBy   string                    `json:"createdBy,omitempty"`
}

type CreateProductionPlanResponse struct {
	ID string `json:"id"`
}

type EditProductionPlanRequest struct{}

type EditProductionPlanResponse struct{}

type DeleteProductionPlanRequest struct {
	ID string `json:"id"`
}

type DeleteProductionPlanResponse struct{}

type ProductionPlan struct {
	ID         string                    `json:"id,omitempty"`
	CustomerID string                    `json:"customerID,omitempty"`
	SalesID    string                    `json:"salesID,omitempty"`
	Thumbnail  string                    `json:"thumbnail,omitempty"`
	Status     enum.ProductionPlanStatus `json:"status,omitempty"`
	Note       string                    `json:"note,omitempty"`
	CreatedBy  string                    `json:"createdBy,omitempty"`
	CreatedAt  time.Time                 `json:"createdAt,omitempty"`
	UpdatedBy  string                    `json:"updatedBy,omitempty"`
	UpdatedAt  time.Time                 `json:"updatedAt,omitempty"`
	DeletedAt  time.Time                 `json:"deletedAt,omitempty"`
	Name       string                    `json:"name,omitempty"`
	CustomData map[string]string         `json:"customData,omitempty"`
}
