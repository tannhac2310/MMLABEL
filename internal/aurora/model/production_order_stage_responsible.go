package model

const (
	ProductionOrderStageResponsibleFieldID              = "id"
	ProductionOrderStageResponsibleFieldPOStageDeviceID = "po_stage_device_id"
	ProductionOrderStageResponsibleFieldUserID          = "user_id"
)

type ProductionOrderStageResponsible struct {
	ID              string `db:"id"`
	POStageDeviceID string `db:"po_stage_device_id"`
	UserID          string `db:"user_id"`
}

func (rcv *ProductionOrderStageResponsible) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		ProductionOrderStageResponsibleFieldID,
		ProductionOrderStageResponsibleFieldPOStageDeviceID,
		ProductionOrderStageResponsibleFieldUserID,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.POStageDeviceID,
		&rcv.UserID,
	}

	return
}

func (*ProductionOrderStageResponsible) TableName() string {
	return "production_order_stage_responsible"
}
