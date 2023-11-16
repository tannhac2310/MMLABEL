package production_order

import (
	"context"
	"fmt"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

func (c *productionOrderService) EditProductionOrder(ctx context.Context, opt *EditProductionOrderOpts) error {
	var err error
	// write code to update production order and production order stage in transaction
	err2 := cockroach.ExecInTx(ctx, func(ctx2 context.Context) error {
		err = c.editProductionOrder(ctx2, opt)
		if err != nil {
			return fmt.Errorf("c.editProductionOrder: %w", err)
		}

		err = c.editProductionOrderStage(ctx2, opt)
		if err != nil {
			return fmt.Errorf("c.editProductionOrderStage: %w", err)
		}

		// write code to edit custom field value
		// delete all custom field value with entity type is production order and entity id is production order id
		err := c.customFieldRepo.DeleteByEntity(ctx2, enum.CustomFieldTypeProductionOrder, opt.ID)
		if err != nil {
			return fmt.Errorf("c.customFieldRepo.DeleteByEntity: %w", err)
		}
		// then insert new custom field value
		for field, customFieldValue := range opt.CustomData {
			err = c.customFieldRepo.Insert(ctx2, &model.CustomField{
				ID:         idutil.ULIDNow(),
				EntityType: enum.CustomFieldTypeProductionOrder,
				EntityID:   opt.ID,
				Field:      field,
				Value:      customFieldValue,
			})
			if err != nil {
				return fmt.Errorf("c.customFieldRepo.Insert: %w", err)
			}
		}
		return nil
	})
	if err2 != nil {
		return fmt.Errorf("cockroach.ExecInTx: %w", err2)
	}
	return nil
}

type EditProductionOrderOpts struct {
	ID                    string
	Name                  string
	QtyPaper              int64
	QtyFinished           int64
	QtyDelivered          int64
	PlannedProductionDate time.Time
	Status                enum.ProductionOrderStatus
	DeliveryDate          time.Time
	DeliveryImage         string
	Note                  string
	ProductionOrderStage  []*ProductionOrderStage
	CustomData            map[string]string
}

func (c *productionOrderService) editProductionOrder(ctx context.Context, opt *EditProductionOrderOpts) error {
	table := model.ProductionOrder{}
	updater := cockroach.NewUpdater(table.TableName(), model.ProductionOrderFieldID, opt.ID)

	updater.Set(model.ProductionOrderFieldName, opt.Name)
	updater.Set(model.ProductionOrderFieldQtyPaper, opt.QtyPaper)
	updater.Set(model.ProductionOrderFieldQtyFinished, opt.QtyFinished)
	updater.Set(model.ProductionOrderFieldQtyDelivered, opt.QtyDelivered)
	updater.Set(model.ProductionOrderFieldPlannedProductionDate, opt.PlannedProductionDate)
	updater.Set(model.ProductionOrderFieldStatus, opt.Status)
	updater.Set(model.ProductionOrderFieldDeliveryDate, opt.DeliveryDate)
	updater.Set(model.ProductionOrderFieldDeliveryImage, opt.DeliveryImage)
	updater.Set(model.ProductionOrderFieldNote, opt.Note)

	updater.Set(model.ProductionOrderFieldUpdatedAt, time.Now())

	err := cockroach.UpdateFields(ctx, updater)
	if err != nil {
		return fmt.Errorf("update productionOrder failed %w", err)
	}
	return nil
}

func (c *productionOrderService) editProductionOrderStage(ctx context.Context, opt *EditProductionOrderOpts) error {
	// if production order stage with id empty then create new production order stage
	// if production order stage with id not empty then update production order stage
	// if production order stage with id not in production order stage then delete production order stage
	orderStageIds := make([]string, 0)
	for _, stage := range opt.ProductionOrderStage {
		if stage.ID != "" {
			orderStageIds = append(orderStageIds, stage.ID)
		}
	}
	err := c.deleteProductionOrderStage(ctx, orderStageIds, opt.ID)
	if err != nil {
		return fmt.Errorf("c.deleteProductionOrderStage: %w", err)
	}
	// update production order stage
	for _, orderStage := range opt.ProductionOrderStage {
		if orderStage.ID == "" {
			// create new production order stage
			err = c.productionOrderStageRepo.Insert(ctx, &model.ProductionOrderStage{
				ID:                  idutil.ULIDNow(),
				ProductionOrderID:   opt.ID,
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
				CreatedAt:           time.Now(),
				UpdatedAt:           time.Now(),
			})
			if err != nil {
				return fmt.Errorf("c.productionOrderStageRepo.Insert: %w", err)
			}
			continue
		}
		// find production order stage by id
		orderStageData, err := c.productionOrderStageRepo.FindByID(ctx, orderStage.ID)
		if err != nil {
			return fmt.Errorf("c.productionOrderStageRepo.FindByID: %w", err)
		}

		table := model.ProductionOrderStage{}
		updater := cockroach.NewUpdater(table.TableName(), model.ProductionOrderStageFieldID, orderStage.ID)
		// update date follow status of production order stage
		if orderStageData.Status != orderStage.Status {
			switch orderStage.Status {
			case enum.ProductionOrderStageStatusWaiting:
				updater.Set(model.ProductionOrderStageFieldWaitingAt, time.Now())
			case enum.ProductionOrderStageStatusReception:
				updater.Set(model.ProductionOrderStageFieldReceptionAt, time.Now())
			case enum.ProductionOrderStageStatusProductionStart:
				updater.Set(model.ProductionOrderStageFieldProductionStartAt, time.Now())
			case enum.ProductionOrderStageStatusProductionCompletion:
				updater.Set(model.ProductionOrderStageFieldProductionCompletionAt, time.Now())
			case enum.ProductionOrderStageStatusProductDelivery:
				updater.Set(model.ProductionOrderStageFieldProductDeliveryAt, time.Now())
			}
		}

		updater.Set(model.ProductionOrderStageFieldEstimatedStartAt, orderStage.EstimatedStartAt)
		updater.Set(model.ProductionOrderStageFieldEstimatedCompleteAt, orderStage.EstimatedCompleteAt)
		updater.Set(model.ProductionOrderStageFieldStartedAt, orderStage.StartedAt)
		updater.Set(model.ProductionOrderStageFieldCompletedAt, orderStage.CompletedAt)
		updater.Set(model.ProductionOrderStageFieldStatus, orderStage.Status)
		updater.Set(model.ProductionOrderStageFieldCondition, orderStage.Condition)
		updater.Set(model.ProductionOrderStageFieldNote, orderStage.Note)
		updater.Set(model.ProductionOrderStageFieldData, orderStage.Data)

		updater.Set(model.ProductionOrderStageFieldUpdatedAt, time.Now())

		err = cockroach.UpdateFields(ctx, updater)
		if err != nil {
			return fmt.Errorf("update productionOrderStage failed %w", err)
		}
	}
	return nil
}
