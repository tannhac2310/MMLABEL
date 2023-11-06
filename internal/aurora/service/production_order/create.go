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

	productionOrder := &model.ProductionOrder{
		ID:                    idutil.ULIDNow(),
		ProductCode:           opt.ProductCode,
		ProductName:           opt.ProductName,
		CustomerID:            opt.CustomerID,
		SalesID:               opt.SalesID,
		QtyPaper:              opt.QtyPaper,
		QtyFinished:           opt.QtyFinished,
		QtyDelivered:          opt.QtyDelivered,
		PlannedProductionDate: opt.PlannedProductionDate,
		DeliveryDate:          opt.DeliveryDate,
		DeliveryImage:         cockroach.String(opt.DeliveryImage),
		Status:                opt.Status,
		Note:                  cockroach.String(opt.Note),
		CreatedBy:             opt.CreatedBy,
		CreatedAt:             now,
		UpdatedAt:             now,
	}

	errTx := cockroach.ExecInTx(ctx, func(ctx2 context.Context) error {
		err := c.productionOrderRepo.Insert(ctx2, productionOrder)
		if err != nil {
			return fmt.Errorf("c.productionOrderRepo.Insert: %w", err)
		}

		return nil
	})
	if errTx != nil {
		return "", errTx
	}
	return productionOrder.ID, nil
}

type CreateProductionOrderOpts struct {
	ProductCode           string
	ProductName           string
	CustomerID            string
	SalesID               string
	QtyPaper              int64
	QtyFinished           int64
	QtyDelivered          int64
	PlannedProductionDate time.Time
	DeliveryDate          time.Time
	DeliveryImage         string
	Status                enum.ProductionOrderStatus
	Note                  string
	CreatedBy             string
}
