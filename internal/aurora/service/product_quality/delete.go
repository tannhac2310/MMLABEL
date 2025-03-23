package product_quality

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

func (c *productQualityService) Delete(ctx context.Context, id string) error {
	return cockroach.ExecInTx(ctx, func(ctx2 context.Context) error {
		err := c.inspectionFormRepo.SoftDelete(ctx2, id)
		if err != nil {
			return fmt.Errorf("c.inspectionFormRepo.SoftDelete: %w", err)
		}

		err = c.inspectionErrorRepo.SoftDeleteByFormID(ctx2, id)
		if err != nil {
			return fmt.Errorf("c.inspectionErrorRepo.SoftDelete: %w", err)
		}

		return nil
	})
}
