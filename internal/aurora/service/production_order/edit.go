package production_order

import (
	"context"
	"fmt"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

func (c *productionOrderService) EditProductionOrder(ctx context.Context, opt *EditProductionOrderOpts) error {
	var err error
	table := model.ProductionOrder{}
	updater := cockroach.NewUpdater(table.TableName(), model.ProductionOrderFieldID, opt.ID)

	updater.Set(model.ProductionOrderFieldQtyPaper, opt.QtyPaper)
	updater.Set(model.ProductionOrderFieldQtyFinished, opt.QtyFinished)
	updater.Set(model.ProductionOrderFieldQtyDelivered, opt.QtyDelivered)
	updater.Set(model.ProductionOrderFieldPlannedProductionDate, opt.PlannedProductionDate)
	updater.Set(model.ProductionOrderFieldStatus, opt.Status)
	updater.Set(model.ProductionOrderFieldDeliveryDate, opt.DeliveryDate)
	updater.Set(model.ProductionOrderFieldDeliveryImage, opt.DeliveryImage)
	updater.Set(model.ProductionOrderFieldNote, opt.Note)

	updater.Set(model.ProductionOrderFieldUpdatedAt, time.Now())

	err = cockroach.UpdateFields(ctx, updater)
	if err != nil {
		return fmt.Errorf("update productionOrder failed %w", err)
	}
	return nil
}

type EditProductionOrderOpts struct {
	ID                    string
	QtyPaper              int64
	QtyFinished           int64
	QtyDelivered          int64
	PlannedProductionDate time.Time
	Status                enum.ProductionOrderStatus
	DeliveryDate          time.Time
	DeliveryImage         string
	Note                  string
}
