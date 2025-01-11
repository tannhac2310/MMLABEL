package dto

import (
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/commondto"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

type ProductionOrderFilter struct {
	IDs                             []string                        `json:"ids"`
	ProductionPlanIDs               []string                        `json:"productionPlanIDs"`
	Name                            string                          `json:"name"`
	CustomerID                      string                          `json:"customerID"`
	Status                          enum.ProductionOrderStatus      `json:"status"`
	Statuses                        []enum.ProductionOrderStatus    `json:"statuses"`
	OrderStageStatus                enum.ProductionOrderStageStatus `json:"orderStageStatus"`
	EstimatedStartAtTo              time.Time                       `json:"estimatedStartAtTo"`
	EstimatedStartAtFrom            time.Time                       `json:"estimatedStartAtFrom"`
	EstimatedCompleteAtFrom         time.Time                       `json:"estimatedCompleteAtFrom"`
	EstimatedCompleteAtTo           time.Time                       `json:"estimatedCompleteAtTo"`
	OrderStageEstimatedStartFrom    time.Time                       `json:"orderStageEstimatedStartFrom"`
	OrderStageEstimatedStartTo      time.Time                       `json:"orderStageEstimatedStartTo"`
	OrderStageEstimatedCompleteFrom time.Time                       `json:"orderStageEstimatedCompleteFrom"`
	OrderStageEstimatedCompleteTo   time.Time                       `json:"orderStageEstimatedCompleteTo"`
	Responsible                     []string                        `json:"responsible"`
	StageIDs                        []string                        `json:"stageIDs"`
	StageInLine                     string                          `json:"stageInLine"` // search lsx mà theo công đoạn StageInLine đang sản xuất: production_start
	DeviceID                        string                          `json:"deviceID"`
}

type FindProductionOrdersRequest struct {
	Filter *ProductionOrderFilter `json:"filter" binding:"required"`
	Paging *commondto.Paging      `json:"paging" binding:"required"`
	Sort   *commondto.Sort        `json:"sort"`
}

type FindProductionOrdersResponse struct {
	ProductionOrders []*ProductionOrder `json:"productionOrders"`
	Analysis         []*Analysis        `json:"analysis"`
	Total            int64              `json:"total"`
}
type Analysis struct {
	Status enum.ProductionOrderStatus `json:"status"`
	Count  int64                      `json:"count"`
}

type FindProductionOrdersWithNoPermissionResponse struct {
	ProductionOrders []*ProductionOrder2 `json:"production_orders"`
	Total            int64               `json:"total"`
}

type ProductionOrder2 struct {
	ID                  string                     `json:"id"`
	Name                string                     `json:"name"`
	ProductCode         string                     `json:"productCode"`
	ProductName         string                     `json:"productName"`
	CustomerID          string                     `json:"customerID"`
	SalesID             string                     `json:"salesID"`
	QtyPaper            int64                      `json:"qtyPaper"`
	QtyFinished         int64                      `json:"qtyFinished"`
	QtyDelivered        int64                      `json:"qtyDelivered"`
	EstimatedStartAt    time.Time                  `json:"estimatedStartAt"`
	EstimatedCompleteAt time.Time                  `json:"estimatedCompleteAt"`
	DeliveryDate        time.Time                  `json:"deliveryDate"`
	DeliveryImage       string                     `json:"deliveredImage"`
	Status              enum.ProductionOrderStatus `json:"status"`
	Note                string                     `json:"note"`
	Version             int64                      `json:"version"`
}
type ProductionOrder struct {
	ID                    string                     `json:"id"`
	Name                  string                     `json:"name"`
	ProductCode           string                     `json:"productCode"`
	ProductName           string                     `json:"productName"`
	CustomerID            string                     `json:"customerID"`
	ProductionPlanID      string                     `json:"productionPlanID"`
	SalesID               string                     `json:"salesID"`
	QtyPaper              int64                      `json:"qtyPaper"`
	QtyFinished           int64                      `json:"qtyFinished"`
	QtyDelivered          int64                      `json:"qtyDelivered"`
	EstimatedStartAt      time.Time                  `json:"estimatedStartAt"`
	EstimatedCompleteAt   time.Time                  `json:"estimatedCompleteAt"`
	DeliveryDate          time.Time                  `json:"deliveryDate"`
	DeliveryImage         string                     `json:"deliveredImage"`
	Status                enum.ProductionOrderStatus `json:"status"`
	Note                  string                     `json:"note"`
	ProductionOrderStages []*OrderStage              `json:"production_order_stages"`
	CustomData            map[string]string          `json:"customData"`
	CustomerData          *Customer                  `json:"customerData"`
	CreatedBy             string                     `json:"createdBy"`
	CreatedAt             time.Time                  `json:"createdAt"`
	UpdatedAt             time.Time                  `json:"updatedAt"`
	Version               int64                      `json:"version"`
	Data                  any                        `json:"data"`
}

type CreateProductionOrderRequest struct {
	Name                  string                     `json:"name" binding:"required"`
	ProductCode           string                     `json:"productCode"  binding:"required"`
	ProductName           string                     `json:"productName"  binding:"required"`
	CustomerID            string                     `json:"customerID"  binding:"required"`
	SalesID               string                     `json:"salesID"`
	QtyPaper              int64                      `json:"qtyPaper"`
	QtyFinished           int64                      `json:"qtyFinished"`
	QtyDelivered          int64                      `json:"qtyDelivered"`
	EstimatedStartAt      time.Time                  `json:"estimatedStartAt"`
	EstimatedCompleteAt   time.Time                  `json:"estimatedCompleteAt"`
	DeliveryDate          time.Time                  `json:"deliveryDate"`
	DeliveryImage         string                     `json:"deliveredImage"`
	ProductionOrderStages []CreateOrderStage         `json:"production_order_stages"`
	Status                enum.ProductionOrderStatus `json:"status"  binding:"required"`
	Note                  string                     `json:"note"`
	CustomField           []*CustomField             `json:"customField"`
}
type CustomField struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type OrderStageDevice struct {
	ID                     string                                `json:"id"`
	ProductionOrderStageID string                                `json:"productionOrderStageID"`
	DeviceID               string                                `json:"deviceID"`
	DeviceName             string                                `json:"deviceName"`
	DeviceData             any                                   `json:"deviceData"`
	ResponsibleObject      []*User                               `json:"responsibleObject"`
	Quantity               int64                                 `json:"quantity"`
	ProcessStatus          enum.ProductionOrderStageDeviceStatus `json:"processStatus"`
	Status                 enum.CommonStatus                     `json:"status"`
	Responsible            []string                              `json:"responsible"`
	Settings               map[string]interface{}                `json:"settings"`
	Note                   string                                `json:"note"`
	EstimatedCompleteAt    time.Time                             `json:"estimatedCompleteAt"`
	AssignedQuantity       int64                                 `json:"assignedQuantity"`
}
type User struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Avatar  string `json:"avatar"`
	Address string `json:"address"`
}
type OrderStage struct {
	ID                     string                          `json:"id"`
	ProductionOrderID      string                          `json:"productionOrderID"`
	StageID                string                          `json:"stageID"`
	StageName              string                          `json:"stageName"`
	StartedAt              time.Time                       `json:"startedAt"`
	CompletedAt            time.Time                       `json:"completedAt"`
	Status                 enum.ProductionOrderStageStatus `json:"status"`
	Condition              string                          `json:"condition"`
	Note                   string                          `json:"note"`
	Data                   map[string]interface{}          `json:"data"`
	SoLuong                int64                           `json:"soLuong"`
	CreatedAt              time.Time                       `json:"createdAt"`
	UpdatedAt              time.Time                       `json:"updatedAt"`
	WaitingAt              time.Time                       `json:"waitingAt"`
	ReceptionAt            time.Time                       `json:"receptionAt"`
	ProductionStartAt      time.Time                       `json:"productionStartAt"`
	ProductionCompletionAt time.Time                       `json:"productionCompletionAt"`
	ProductDeliveryAt      time.Time                       `json:"productDeliveryAt"`
	EstimatedStartAt       time.Time                       `json:"estimatedStartAt"`
	EstimatedCompleteAt    time.Time                       `json:"estimatedCompleteAt"`
	Sorting                int16                           `json:"sorting"`
	OrderStageDevices      []*OrderStageDevice             `json:"orderStageDevices"`
}

