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

func (p productService) UpdateProduct(ctx context.Context, opt *UpdateProductOpts) error {
	errTx := cockroach.ExecInTx(ctx, func(tx context.Context) error {
		// 1. Update product
		table := model.Product{}
		updater := cockroach.NewUpdater(table.TableName(), model.ProductFieldID, opt.ID)
		updater.Set(model.ProductFieldName, opt.Name)
		updater.Set(model.ProductFieldCode, opt.Code)
		updater.Set(model.ProductFieldCustomerID, opt.CustomerID)
		updater.Set(model.ProductFieldSaleID, opt.SaleID)
		updater.Set(model.ProductFieldDescription, opt.Description)
		updater.Set(model.ProductFieldData, opt.Data)
		updater.Set(model.ProductFieldUpdatedBy, opt.UpdatedBy)
		updater.Set(model.ProductFieldUpdatedAt, time.Now())

		if err := cockroach.UpdateFields(ctx, updater); err != nil {
			return fmt.Errorf("update product failed: %w", err)
		}

		// 2. Delete all user fields
		err := p.userFieldRepo.DeleteByEntity(ctx, enum.CustomFieldTypeProduct, opt.ID)
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
