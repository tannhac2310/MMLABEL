package production_order

import (
	"context"
	"fmt"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"time"
)

func (c *productionOrderService) AcceptAndChangeNextStage(ctx context.Context, id string) error {
	// find all stage of
	lastDoingStages, err := c.productionOrderStageRepo.Search(ctx, &repository.SearchProductionOrderStagesOpts{
		IDs:                        nil,
		ProductionOrderID:          id,
		ProductionOrderStageStatus: enum.ProductionOrderStageStatusProductionCompletion,
		Limit:                      1,
		Offset:                     0,
		Sort: &repository.Sort{
			Order: "DESC",
			By:    "sorting",
		},
	})
	if err != nil {
		return err
	}

	firstNoneStages, err := c.productionOrderStageRepo.Search(ctx, &repository.SearchProductionOrderStagesOpts{
		IDs:                        nil,
		ProductionOrderID:          id,
		ProductionOrderStageStatus: enum.ProductionOrderStageStatusWaiting,
		Limit:                      1,
		Offset:                     0,
		Sort: &repository.Sort{
			Order: "DESC",
			By:    "sorting",
		},
	})
	if err != nil {
		return err
	}

	// update doing stage to done
	// update none stage to doing
	err = cockroach.ExecInTx(ctx, func(ctx2 context.Context) error {
		if len(lastDoingStages) > 0 {
			model := lastDoingStages[0]
			model.Status = enum.ProductionOrderStageStatusProductDelivery
			model.ProductionCompletionAt = cockroach.Time(time.Now()) // cap nhat ngay hoan thanh ban giao san pham
			err = c.productionOrderStageRepo.Update(ctx2, model)
			if err != nil {
				return fmt.Errorf("c.productionOrderStageRepo.Update: %w", err)
			}
		}

		model := firstNoneStages[0]
		model.Status = enum.ProductionOrderStageStatusReception
		err = c.productionOrderStageRepo.Update(ctx2, model)
		model.ReceptionAt = cockroach.Time(time.Now()) // cap nhat ngay nhan san pham
		if err != nil {
			return fmt.Errorf("c.productionOrderStageRepo.Update: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("cockroach.ExecInTx: %w", err)
	}
	return nil
}
