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
	Name        string
	CustomerID  string
	SalesID     string
	Thumbnail   string
	Status      enum.ProductionPlanStatus
	Note        string
	CustomField []*CustomField
	CreatedBy   string
}

type CustomField struct {
	Field string
	Value string
}

func (c *productionPlanService) CreateProductionPlan(ctx context.Context, opt *CreateProductionPlanOpts) (string, error) {
	now := time.Now()
	id := idutil.ULIDNow()
	productionPlan := &model.ProductionPlan{
		ID:         id,
		CustomerID: opt.CustomerID,
		SalesID:    opt.SalesID,
		Thumbnail:  cockroach.String(opt.Thumbnail),
		Status:     opt.Status,
		Note:       cockroach.String(opt.Note),
		CreatedBy:  opt.CreatedBy,
		CreatedAt:  now,
		UpdatedAt:  now,
		Name:       opt.Name,
	}

	errTx := cockroach.ExecInTx(ctx, func(ctx2 context.Context) error {
		if err := c.productionPlanRepo.Insert(ctx2, productionPlan); err != nil {
			return fmt.Errorf("c.productionPlanRepo.Insert: %w", err)
		}

		for _, val := range opt.CustomField {
			// add custom field table
			err := c.customFieldRepo.Insert(ctx2, &model.CustomField{
				ID:         idutil.ULIDNow(),
				EntityType: enum.CustomFieldTypeProductionPlan,
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

	return productionPlan.ID, nil
}
