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

type ProcessProductionOrderOpts struct {
	ID                  string
	Stages              []*ProductionOrderStage
	EstimatedStartAt    time.Time
	EstimatedCompleteAt time.Time
	CreatedBy           string
}

type ProductionOrderStage struct {
	StageID             string
	EstimatedStartAt    time.Time
	EstimatedCompleteAt time.Time
	StartedAt           time.Time
	CompletedAt         time.Time
	Status              enum.ProductionOrderStageStatus
	Condition           string
	Note                string
	Data                map[string]interface{}
	ID                  string
	Sorting             int16
}

func (c *productionPlanService) ProcessProductionOrder(ctx context.Context, opt *ProcessProductionOrderOpts) (string, error) {
	now := time.Now()
	id := idutil.ULIDNow()

	plan, err := c.productionPlanRepo.FindByID(ctx, opt.ID)
	if err != nil {
		return "", fmt.Errorf("not found any plans with id %s: %w", opt.ID, err)
	}

	customFields, err := c.customFieldRepo.Search(ctx, &repository.SearchCustomFieldsOpts{
		EntityId:   plan.ID,
		EntityType: enum.CustomFieldTypeProductionPlan,
		Limit:      1000,
		Offset:     0,
	})
	if err != nil {
		return "", fmt.Errorf("query custom fields for production plan %s failed: %w", plan.ID, err)
	}
	customFieldMap := generic.ToMap(customFields, func(f *repository.CustomFieldData) string {
		return f.Field
	})

	productionOrder := &model.ProductionOrder{
		ID:                  id,
		ProductCode:         customFieldMap[ProductionPlanCustomField_ma_sp].Value,
		ProductName:         plan.ProductionPlan.Name,
		CustomerID:          plan.CustomerID,
		SalesID:             plan.SalesID,
		QtyPaper:            0, // FIXME
		QtyFinished:         0,
		QtyDelivered:        0,
		DeliveryDate:        time.Time{},
		DeliveryImage:       plan.Thumbnail,
		Status:              enum.ProductionOrderStatusWaiting,
		Note:                plan.Note,
		CreatedBy:           opt.CreatedBy,
		CreatedAt:           now,
		UpdatedAt:           now,
		Name:                plan.Name,
		EstimatedStartAt:    cockroach.Time(opt.EstimatedStartAt),
		EstimatedCompleteAt: cockroach.Time(opt.EstimatedCompleteAt),
	}

	errTx := cockroach.ExecInTx(ctx, func(ctx2 context.Context) error {
		err := c.productionOrderRepo.Insert(ctx2, productionOrder)
		if err != nil {
			return fmt.Errorf("c.productionOrderRepo.Insert: %w", err)
		}
		for _, orderStage := range opt.Stages {
			err = c.productionOrderStageRepo.Insert(ctx2, &model.ProductionOrderStage{
				ID:                  idutil.ULIDNow(),
				ProductionOrderID:   id,
				Sorting:             orderStage.Sorting,
				StageID:             orderStage.StageID,
				EstimatedStartAt:    cockroach.Time(orderStage.EstimatedStartAt),
				EstimatedCompleteAt: cockroach.Time(orderStage.EstimatedCompleteAt),
				StartedAt:           cockroach.Time(orderStage.StartedAt),
				CompletedAt:         cockroach.Time(orderStage.CompletedAt),
				Status:              orderStage.Status,
				Condition:           cockroach.String(orderStage.Condition),
				Note:                cockroach.String(orderStage.Note),
				Data:                orderStage.Data,
				CreatedAt:           now,
				UpdatedAt:           now,
			})
			if err != nil {
				return fmt.Errorf("add order stage: %w", err)
			}
		}
		for _, val := range customFields {
			// add custom field table
			err = c.customFieldRepo.Insert(ctx2, &model.CustomField{
				ID:         idutil.ULIDNow(),
				EntityType: enum.CustomFieldTypeProductionOrder,
				EntityID:   id,
				Field:      val.Field,
				Value:      val.Value,
			})
			if err != nil {
				return fmt.Errorf("c.customFieldRepo.Insert: %w", err)
			}
		}
		return nil
	})
	if errTx != nil {
		return "", errTx
	}
	return productionOrder.ID, nil
}
