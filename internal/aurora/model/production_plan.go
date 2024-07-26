package model

import (
	"database/sql"
	"time"
)

const (
	ProductionPlanFieldID         = "id"
	ProductionPlanFieldCustomerID = "customer_id"
	ProductionPlanFieldSalesID    = "sales_id"
	ProductionPlanFieldThumbnail  = "thumbnail"
	ProductionPlanFieldStatus     = "status"
	ProductionPlanFieldNote       = "note"
	ProductionPlanFieldCreatedBy  = "created_by"
	ProductionPlanFieldCreatedAt  = "created_at"
	ProductionPlanFieldUpdatedAt  = "updated_at"
	ProductionPlanFieldDeletedAt  = "deleted_at"
	ProductionPlanFieldName       = "name"
)

type ProductionPlan struct {
	ID         string         `db:"id"`
	CustomerID string         `db:"customer_id"`
	SalesID    string         `db:"sales_id"`
	Thumbnail  sql.NullString `db:"thumbnail"`
	Status     int16          `db:"status"`
	Note       sql.NullString `db:"note"`
	CreatedBy  string         `db:"created_by"`
	CreatedAt  time.Time      `db:"created_at"`
	UpdatedAt  time.Time      `db:"updated_at"`
	DeletedAt  sql.NullTime   `db:"deleted_at"`
	Name       string         `db:"name"`
}

func (rcv *ProductionPlan) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		ProductionPlanFieldID,
		ProductionPlanFieldCustomerID,
		ProductionPlanFieldSalesID,
		ProductionPlanFieldThumbnail,
		ProductionPlanFieldStatus,
		ProductionPlanFieldNote,
		ProductionPlanFieldCreatedBy,
		ProductionPlanFieldCreatedAt,
		ProductionPlanFieldUpdatedAt,
		ProductionPlanFieldDeletedAt,
		ProductionPlanFieldName,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.CustomerID,
		&rcv.SalesID,
		&rcv.Thumbnail,
		&rcv.Status,
		&rcv.Note,
		&rcv.CreatedBy,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
		&rcv.DeletedAt,
		&rcv.Name,
	}

	return
}

func (*ProductionPlan) TableName() string {
	return "production_plans"
}
