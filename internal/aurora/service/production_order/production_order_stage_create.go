package production_order

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
)

func (c *productionOrderService) CreateProductionOrderStage(ctx context.Context, poId string, opt *ProductionOrderStage) (string, error) {
	id := idutil.ULIDNow()
	err := c.productionOrderStageRepo.Insert(ctx, &model.ProductionOrderStage{
		ID:                  id,
		ProductionOrderID:   poId,
		StageID:             opt.StageID,
		EstimatedStartAt:    cockroach.Time(opt.EstimatedStartAt),
		EstimatedCompleteAt: cockroach.Time(opt.EstimatedCompleteAt),
		StartedAt:           cockroach.Time(opt.StartedAt),
		CompletedAt:         cockroach.Time(opt.CompletedAt),
		Status:              opt.Status,
		Condition:           cockroach.String(opt.Condition),
		Note:                cockroach.String(opt.Note),
		Data:                opt.Data,
	})
	if err != nil {
		return "", err
	}
	return id, nil
}
func (c *productionOrderService) EditProductionOrderStage(ctx context.Context, opt *ProductionOrderStage) error {
	table := model.ProductionOrderStage{}

	updater := cockroach.NewUpdater(table.TableName(), model.ProductionOrderStageFieldID, opt.ID)
	updater.Set(model.ProductionOrderStageFieldEstimatedStartAt, cockroach.Time(opt.EstimatedStartAt))
	if opt.StageID != "" {
		updater.Set(model.ProductionOrderStageFieldStageID, opt.StageID)
	}
	updater.Set(model.ProductionOrderStageFieldEstimatedCompleteAt, cockroach.Time(opt.EstimatedCompleteAt))
	updater.Set(model.ProductionOrderStageFieldStartedAt, cockroach.Time(opt.StartedAt))
	updater.Set(model.ProductionOrderStageFieldCompletedAt, cockroach.Time(opt.CompletedAt))
	updater.Set(model.ProductionOrderStageFieldStatus, opt.Status)
	updater.Set(model.ProductionOrderStageFieldCondition, cockroach.String(opt.Condition))
	updater.Set(model.ProductionOrderStageFieldNote, cockroach.String(opt.Note))
	updater.Set(model.ProductionOrderStageFieldData, opt.Data)
	updater.Set(model.ProductionOrderStageFieldSorting, opt.Sorting)

	if opt.Status == enum.ProductionOrderStageStatusProductionStart {
		updater.Set(model.ProductionOrderStageFieldProductionStartAt, cockroach.Time(time.Now()))
	}
	if opt.Status == enum.ProductionOrderStageStatusProductionCompletion {
		updater.Set(model.ProductionOrderStageFieldProductionCompletionAt, cockroach.Time(time.Now()))
	}

	err := cockroach.UpdateFields(ctx, updater)
	if err != nil {
		return fmt.Errorf("update productionOrderStage failed %w", err)
	}
	return nil
}
func (c *productionOrderService) DeleteProductionOrderStage(ctx context.Context, id string) error {
	return c.productionOrderStageRepo.SoftDelete(ctx, id)
}
