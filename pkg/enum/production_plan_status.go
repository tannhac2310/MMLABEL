package enum

import (
	"bytes"
	"fmt"
)

type ProductionPlanStatus uint8

const (
	ProductionPlanStatus_Sale_CollectInfo ProductionPlanStatus = iota + 1
	ProductionPlanStatus_Develope_New
	ProductionPlanStatus_Develope_Feedback
	ProductionPlanStatus_Develope_Feedbacked
	ProductionPlanStatus_Develope_Done
	ProductionPlanStatus_Design_New
	ProductionPlanStatus_Design_Processing
	ProductionPlanStatus_Design_Done
	ProductionPlanStatus_Price_Processing
	ProductionPlanStatus_Price_Wait_For_Approval
	ProductionPlanStatus_Price_Done
	ProductionPlanStatus_PO_Processing
	ProductionPlanStatus_PO_Done
	ProductionPlanStatus_PO_Archived
)

const (
	ProductionPlanStageSale       = 1  // 00001
	ProductionPlanStageDevelop    = 2  // 00010
	ProductionPlanStageDesign     = 4  // 00100
	ProductionPlanStageFinance    = 8  // 01000
	ProductionPlanStageProduction = 16 // 10000
)

var ProductionPlanStatusName = map[ProductionPlanStatus]string{
	ProductionPlanStatus_Sale_CollectInfo:        "sale_collect_info",
	ProductionPlanStatus_Develope_New:            "develop_new",
	ProductionPlanStatus_Develope_Feedback:       "develop_feedback",
	ProductionPlanStatus_Develope_Feedbacked:     "develop_feedbacked",
	ProductionPlanStatus_Develope_Done:           "develop_done",
	ProductionPlanStatus_Design_New:              "design_new",
	ProductionPlanStatus_Design_Processing:       "design_processing",
	ProductionPlanStatus_Design_Done:             "design_done",
	ProductionPlanStatus_Price_Processing:        "price_processing",
	ProductionPlanStatus_Price_Wait_For_Approval: "price_wait_for_approval",
	ProductionPlanStatus_Price_Done:              "price_done",
	ProductionPlanStatus_PO_Processing:           "po_processing",
	ProductionPlanStatus_PO_Done:                 "po_done",
	ProductionPlanStatus_PO_Archived:             "po_archived",
}

var ProductionPlanStatusSage = map[ProductionPlanStatus]int{
	ProductionPlanStatus_Sale_CollectInfo:        ProductionPlanStageSale,
	ProductionPlanStatus_Develope_New:            ProductionPlanStageDevelop,
	ProductionPlanStatus_Develope_Feedback:       ProductionPlanStageDevelop,
	ProductionPlanStatus_Develope_Feedbacked:     ProductionPlanStageDevelop,
	ProductionPlanStatus_Develope_Done:           ProductionPlanStageDevelop,
	ProductionPlanStatus_Design_New:              ProductionPlanStageDesign,
	ProductionPlanStatus_Design_Processing:       ProductionPlanStageDesign,
	ProductionPlanStatus_Design_Done:             ProductionPlanStageDesign,
	ProductionPlanStatus_Price_Processing:        ProductionPlanStageDesign,
	ProductionPlanStatus_Price_Wait_For_Approval: ProductionPlanStageFinance,
	ProductionPlanStatus_Price_Done:              ProductionPlanStageFinance,
	ProductionPlanStatus_PO_Processing:           ProductionPlanStageProduction,
	ProductionPlanStatus_PO_Done:                 ProductionPlanStageProduction,
	ProductionPlanStatus_PO_Archived:             ProductionPlanStageProduction,
}

var ProductionPlanStatusValue = func() map[string]ProductionPlanStatus {
	value := map[string]ProductionPlanStatus{}
	for k, v := range ProductionPlanStatusName {
		value[v] = k
		value[fmt.Sprintf("%v", k)] = k
	}

	return value
}()

func (e ProductionPlanStatus) MarshalJSON() ([]byte, error) {
	v, ok := ProductionPlanStatusName[e]
	if !ok {
		return nil, fmt.Errorf("invalid values of ProductionPlanStatus")
	}

	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(v)
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (e *ProductionPlanStatus) UnmarshalJSON(data []byte) error {
	data = bytes.Trim(data, "\"")
	v, ok := ProductionPlanStatusValue[string(data)]
	if !ok {
		return fmt.Errorf("enum '%s' is not register, must be one of: %v", data, e.EnumDescriptions())
	}

	*e = v

	return nil
}

func (*ProductionPlanStatus) EnumDescriptions() []string {
	vals := []string{}

	for _, name := range ProductionPlanStatusName {
		vals = append(vals, name)
	}

	return vals
}
