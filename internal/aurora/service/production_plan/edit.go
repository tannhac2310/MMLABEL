package production_plan

import (
	"context"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

type EditProductionPlanOpts struct {
	ID          string
	Name        string
	CustomerID  string
	SalesID     string
	Thumbnail   string
	Status      enum.ProductionPlanStatus
	Note        string
	CustomField []*CustomField
	CreatedBy   string
}

func (c *productionPlanService) EditProductionPlan(ctx context.Context, opt *EditProductionPlanOpts) error {
	return nil
}
