package production_plan

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
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
	newProductionOrderID := idutil.ULIDNow()

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

	deliveryDate := time.Time{}
	for _, val := range customFields {
		if val.Field == "delivery_date" && val.Value != "" {
			deliveryDate, _ = time.Parse(time.RFC3339, val.Value)
		}
	}
	productionOrder := &model.ProductionOrder{
		ID:                  newProductionOrderID,
		ProductCode:         plan.ProductCode,
		ProductName:         plan.ProductName,
		QtyPaper:            plan.QtyPaper,
		QtyFinished:         plan.QtyFinished,
		QtyDelivered:        plan.QtyDelivered,
		DeliveryDate:        deliveryDate,
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

		plan.PoStages = model.ProductionStageInfo{
			Items: make([]*model.ProductionStageItem, 0),
		}
		for _, orderStage := range opt.Stages {
			stageID := idutil.ULIDNow()

			// update stages of production plan
			plan.PoStages.Items = append(plan.PoStages.Items, &model.ProductionStageItem{
				ID:                  stageID,
				ProductionPlanID:    plan.ID,
				StageID:             orderStage.StageID,
				Note:                orderStage.Note,
				CreatedAt:           now,
				UpdatedAt:           now,
				EstimatedStartAt:    orderStage.EstimatedStartAt,
				EstimatedCompleteAt: orderStage.EstimatedCompleteAt,
				Sorting:             orderStage.Sorting,
			})

			err = c.productionOrderStageRepo.Insert(ctx2, &model.ProductionOrderStage{
				ID:                  stageID,
				ProductionOrderID:   newProductionOrderID,
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
				EntityID:   newProductionOrderID,
				Field:      val.Field,
				Value:      val.Value,
			})
			if err != nil {
				return fmt.Errorf("c.customFieldRepo.Insert: %w", err)
			}
		}

		// update production plan
		plan.CurrentStage = enum.ProductionPlanStagePOProcessing
		plan.Status = enum.ProductionPlanStatusPOProcessing
		plan.UpdatedAt = now
		if err := c.productionPlanRepo.Update(ctx, plan.ProductionPlan); err != nil {
			return fmt.Errorf("c.productionPlanRepo.Update: %w", err)
		}

		// copy production_order_device_config with production_plan_device_config.production_plan_id to production_order_device_config.production_order_id
		deviceConfigs, err := c.deviceConfigRepo.Search(ctx, &repository.SearchProductionOrderDeviceConfigOpts{
			ProductionPlanID: plan.ID,
			Limit:            1000,
			Offset:           0,
		})
		if err != nil {
			return fmt.Errorf("c.deviceConfigRepo.Search: %w", err)
		}
		for _, deviceConfig := range deviceConfigs {
			// TODO use device_config_service.CreateDeviceConfig
			err = c.deviceConfigRepo.Insert(ctx2, &model.ProductionOrderDeviceConfig{
				ID: idutil.ULIDNow(),
				//ProductionOrderID: newProductionOrderID,
				//ProductionPlanID:  deviceConfig.ProductionPlanID,
				DeviceID:     deviceConfig.DeviceID,
				Color:        deviceConfig.Color,
				Description:  deviceConfig.Description,
				Search:       deviceConfig.Search,
				DeviceConfig: deviceConfig.DeviceConfig,
				CreatedBy:    deviceConfig.CreatedBy,
				CreatedAt:    now,
				UpdatedBy:    deviceConfig.UpdatedBy,
				UpdatedAt:    now,
			})
			if err != nil {
				return fmt.Errorf("c.deviceConfigRepo.Insert: %w", err)
			}
		}

		// update production plan with production order id
		plan.ProductionOrderID = cockroach.String(newProductionOrderID)
		if err := c.productionPlanRepo.Update(ctx, plan.ProductionPlan); err != nil {
			return fmt.Errorf("c.productionPlanRepo.Update: %w", err)
		}

		return nil
	})
	if errTx != nil {
		return "", fmt.Errorf("process production order: %w", errTx)
	}
	return productionOrder.ID, nil
}
