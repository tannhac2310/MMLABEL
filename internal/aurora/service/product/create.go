package product

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
)

func (p productService) CreateProduct(ctx context.Context, opt *CreateProductOpts) (string, error) {
	if opt.Name == "" {
		return "", fmt.Errorf("name is required")
	}
	productID := idutil.ULIDNow()

	now := time.Now()
	errTx := cockroach.ExecInTx(ctx, func(tx context.Context) error {
		// 1. Create product
		err := p.productRepo.Insert(ctx, &model.Product{
			ID:          productID,
			Name:        opt.Name,
			Code:        opt.Code,
			CustomerID:  opt.CustomerID,
			SaleID:      opt.SaleID,
			Description: opt.Description,
			Data:        opt.Data,
			CreatedAt:   now,
			UpdatedAt:   now,
			CreatedBy:   opt.CreatedBy,
			UpdatedBy:   opt.CreatedBy,
		})
		if err != nil {
			return fmt.Errorf("create product failed: %w", err)
		}
		// 2. Create user fields
		for _, uf := range opt.UserField {
			err = p.userFieldRepo.Insert(ctx, &model.CustomField{
				ID:         idutil.ULIDNow(),
				EntityID:   productID,
				EntityType: enum.CustomFieldTypeProduct,
				Field:      uf.Key,
				Value:      uf.Value,
			})
			if err != nil {
				return fmt.Errorf("create user field failed: %w", err)
			}
		}
		return nil
	})

	if errTx != nil {
		return "", fmt.Errorf("create product failed: %w", errTx)
	}

	return productID, nil
}
