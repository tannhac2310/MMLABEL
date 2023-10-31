package category

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
)

func (b categoryService) CreateCategory(ctx context.Context, opt *CreateCategoryOpts) (string, error) {
	now := time.Now()

	category := &model.Category{
		ID:          idutil.ULIDNow(),
		Name:        opt.Name,
		Description: opt.Description,
		CreatedBy:   opt.CreatedBy,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	err := b.categoryRepo.Insert(ctx, category)
	if err != nil {
		return "", fmt.Errorf("p.categoryRepo.Insert: %w", err)
	}

	return category.ID, nil
}

type CreateCategoryOpts struct {
	Name        string
	Description string
	CreatedBy   string
}
