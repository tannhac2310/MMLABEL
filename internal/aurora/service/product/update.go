package product

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
)

func (p productService) UpdateProduct(ctx context.Context, opt *UpdateProductOpts) error {
	errTx := cockroach.ExecInTx(ctx, func(tx context.Context) error {
		// 1. Update product
		err := p.productRepo.Update(ctx, &model.Product{
			ID:          opt.ID,
			Name:        opt.Name,
			Code:        opt.Code,
			CustomerID:  opt.CustomerID,
			SaleID:      opt.SaleID,
			Description: opt.Description,
			Data:        opt.Data,
			UpdatedBy:   opt.UpdatedBy,
		})
		if err != nil {
			return fmt.Errorf("update product failed: %w", err)
		}
		// 2. Delete all user fields
		err = p.userFieldRepo.DeleteByEntity(ctx, enum.CustomFieldTypeProduct, opt.ID)
		if err != nil {
			return fmt.Errorf("delete user fields failed: %w", err)
		}
		// 3. Create user fields
		for _, uf := range opt.UserField {
			err = p.userFieldRepo.Insert(ctx, &model.CustomField{
				ID:         idutil.ULIDNow(),
				EntityID:   opt.ID,
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
		return fmt.Errorf("update product failed: %w", errTx)
	}

	return nil
}
