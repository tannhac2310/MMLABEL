package enum

import (
	"bytes"
	"fmt"
)

type ProductionPlanStage uint8

const (
	ProductionPlanStageSale ProductionPlanStage = iota + 1
	ProductionPlanStageRandD
	ProductionPlanStageDesign
	ProductionPlanStagePOProcessing
	ProductionPlanStagePOCompleted
	ProductionPlanStagePOArchived
)

var ProductionPlanStageName = map[ProductionPlanStage]string{
	ProductionPlanStageSale:         "sale",
	ProductionPlanStageRandD:        "r_and_d",
	ProductionPlanStageDesign:       "design",
	ProductionPlanStagePOProcessing: "po_processing",
	ProductionPlanStagePOCompleted:  "po_completed",
	ProductionPlanStagePOArchived:   "archived",
}

var ProductionPlanStageValue = func() map[string]ProductionPlanStage {
	value := map[string]ProductionPlanStage{}
	for k, v := range ProductionPlanStageName {
		value[v] = k
		value[fmt.Sprintf("%v", k)] = k
	}

	return value
}()

func (e ProductionPlanStage) MarshalJSON() ([]byte, error) {
	v, ok := ProductionPlanStageName[e]
	if !ok {
		return nil, fmt.Errorf("invalid values of ProductionPlanStage")
	}

	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(v)
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (e *ProductionPlanStage) UnmarshalJSON(data []byte) error {
	data = bytes.Trim(data, "\"")
	v, ok := ProductionPlanStageValue[string(data)]
	if !ok {
		return fmt.Errorf("enum '%s' is not register, must be one of: %v", data, e.EnumDescriptions())
	}

	*e = v

	return nil
}

func (*ProductionPlanStage) EnumDescriptions() []string {
	vals := []string{}

	for _, name := range ProductionPlanStageName {
		vals = append(vals, name)
	}

	return vals
}
