package production_plan

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

func (c *productionPlanService) UpdateCurrentStage(ctx context.Context, productionPlanID string, stageID enum.ProductionPlanStage) error {
	plan, err := c.productionPlanRepo.FindByID(ctx, productionPlanID)
	if err != nil {
		return fmt.Errorf("find production plan failed: %w", err)
	}

	plan.CurrentStage = stageID
	if err := c.productionPlanRepo.Update(ctx, plan.ProductionPlan); err != nil {
		return fmt.Errorf("update production plan failed: %w", err)
	}

	return nil
}
