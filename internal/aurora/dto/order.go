package dto

import (
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/commondto"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

type CreateOrderRequest struct {
	Order OrderData        `json:"order"`
	Items []*OrderItemData `json:"items"`
}

type UpdateOrderRequest struct {
	Order OrderData        `json:"order"`
	Items []*OrderItemData `json:"items"`
}

type OrderData struct {
	ID                  string           `json:"id"`
	Title               string           `json:"title"`
	Code                string           `json:"code"`
	SaleName            string           `json:"saleName"`
	SaleAdminName       string           `json:"saleAdminName"`
	ProductCode         string           `json:"productCode"`
	ProductName         string           `json:"productName"`
	CustomerID          string           `json:"customerID"`
	CustomerProductCode string           `json:"customerProductCode"`
	CustomerProductName string           `json:"customerProductName"`
	Status              enum.OrderStatus `json:"status"`
}

type OrderItemData struct {
	ID                      string            `json:"id"`
	ProductionPlanProductID string            `json:"productionPlanProductID"`
	ProductionPlanID        string            `json:"productionPlanID"`
	ProductionQuantity      int64             `json:"productionQuantity"`
	Quantity                int64             `json:"quantity"`
	UnitPrice               float64           `json:"unitPrice"`
	DeliveredQuantity       int64             `json:"deliveredQuantity"`
	EstimatedDeliveryDate   time.Time         `json:"estimatedDeliveryDate"`
	DeliveredDate           time.Time         `json:"deliveredDate"`
	Status                  string            `json:"status"`
	Attachment              map[string]string `json:"attachment"`
	Note                    string            `json:"note"`
}

type CreateOrderResponse struct {
	ID string `json:"id"`
}

type UpdateOrderResponse struct {
	ID string `json:"id"`
}

type DeleteOrderRequest struct {
	ID string `json:"id"`
}

type DeleteOrderResponse struct {
	ID string `json:"id"`
}

type SearchOrderFilter struct {
	IDs      []string `json:"ids"`
	Title    string   `json:"title"`
	Code     string   `json:"code"`
	SaleName string   `json:"saleName"`
}

type SearchOrderRequest struct {
	Filter SearchOrderFilter `json:"filter" binding:"required"`
	Paging *commondto.Paging `json:"paging" binding:"required"`
}

type SearchOrderResponse struct {
	Orders []*OrderWithItems `json:"orders"`
	Total  int64             `json:"total"`
}

type OrderWithItems struct {
	Order OrderData
	Items []*OrderItemData
}
