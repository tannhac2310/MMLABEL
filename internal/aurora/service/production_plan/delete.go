package production_plan

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

func (c *productionPlanService) DeleteProductionPlan(ctx context.Context, id string) error {
	// exec in transaction
	return cockroach.ExecInTx(ctx, func(ctx2 context.Context) error {
		err := c.productionPlanRepo.SoftDelete(ctx2, id)
		if err != nil {
			return fmt.Errorf("c.productionPlanRepo.SoftDelete: %w", err)
		}

		err = c.customFieldRepo.DeleteByEntity(ctx, enum.CustomFieldTypeProductionPlan, id)
		if err != nil {
			return fmt.Errorf("c.customFieldRepo.DeleteByEntity: %w", err)
		}

		return nil
	})
}
