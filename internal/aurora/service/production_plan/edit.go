package production_plan

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
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
	now := time.Now()

	plan, err := c.productionPlanRepo.FindByID(ctx, opt.ID)
	if err != nil {
		return err
	} else if !plan.Editable() {
		return fmt.Errorf("production plan cannot be edit")
	}

	execTx := cockroach.ExecInTx(ctx, func(ctx2 context.Context) error {
		plan.UpdatedBy = opt.CreatedBy
		plan.UpdatedAt = now
		plan.Thumbnail = cockroach.String(opt.Thumbnail)
		plan.Note = cockroach.String(opt.Note)
		if !plan.CanChangeStatusTo(opt.Status) {
			return fmt.Errorf("cannot change production plan status to %s", enum.ProductionPlanStatusName[opt.Status])
		}

		if err := c.productionPlanRepo.Update(ctx2, plan.ProductionPlan); err != nil {
			return fmt.Errorf("update production plan failed: %w", err)
		}

		if err := c.customFieldRepo.DeleteByEntity(ctx, enum.CustomFieldTypeProductionPlan, plan.ID); err != nil {
			return fmt.Errorf("override existing custom fields for production plan failed: %w", err)
		}

		for _, field := range opt.CustomField {
			if err := c.customFieldRepo.Insert(ctx2, &model.CustomField{
				ID:         idutil.ULIDNow(),
				EntityID:   plan.ID,
				EntityType: enum.CustomFieldTypeProductionPlan,
				Field:      field.Field,
				Value:      field.Value,
			}); err != nil {
				return fmt.Errorf("create new custom field for production plan failed: %w", err)
			}
		}

		return nil
	})
	if execTx != nil {
		return execTx
	}

	return nil
}
