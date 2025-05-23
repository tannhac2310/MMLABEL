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
	// find last doing stage
	lastDoingStages, err := c.productionOrderStageRepo.Search(ctx, &repository.SearchProductionOrderStagesOpts{
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

	// update complete stage to delivery
	// update none stage to reception
	err = cockroach.ExecInTx(ctx, func(ctx2 context.Context) error {
		if len(lastDoingStages) > 0 {
			model := lastDoingStages[0]
			model.Status = enum.ProductionOrderStageStatusProductDelivery
			model.ProductDeliveryAt = cockroach.Time(time.Now()) // cap nhat ngay giao san pham
			err = c.productionOrderStageRepo.Update(ctx2, model)
			if err != nil {
				return fmt.Errorf("c.productionOrderStageRepo.Update: %w", err)
			}
		}
		if len(firstNoneStages) > 0 {
			model := firstNoneStages[0]
			model.ReceptionAt = cockroach.Time(time.Now()) // cap nhat ngay nhan san pham
			model.Status = enum.ProductionOrderStageStatusReception
			err = c.productionOrderStageRepo.Update(ctx2, model)
			model.ReceptionAt = cockroach.Time(time.Now()) // cap nhat ngay nhan san pham
			if err != nil {
				return fmt.Errorf("c.productionOrderStageRepo.Update: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("cockroach.ExecInTx: %w", err)
	}
	return nil
}
