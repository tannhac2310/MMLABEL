package product_quality

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
)

func (c *productQualityService) CreateProductQuality(ctx context.Context, opt *CreateProductQualityOpts) (string, error) {
	now := time.Now()

	productQuality := &model.ProductQuality{
		ID:                idutil.ULIDNow(),
		ProductionOrderID: cockroach.String(opt.ProductionOrderID),
		ProductID:         cockroach.String(opt.ProductID),
		DeviceIDs:         opt.DeviceIDs,
		DefectType:        cockroach.String(opt.DefectType),
		DefectCode:        cockroach.String(opt.DefectCode),
		DefectLevel:       opt.DefectLevel,
		ProductionStageID: cockroach.String(opt.ProductionStageID),
		DefectiveQuantity: opt.DefectiveQuantity,
		GoodQuantity:      opt.GoodQuantity,
		Description:       cockroach.String(opt.Description),
		CreatedBy:         opt.CreatedBy,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	errTx := cockroach.ExecInTx(ctx, func(ctx2 context.Context) error {
		err := c.productQualityRepo.Insert(ctx2, productQuality)
		if err != nil {
			return fmt.Errorf("c.productQualityRepo.Insert: %w", err)
		}

		return nil
	})
	if errTx != nil {
		return "", errTx
	}
	return productQuality.ID, nil
}

type CreateProductQualityOpts struct {
	ProductionOrderID string
	ProductID         string
	DeviceIDs         []string
	DefectType        string
	DefectCode        string
	DefectLevel       int16
	ProductionStageID string
	DefectiveQuantity int64
	GoodQuantity      int64
	Description       string
	CreatedBy         string
}
