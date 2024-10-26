package model

import "time"

const (
	ProductionPlanProductFieldID               = "id"
	ProductionPlanProductFieldProductionPlanID = "production_plan_id"
	ProductionPlanProductFieldProductID        = "product_id"
	ProductionPlanProductFieldQuantity         = "quantity"
	ProductionPlanProductFieldCreatedAt        = "created_at"
	ProductionPlanProductFieldUpdatedAt        = "updated_at"
	ProductionPlanProductFieldCreatedBy        = "created_by"
	ProductionPlanProductFieldUpdatedBy        = "updated_by"
)

type ProductionPlanProduct struct {
	ID               string    `db:"id"`
	ProductionPlanID string    `db:"production_plan_id"`
	ProductID        string    `db:"product_id"`
	Quantity         int64     `db:"quantity"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
	CreatedBy        string    `db:"created_by"`
	UpdatedBy        string    `db:"updated_by"`
}

func (rcv *ProductionPlanProduct) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		ProductionPlanProductFieldID,
		ProductionPlanProductFieldProductionPlanID,
		ProductionPlanProductFieldProductID,
		ProductionPlanProductFieldQuantity,
		ProductionPlanProductFieldCreatedAt,
		ProductionPlanProductFieldUpdatedAt,
		ProductionPlanProductFieldCreatedBy,
		ProductionPlanProductFieldUpdatedBy,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.ProductionPlanID,
		&rcv.ProductID,
		&rcv.Quantity,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
		&rcv.CreatedBy,
		&rcv.UpdatedBy,
	}

	return
}

func (*ProductionPlanProduct) TableName() string {
	return "production_plan_products"
}
