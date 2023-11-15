package production_order

import (
	"context"
	"fmt"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

func (c *productionOrderService) Delete(ctx context.Context, id string) error {
	// exec in transaction
	return cockroach.ExecInTx(ctx, func(ctx2 context.Context) error {
		err := c.productionOrderRepo.SoftDelete(ctx, id)
		if err != nil {
			return fmt.Errorf("c.productionOrderRepo.SoftDelete: %w", err)
		}
		// delete related production order stage
		err = c.productionOrderStageRepo.DeleteByProductionOrderID(ctx, id)
		if err != nil {
			return fmt.Errorf("c.productionOrderStageRepo.DeleteByProductionOrderID: %w", err)
		}
		// delete related custom field
		err = c.customFieldRepo.DeleteByEntity(ctx, enum.CustomFieldTypeProductionOrder, id)
		if err != nil {
			return fmt.Errorf("c.customFieldRepo.DeleteByEntity: %w", err)
		}
		return nil
	})

}
