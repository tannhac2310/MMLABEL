package enum

import (
	"bytes"
	"fmt"
)

type ProductionPlanStatus uint8

const (
	// sale
	ProductionPlanStatusSaleNew = 1

	ProductionPlanStatusAandDNew            = 21
	ProductionPlanStatusAandDWaitingComment = 22
	ProductionPlanStatusAandDCommented      = 23
	ProductionPlanStatusAandDFinished       = 24

	ProductionPlanStatusDesignNew       = 35
	ProductionPlanStatusDesignDesigning = 36
	ProductionPlanStatusDesignCompleted = 37

	ProductionPlanStatusPOProcessing = 41
	ProductionPlanStatusPOCompleted  = 42

	ProductionPlanStatusArchivedNew = 51
)

var ProductionPlanStatusName = map[ProductionPlanStatus]string{
	ProductionPlanStatusSaleNew: "sale_new",

	// a_and_d
	ProductionPlanStatusAandDNew:            "a_and_d_new",
	ProductionPlanStatusAandDWaitingComment: "a_and_d_waiting_comment",
	ProductionPlanStatusAandDCommented:      "a_and_d_commented",
	ProductionPlanStatusAandDFinished:       "a_and_d_finished",
	// design
	ProductionPlanStatusDesignNew:       "design_new",
	ProductionPlanStatusDesignDesigning: "design_designing",
	ProductionPlanStatusDesignCompleted: "design_completed",
	// po
	ProductionPlanStatusPOProcessing: "po_processing",
	ProductionPlanStatusPOCompleted:  "po_completed",
	// archived
	ProductionPlanStatusArchivedNew: "archived_new",
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
