package category

import (
	"context"
	"fmt"
)

func (b *categoryService) EditCategory(ctx context.Context, opt *EditCategoryOpts) error {
	category, err := b.categoryRepo.FindByID(ctx, opt.ID)
	if err != nil {
		return fmt.Errorf("f.pondRepo.FindByID: %w", err)
	}

	category.Name = opt.Name
	category.Description = opt.Description
	category.UpdatedBy = opt.UpdatedBy

	err = b.categoryRepo.Update(ctx, category)
	if err != nil {
		return fmt.Errorf("p.pondRepo.Update: %w", err)
	}

	return nil
}

type EditCategoryOpts struct {
	ID          string
	Name        string
	Description string
	UpdatedBy   string
}
