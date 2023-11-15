package dto

import (
	"mmlabel.gitlab.com/mm-printing-backend/pkg/commondto"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"time"
)

type ProductionOrderFilter struct {
	IDs                       []string
	ProductCode               string                     `json:"productCode"`
	ProductName               string                     `json:"productName"`
	CustomerID                string                     `json:"customerID"`
	Status                    enum.ProductionOrderStatus `json:"status"`
	PlannedProductionDateFrom time.Time                  `json:"plannedProductionDateFrom"`
	PlannedProductionDateTo   time.Time                  `json:"plannedProductionDateTo"`
	DeliveryDateFrom          time.Time                  `json:"deliveryDateFrom"`
	DeliveryDateTo            time.Time                  `json:"deliveryDateTo"`
}

type FindProductionOrdersRequest struct {
	Filter *ProductionOrderFilter `json:"filter" binding:"required"`
	Paging *commondto.Paging      `json:"paging" binding:"required"`
}

type FindProductionOrdersResponse struct {
	ProductionOrders []*ProductionOrder `json:"production_orders"`
	Total            int64              `json:"total"`
}
type ProductionOrder struct {
	ID                    string                     `json:"id"`
	ProductCode           string                     `json:"productCode"`
	ProductName           string                     `json:"productName"`
	CustomerID            string                     `json:"customerID"`
	SalesID               string                     `json:"salesID"`
	QtyPaper              int64                      `json:"qtyPaper"`
	QtyFinished           int64                      `json:"qtyFinished"`
	QtyDelivered          int64                      `json:"qtyDelivered"`
	PlannedProductionDate time.Time                  `json:"plannedProductionDate"`
	DeliveryDate          time.Time                  `json:"deliveryDate"`
	DeliveryImage         string                     `json:"deliveredImage"`
	Status                enum.ProductionOrderStatus `json:"status"`
	Note                  string                     `json:"note"`
	ProductionOrderStages []*OrderStage              `json:"production_order_stages"`
	CreatedBy             string                     `json:"createdBy"`
	CreatedAt             time.Time                  `json:"createdAt"`
	UpdatedAt             time.Time                  `json:"updatedAt"`
}

type CreateProductionOrderRequest struct {
	ProductCode           string                     `json:"productCode"  binding:"required"`
	ProductName           string                     `json:"productName"  binding:"required"`
	CustomerID            string                     `json:"customerID"  binding:"required"`
	SalesID               string                     `json:"salesID"`
	QtyPaper              int64                      `json:"qtyPaper"`
	QtyFinished           int64                      `json:"qtyFinished"`
	QtyDelivered          int64                      `json:"qtyDelivered"`
	PlannedProductionDate time.Time                  `json:"plannedProductionDate"`
	DeliveryDate          time.Time                  `json:"deliveryDate"`
	DeliveryImage         string                     `json:"deliveredImage"`
	ProductionOrderStages []CreateOrderStage         `json:"production_order_stages"`
	Status                enum.ProductionOrderStatus `json:"status"  binding:"required"`
	Note                  string                     `json:"note"`
}

type OrderStageDevice struct {
	ID                     string                                `json:"id"`
	ProductionOrderStageID string                                `json:"productionOrderStageID"`
	DeviceID               string                                `json:"deviceID"`
	Quantity               int64                                 `json:"quantity"`
	ProcessStatus          enum.ProductionOrderStageDeviceStatus `json:"processStatus"`
	Status                 enum.CommonStatus                     `json:"status"`
	Responsible            []string                              `json:"responsible"`
	Settings               map[string]interface{}                `json:"settings"`
	Note                   string                                `json:"note"`
}

type OrderStage struct {
	ID                  string                          `json:"id"`
	StageID             string                          `json:"stageID"`
	EstimatedStartAt    time.Time                       `json:"estimatedStartAt"`
	EstimatedCompleteAt time.Time                       `json:"estimatedCompleteAt"`
	StartedAt           time.Time                       `json:"startedAt"`
	CompletedAt         time.Time                       `json:"completedAt"`
	Status              enum.ProductionOrderStageStatus `json:"status"`
	Condition           string                          `json:"condition"`
	Note                string                          `json:"note"`
	Data                map[string]interface{}          `json:"data"`
	OrderStageDevices   []*OrderStageDevice             `json:"order_stage_devices"`
}

type CreateOrderStage struct {
	StageID             string                          `json:"stageID"`
	EstimatedStartAt    time.Time                       `json:"estimatedStartAt"`
	EstimatedCompleteAt time.Time                       `json:"estimatedCompleteAt"`
	StartedAt           time.Time                       `json:"startedAt"`
	CompletedAt         time.Time                       `json:"completedAt"`
	Status              enum.ProductionOrderStageStatus `json:"status"`
	Condition           string                          `json:"condition"`
	Note                string                          `json:"note"`
	Data                map[string]interface{}          `json:"data"`
}
type CreateProductionOrderResponse struct {
	ID string `json:"id"`
}

type EditProductionOrderRequest struct {
	ID                    string                     `json:"id" binding:"required"`
	ProductCode           string                     `json:"productCode"  binding:"required"`
	ProductName           string                     `json:"productName"  binding:"required"`
	QtyPaper              int64                      `json:"qtyPaper"`
	QtyFinished           int64                      `json:"qtyFinished"`
	QtyDelivered          int64                      `json:"qtyDelivered"`
	PlannedProductionDate time.Time                  `json:"plannedProductionDate"`
	DeliveryDate          time.Time                  `json:"deliveryDate"`
	DeliveryImage         string                     `json:"deliveryImage"`
	Status                enum.ProductionOrderStatus `json:"status"  binding:"required"`
	Note                  string                     `json:"note"`
	ProductionOrderStages []EditOrderStage           `json:"production_order_stages"`
}
type EditOrderStage struct {
	ID                  string                          `json:"id"`
	StageID             string                          `json:"stageID"`
	EstimatedStartAt    time.Time                       `json:"estimatedStartAt"`
	EstimatedCompleteAt time.Time                       `json:"estimatedCompleteAt"`
	StartedAt           time.Time                       `json:"startedAt"`
	CompletedAt         time.Time                       `json:"completedAt"`
	Status              enum.ProductionOrderStageStatus `json:"status"`
	Condition           string                          `json:"condition"`
	Note                string                          `json:"note"`
	Data                map[string]interface{}          `json:"data"`
}

type EditProductionOrderResponse struct {
}

type DeleteProductionOrderRequest struct {
	ID string `json:"id"`
}

type DeleteProductionOrderResponse struct {
}
