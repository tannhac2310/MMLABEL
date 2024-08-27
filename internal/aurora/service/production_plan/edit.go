package production_plan

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/generic"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
)

type EditProductionPlanOpts struct {
	ID           string
	Name         string
	CustomerID   string
	SalesID      string
	ProductName  string
	ProductCode  string
	QtyPaper     int64
	QtyFinished  int64
	QtyDelivered int64
	Thumbnail    string
	Status       enum.ProductionPlanStatus
	Note         string
	CustomField  []*CustomField
	Workflow     any
	CreatedBy    string
}

func (c *productionPlanService) EditProductionPlan(ctx context.Context, opt *EditProductionPlanOpts) error {
	now := time.Now()

	plan, err := c.productionPlanRepo.FindByID(ctx, opt.ID)
	if err != nil {
		return err
	}
	if !plan.ProductionPlan.Editable() {
		return fmt.Errorf("không thể chỉnh sửa kế hoạch đã được đưa vào sản xuất")
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
		plan.ProductionPlan.Name = opt.Name
		plan.ProductionPlan.CustomerID = opt.CustomerID
		plan.ProductionPlan.SalesID = opt.SalesID
		plan.ProductionPlan.ProductName = opt.ProductName
		plan.ProductionPlan.ProductCode = opt.ProductCode
		plan.ProductionPlan.Workflow = opt.Workflow

		plan.UpdatedBy = opt.CreatedBy
		plan.UpdatedAt = now
		plan.QtyPaper = opt.QtyPaper
		plan.QtyFinished = opt.QtyFinished
		plan.QtyDelivered = opt.QtyDelivered
		plan.Thumbnail = cockroach.String(opt.Thumbnail)
		plan.Note = cockroach.String(opt.Note)
		plan.Status = opt.Status
		//plan.CurrentStage = opt.CurrentStage // TODO: have to update this field

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
			// if field.key contains sale_survey then ignore delete

			if strings.Contains(field.Field, "sale_survey") {
				//continue
			}

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

func (c *productionPlanService) UpdateCustomFields(ctx context.Context, productionPlanID string, fieldValues []*CustomField) error {
	plan, err := c.productionPlanRepo.FindByID(ctx, productionPlanID)
	if err != nil {
		return err
	}
	if !plan.ProductionPlan.Editable() {
		return fmt.Errorf("không thể chỉnh sửa kế hoạch đã được đưa vào sản xuất")
	}
	errTx := cockroach.ExecInTx(ctx, func(ctx2 context.Context) error {

		customFields, err := c.customFieldRepo.Search(ctx, &repository.SearchCustomFieldsOpts{
			EntityId:   plan.ID,
			EntityType: enum.CustomFieldTypeProductionPlan,
			Limit:      1000,
			Offset:     0,
		})

		if err != nil {
			return fmt.Errorf("search custom fields failed: %w", err)
		}

		customFieldMap := generic.ToMap(customFields, func(f *repository.CustomFieldData) string {
			return f.Field
		})

		newCustomFields := make([]*model.CustomField, 0)
		updatedCustomFields := make([]*model.CustomField, 0)
		for _, field := range fieldValues {
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

		return nil

	})

	if errTx != nil {
		return fmt.Errorf("update custom fields failed: %w", errTx)
	}

	return nil
}
