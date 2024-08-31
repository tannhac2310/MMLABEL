package enum

import (
	"bytes"
	"fmt"
)

type CommentTarget uint8

const (
	CommentTarget_Unknown CommentTarget = iota
	CommentTarget_ProductionPlanStageSale
	CommentTarget_ProductionPlanStageRandD
	CommentTarget_ProductionPlanStageDesign
	CommentTarget_ProductionPlanStagePrice
)

var CommentTargetName = map[CommentTarget]string{
	CommentTarget_ProductionPlanStageSale:   "pl_sale",
	CommentTarget_ProductionPlanStageRandD:  "pl_r_and_d",
	CommentTarget_ProductionPlanStageDesign: "pl_design",
	CommentTarget_ProductionPlanStagePrice:  "pl_price",
}

var CommentTargetValue = func() map[string]CommentTarget {
	value := map[string]CommentTarget{}
	for k, v := range CommentTargetName {
		value[v] = k
		value[fmt.Sprintf("%v", k)] = k
	}

	return value
}()

func (e CommentTarget) MarshalJSON() ([]byte, error) {
	v, ok := CommentTargetName[e]
	if !ok {
		return nil, fmt.Errorf("invalid values of CommentTarget")
	}

	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(v)
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (e *CommentTarget) UnmarshalJSON(data []byte) error {
	data = bytes.Trim(data, "\"")
	v, ok := CommentTargetValue[string(data)]
	if !ok {
		return fmt.Errorf("enum '%s' is not register, must be one of: %v", data, e.EnumDescriptions())
	}

	*e = v

	return nil
}

func (*CommentTarget) EnumDescriptions() []string {
	vals := []string{}

	for _, name := range CommentTargetName {
		vals = append(vals, name)
	}

	return vals
}
