package model

const (
	ProductionPlanDetailFieldID                        = "id"
	ProductionPlanDetailFieldProductionPlanID          = "production_plan_id"
	ProductionPlanDetailFieldProductionPlanAttributeID = "production_plan_attribute_id"
)

type ProductionPlanDetail struct {
	ID                        string `db:"id"`
	ProductionPlanID          string `db:"production_plan_id"`
	ProductionPlanAttributeID string `db:"production_plan_attribute_id"`
}

func (rcv *ProductionPlanDetail) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		ProductionPlanDetailFieldID,
		ProductionPlanDetailFieldProductionPlanID,
		ProductionPlanDetailFieldProductionPlanAttributeID,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.ProductionPlanID,
		&rcv.ProductionPlanAttributeID,
	}

	return
}

func (*ProductionPlanDetail) TableName() string {
	return "production_plan_details"
}
