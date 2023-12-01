package production_order

import (
	"context"
	"fmt"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
)

func (c *productionOrderService) CreateProductionOrder(ctx context.Context, opt *CreateProductionOrderOpts) (string, error) {
	now := time.Now()
	id := idutil.ULIDNow()
	productionOrder := &model.ProductionOrder{
		ID:                  id,
		ProductCode:         opt.ProductCode,
		ProductName:         opt.ProductName,
		CustomerID:          opt.CustomerID,
		SalesID:             opt.SalesID,
		QtyPaper:            opt.QtyPaper,
		QtyFinished:         opt.QtyFinished,
		QtyDelivered:        opt.QtyDelivered,
		DeliveryDate:        opt.DeliveryDate,
		DeliveryImage:       cockroach.String(opt.DeliveryImage),
		Status:              opt.Status,
		Note:                cockroach.String(opt.Note),
		CreatedBy:           opt.CreatedBy,
		CreatedAt:           now,
		UpdatedAt:           now,
		Name:                opt.Name,
		EstimatedStartAt:    cockroach.Time(opt.EstimatedStartAt),
		EstimatedCompleteAt: cockroach.Time(opt.EstimatedCompleteAt),
	}

	errTx := cockroach.ExecInTx(ctx, func(ctx2 context.Context) error {
		err := c.productionOrderRepo.Insert(ctx2, productionOrder)
		if err != nil {
			return fmt.Errorf("c.productionOrderRepo.Insert: %w", err)
		}
		for _, orderStage := range opt.ProductionOrderStage {
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
		for _, val := range opt.CustomField {
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

type CreateProductionOrderOpts struct {
	Name                 string
	ProductCode          string
	ProductName          string
	CustomerID           string
	SalesID              string
	QtyPaper             int64
	QtyFinished          int64
	QtyDelivered         int64
	EstimatedStartAt     time.Time
	EstimatedCompleteAt  time.Time
	DeliveryDate         time.Time
	DeliveryImage        string
	Status               enum.ProductionOrderStatus
	Note                 string
	ProductionOrderStage []*ProductionOrderStage
	CustomField          []*CustomField
	CreatedBy            string
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

type CustomField struct {
	Field string
	Value string
}
