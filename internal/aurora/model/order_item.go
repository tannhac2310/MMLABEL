package model

import (
	"database/sql"
	"time"
)

const (
	OrderItemFieldID                      = "id"
	OrderItemFieldOrderID                 = "order_id"
	OrderItemFieldProductionPlanProductID = "production_plan_product_id"
	OrderItemFieldProductionPlanID        = "production_plan_id"
	OrderItemFieldProductionQuantity      = "production_quantity"
	OrderItemFieldQuantity                = "quantity"
	OrderItemFieldUnitPrice               = "unit_price"
	OrderItemFieldDeliveredQuantity       = "delivered_quantity"
	OrderItemFieldEstimatedDeliveryDate   = "estimated_delivery_date"
	OrderItemFieldDeliveredDate           = "delivered_date"
	OrderItemFieldStatus                  = "status"
	OrderItemFieldAttachment              = "attachment"
	OrderItemFieldNote                    = "note"
	OrderItemFieldCreatedBy               = "created_by"
	OrderItemFieldUpdatedBy               = "updated_by"
	OrderItemFieldCreatedAt               = "created_at"
	OrderItemFieldUpdatedAt               = "updated_at"
)

type OrderItem struct {
	ID                      string            `db:"id"`
	OrderID                 string            `db:"order_id"`
	ProductionPlanProductID string            `db:"production_plan_product_id"`
	ProductionPlanID        string            `db:"production_plan_id"`
	ProductionQuantity      int64             `db:"production_quantity"`
	Quantity                int64             `db:"quantity"`
	UnitPrice               float64           `db:"unit_price"`
	DeliveredQuantity       int64             `db:"delivered_quantity"`
	EstimatedDeliveryDate   sql.NullTime      `db:"estimated_delivery_date"`
	DeliveredDate           sql.NullTime      `db:"delivered_date"`
	Status                  string            `db:"status"`
	Attachment              map[string]string `db:"attachment"`
	Note                    string            `db:"note"`
	CreatedBy               string            `db:"created_by"`
	UpdatedBy               string            `db:"updated_by"`
	CreatedAt               time.Time         `db:"created_at"`
	UpdatedAt               time.Time         `db:"updated_at"`
}

func (rcv *OrderItem) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		OrderItemFieldID,
		OrderItemFieldOrderID,
		OrderItemFieldProductionPlanProductID,
		OrderItemFieldProductionPlanID,
		OrderItemFieldProductionQuantity,
		OrderItemFieldQuantity,
		OrderItemFieldUnitPrice,
		OrderItemFieldDeliveredQuantity,
		OrderItemFieldEstimatedDeliveryDate,
		OrderItemFieldDeliveredDate,
		OrderItemFieldStatus,
		OrderItemFieldAttachment,
		OrderItemFieldNote,
		OrderItemFieldCreatedBy,
		OrderItemFieldUpdatedBy,
		OrderItemFieldCreatedAt,
		OrderItemFieldUpdatedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.OrderID,
		&rcv.ProductionPlanProductID,
		&rcv.ProductionPlanID,
		&rcv.ProductionQuantity,
		&rcv.Quantity,
		&rcv.UnitPrice,
		&rcv.DeliveredQuantity,
		&rcv.EstimatedDeliveryDate,
		&rcv.DeliveredDate,
		&rcv.Status,
		&rcv.Attachment,
		&rcv.Note,
		&rcv.CreatedBy,
		&rcv.UpdatedBy,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
	}

	return
}

func (*OrderItem) TableName() string {
	return "order_items"
}
