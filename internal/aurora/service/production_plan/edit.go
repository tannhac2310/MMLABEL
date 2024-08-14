package production_plan

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/generic"
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
	}

	customFields, err := c.customFieldRepo.Search(ctx, &repository.SearchCustomFieldsOpts{
		EntityId:   plan.ID,
		EntityType: enum.CustomFieldTypeProductionPlan,
		Limit:      1000,
		Offset:     0,
	})
	if err != nil {
		return err
	}
	requiredCustomFields := c.GetCustomField()
	for _, val := range opt.CustomField {
		if _, ok := requiredCustomFields[val.Field]; !ok {
			return fmt.Errorf("thông tin %s không hợp lệ", val.Field)
		}
	}

	customFieldMap := generic.ToMap(customFields, func(f *repository.CustomFieldData) string {
		return f.Field
	})

	newCustomFields := make([]*model.CustomField, 0)
	updatedCustomFields := make([]*model.CustomField, 0)
	for _, field := range opt.CustomField {
		if _, ok := customFieldMap[field.Field]; ok {
			customFieldMap[field.Field].Value = field.Value
			updatedCustomFields = append(updatedCustomFields, customFieldMap[field.Field].CustomField)
			delete(customFieldMap, field.Field)
		} else {
			newCustomFields = append(newCustomFields, &model.CustomField{
				ID:         idutil.ULIDNow(),
				EntityID:   plan.ID,
				EntityType: enum.CustomFieldTypeProductionPlan,
				Field:      field.Field,
				Value:      field.Value,
			})
		}
	}
	deletedCustomFields := make([]*model.CustomField, 0)
	for _, field := range customFieldMap {
		deletedCustomFields = append(deletedCustomFields, field.CustomField)
	}

	execTx := cockroach.ExecInTx(ctx, func(ctx2 context.Context) error {
		plan.UpdatedBy = opt.CreatedBy
		plan.UpdatedAt = now
		plan.Thumbnail = cockroach.String(opt.Thumbnail)
		plan.Note = cockroach.String(opt.Note)
		if !plan.CanChangeStatusTo(opt.Status) {
			return fmt.Errorf("cannot change production plan status to %s", enum.ProductionPlanStatusName[opt.Status])
		}
		plan.Status = opt.Status
		currentStage := enum.ProductionPlanStatusSage[plan.Status]
		plan.CurrentStage = int(enum.ProductionPlanStageSale) | currentStage // only sale and current stage can view the production plan

		if err := c.productionPlanRepo.Update(ctx2, plan.ProductionPlan); err != nil {
			return fmt.Errorf("update production plan failed: %w", err)
		}

		for _, field := range updatedCustomFields {
			if err := c.customFieldRepo.Update(ctx, field); err != nil {
				return fmt.Errorf("update custom field for production plan failed: %w", err)
			}
		}
		for _, field := range newCustomFields {
			if err := c.customFieldRepo.Insert(ctx, field); err != nil {
				return fmt.Errorf("create custom field for production plan failed: %w", err)
			}
		}
		for _, field := range deletedCustomFields {
			if err := c.customFieldRepo.Delete(ctx, field.ID); err != nil {
				return fmt.Errorf("delete custom field for production plan failed: %w", err)
			}
		}

		return nil
	})
	if execTx != nil {
		return execTx
	}

	return nil
}
