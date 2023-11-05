package model

import (
	"database/sql"
	"time"
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
	ProductionOrderFieldDeliveredImage        = "delivered_image"
	ProductionOrderFieldStatus                = "status"
	ProductionOrderFieldNote                  = "note"
	ProductionOrderFieldCreatedBy             = "created_by"
	ProductionOrderFieldCreatedAt             = "created_at"
	ProductionOrderFieldUpdatedAt             = "updated_at"
	ProductionOrderFieldDeletedAt             = "deleted_at"
)

type ProductionOrder struct {
	ID                    string         `db:"id"`
	ProductCode           string         `db:"product_code"`
	ProductName           string         `db:"product_name"`
	CustomerID            string         `db:"customer_id"`
	SalesID               string         `db:"sales_id"`
	QtyPaper              int64          `db:"qty_paper"`
	QtyFinished           int64          `db:"qty_finished"`
	QtyDelivered          int64          `db:"qty_delivered"`
	PlannedProductionDate time.Time      `db:"planned_production_date"`
	DeliveryDate          time.Time      `db:"delivery_date"`
	DeliveredImage        sql.NullString `db:"delivered_image"`
	Status                int16          `db:"status"`
	Note                  sql.NullString `db:"note"`
	CreatedBy             string         `db:"created_by"`
	CreatedAt             time.Time      `db:"created_at"`
	UpdatedAt             time.Time      `db:"updated_at"`
	DeletedAt             sql.NullTime   `db:"deleted_at"`
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
		ProductionOrderFieldDeliveredImage,
		ProductionOrderFieldStatus,
		ProductionOrderFieldNote,
		ProductionOrderFieldCreatedBy,
		ProductionOrderFieldCreatedAt,
		ProductionOrderFieldUpdatedAt,
		ProductionOrderFieldDeletedAt,
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
		&rcv.DeliveredImage,
		&rcv.Status,
		&rcv.Note,
		&rcv.CreatedBy,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
		&rcv.DeletedAt,
	}

	return
}

func (*ProductionOrder) TableName() string {
	return "production_orders"
}
