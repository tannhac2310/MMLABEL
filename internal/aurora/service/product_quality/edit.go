package product_quality

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

func (c *productQualityService) EditProductQuality(ctx context.Context, opt *EditProductQualityOpts) error {
	var err error
	table := model.ProductQuality{}
	updater := cockroach.NewUpdater(table.TableName(), model.ProductQualityFieldID, opt.ID)

	updater.Set(model.ProductQualityFieldDefectType, opt.DefectType)
	updater.Set(model.ProductQualityFieldDefectCode, opt.DefectCode)
	updater.Set(model.ProductQualityFieldDefectLevel, opt.DefectLevel)
	updater.Set(model.ProductQualityFieldProductionStageID, opt.ProductionStageID)
	updater.Set(model.ProductQualityFieldDefectiveQuantity, opt.DefectiveQuantity)
	updater.Set(model.ProductQualityFieldGoodQuantity, opt.GoodQuantity)
	updater.Set(model.ProductQualityFieldDescription, opt.Description)

	updater.Set(model.ProductQualityFieldUpdatedAt, time.Now())

	err = cockroach.UpdateFields(ctx, updater)
	if err != nil {
		return fmt.Errorf("update productQuality failed %w", err)
	}
	return nil
}

type EditProductQualityOpts struct {
	ID                string
	DefectType        string
	DefectCode        string
	DeviceIDs         []string
	DefectLevel       string
	ProductionStageID string
	DefectiveQuantity int64
	GoodQuantity      int64
	Description       string
}