type CreateOrderStage struct {
	StageID             string                          `json:"stageID"`
	EstimatedStartAt    time.Time                       `json:"estimatedStartAt"`
	EstimatedCompleteAt time.Time                       `json:"estimatedCompleteAt"`
	StartedAt           time.Time                       `json:"startedAt"`
	CompletedAt         time.Time                       `json:"completedAt"`
	Status              enum.ProductionOrderStageStatus `json:"status"`
	Condition           string                          `json:"condition"`
	SoLuong             int64                           `json:"soLuong"`
	Note                string                          `json:"note"`
	Data                map[string]interface{}          `json:"data"`
}
type CreateProductionOrderResponse struct {
	ID string `json:"id"  binding:"required"`
}
type AcceptAndChangeNextStageRequest struct {
	ID string `json:"id"  binding:"required"`
}
type AcceptAndChangeNextStageResponse struct {
}

type EditProductionOrderRequest struct {
	ID                  string                     `json:"id" binding:"required"`
	Name                string                     `json:"name" binding:"required"`
	ProductCode         string                     `json:"productCode"  binding:"required"`
	ProductName         string                     `json:"productName"  binding:"required"`
	QtyPaper            int64                      `json:"qtyPaper"`
	QtyFinished         int64                      `json:"qtyFinished"`
	QtyDelivered        int64                      `json:"qtyDelivered"`
	EstimatedStartAt    time.Time                  `json:"estimatedStartAt"`
	EstimatedCompleteAt time.Time                  `json:"estimatedCompleteAt"`
	DeliveryDate        time.Time                  `json:"deliveryDate"`
	DeliveryImage       string                     `json:"deliveredImage"`
	Status              enum.ProductionOrderStatus `json:"status"  binding:"required"`
	SoLuong             int64                      `json:"soLuong"`
	Note                string                     `json:"note"`
	//ProductionOrderStages []EditOrderStage           `json:"production_order_stages" binding:"required"`
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
