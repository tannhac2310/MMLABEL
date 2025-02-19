package model

import (
	"database/sql"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

const (
	ProductionOrderFieldID                    = "id"
	ProductionOrderFieldProductCode           = "product_code"
	ProductionOrderFieldProductName           = "product_name"
	ProductionOrderFieldCustomerID            = "customer_id"
	ProductionOrderFieldSalesID               = "sales_id"
	ProductionOrderFieldQtyPaper              = "qty_paper"
	ProductionOrderFieldQtyFinished           = "qty_finished"
	ProductionOrderFieldQtyDelivered          = "qty_delivered"
	ProductionOrderFieldPlannedProductionDate = "planned_production_date"
	ProductionOrderFieldDeliveryDate          = "delivery_date"
	ProductionOrderFieldDeliveryImage         = "delivery_image"
	ProductionOrderFieldStatus                = "status"
	ProductionOrderFieldNote                  = "note"
	ProductionOrderFieldCreatedBy             = "created_by"
	ProductionOrderFieldCreatedAt             = "created_at"
	ProductionOrderFieldUpdatedAt             = "updated_at"
	ProductionOrderFieldDeletedAt             = "deleted_at"
	ProductionOrderFieldName                  = "name"
	ProductionOrderFieldEstimatedStartAt      = "estimated_start_at"
	ProductionOrderFieldEstimatedCompleteAt   = "estimated_complete_at"
	ProductionOrderFieldVersion               = "version"
	ProductionOrderFieldData                  = "data"
	ProductionOrderFieldProductionPlanID      = "production_plan_id"
	ProductionOrderFieldOrderID               = "order_id"
)

type ProductionOrder struct {
	ID                    string                     `db:"id"`
	ProductCode           string                     `db:"product_code"`
	ProductName           string                     `db:"product_name"`
	CustomerID            string                     `db:"customer_id"`
	SalesID               string                     `db:"sales_id"`
	QtyPaper              int64                      `db:"qty_paper"`
	QtyFinished           int64                      `db:"qty_finished"`
	QtyDelivered          int64                      `db:"qty_delivered"`
	PlannedProductionDate time.Time                  `db:"planned_production_date"`
	DeliveryDate          time.Time                  `db:"delivery_date"`
	DeliveryImage         sql.NullString             `db:"delivery_image"`
	Status                enum.ProductionOrderStatus `db:"status"`
	Note                  sql.NullString             `db:"note"`
	CreatedBy             string                     `db:"created_by"`
	CreatedAt             time.Time                  `db:"created_at"`
	UpdatedAt             time.Time                  `db:"updated_at"`
	DeletedAt             sql.NullTime               `db:"deleted_at"`
	Name                  string                     `db:"name"`
	EstimatedStartAt      sql.NullTime               `db:"estimated_start_at"`
	EstimatedCompleteAt   sql.NullTime               `db:"estimated_complete_at"`
	Version               int64                      `db:"version"`
	Data                  any                        `db:"data"`
	ProductionPlanID      sql.NullString             `db:"production_plan_id"`
	OrderID               sql.NullString             `db:"order_id"`
}

func (rcv *ProductionOrder) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		ProductionOrderFieldID,
		ProductionOrderFieldProductCode,
		ProductionOrderFieldProductName,
		ProductionOrderFieldCustomerID,
		ProductionOrderFieldSalesID,
		ProductionOrderFieldQtyPaper,
		ProductionOrderFieldQtyFinished,
		ProductionOrderFieldQtyDelivered,
		ProductionOrderFieldPlannedProductionDate,
		ProductionOrderFieldDeliveryDate,
		ProductionOrderFieldDeliveryImage,
		ProductionOrderFieldStatus,
		ProductionOrderFieldNote,
		ProductionOrderFieldCreatedBy,
		ProductionOrderFieldCreatedAt,
		ProductionOrderFieldUpdatedAt,
		ProductionOrderFieldDeletedAt,
		ProductionOrderFieldName,
		ProductionOrderFieldEstimatedStartAt,
		ProductionOrderFieldEstimatedCompleteAt,
		ProductionOrderFieldVersion,
		ProductionOrderFieldData,
		ProductionOrderFieldProductionPlanID,
		ProductionOrderFieldOrderID,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.ProductCode,
		&rcv.ProductName,
		&rcv.CustomerID,
		&rcv.SalesID,
		&rcv.QtyPaper,
		&rcv.QtyFinished,
		&rcv.QtyDelivered,
		&rcv.PlannedProductionDate,
		&rcv.DeliveryDate,
		&rcv.DeliveryImage,
		&rcv.Status,
		&rcv.Note,
		&rcv.CreatedBy,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
		&rcv.DeletedAt,
		&rcv.Name,
		&rcv.EstimatedStartAt,
		&rcv.EstimatedCompleteAt,
		&rcv.Version,
		&rcv.Data,
		&rcv.ProductionPlanID,
		&rcv.OrderID,
	}

	return
}

func (*ProductionOrder) TableName() string {
	return "production_orders"
}
