package production_plan

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
)

type CreateProductionPlanOpts struct {
	Name         string
	ProductName  string
	ProductCode  string
	QtyPaper     int64
	QtyFinished  int64
	QtyDelivered int64
	Thumbnail    string
	Workflow     any
	Note         string
	CustomField  []*CustomField
	CreatedBy    string
}

type CustomField struct {
	Field string
	Value string
}

func (c *productionPlanService) CreateProductionPlan(ctx context.Context, opt *CreateProductionPlanOpts) (string, error) {
	now := time.Now()

	productionPlanID := fmt.Sprintf("in-%d", c.productionPlanRepo.CountRows(ctx)+1)
	//structToArray

	search := fmt.Sprintf("%s %s %s %s", opt.Name, opt.ProductName, opt.ProductCode, opt.Note)
	for _, field := range opt.CustomField {
		search += fmt.Sprintf(" %s", field.Value)
	}
	errTx := cockroach.ExecInTx(ctx, func(ctx2 context.Context) error {
		productionPlan := &model.ProductionPlan{
			ID:            productionPlanID,
			ProductName:   opt.ProductName,
			ProductCode:   opt.ProductCode,
			QtyPaper:      opt.QtyPaper,
			QtyFinished:   opt.QtyFinished,
			QtyDelivered:  opt.QtyDelivered,
			Workflow:      opt.Workflow,
			Thumbnail:     cockroach.String(opt.Thumbnail),
			Status:        enum.ProductionPlanStatusSaleNew,
			Note:          cockroach.String(opt.Note),
			CreatedBy:     opt.CreatedBy,
			CreatedAt:     now,
			UpdatedBy:     opt.CreatedBy,
			UpdatedAt:     now,
			Name:          opt.Name,
			SearchContent: search,
			CurrentStage:  enum.ProductionPlanStageSale,
		}
		if err := c.productionPlanRepo.Insert(ctx2, productionPlan); err != nil {
			return fmt.Errorf("c.productionPlanRepo.Insert: %w", err)
		}

		for _, val := range opt.CustomField {
			// add custom field table
			err := c.customFieldRepo.Insert(ctx2, &model.CustomField{
				ID:         idutil.ULIDNow(),
				EntityType: enum.CustomFieldTypeProductionPlan,
				EntityID:   productionPlanID,
				Field:      val.Field,
				Value:      val.Value,
			})
			if err != nil {
				return fmt.Errorf("c.customFieldRepo.Insert:  %w %s %s %s ", err, productionPlanID, val.Field, val.Value)
			}
		}
		return nil
	})
	if errTx != nil {
		return "", errTx
	}

	return productionPlanID, nil
}
