package product

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

func (p productService) DeleteProduct(ctx context.Context, id string) error {

	errTx := cockroach.ExecInTx(ctx, func(tx context.Context) error {
		// 1. Soft delete product
		err := p.productRepo.SoftDelete(ctx, id)
		if err != nil {
			return fmt.Errorf("delete product failed: %w", err)
		}
		// 2. Delete all user fields
		err = p.userFieldRepo.DeleteByEntity(ctx, enum.CustomFieldTypeProduct, id)
		if err != nil {
			return fmt.Errorf("delete user fields failed: %w", err)
		}
		return nil
	})
	if errTx != nil {
		return fmt.Errorf("delete product failed: %w", errTx)
	}

	return nil
}
