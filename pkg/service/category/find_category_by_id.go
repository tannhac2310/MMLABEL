package category

import (
	"context"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
)

func (b *categoryService) FindCategoryByID(ctx context.Context, id string) (*model.Category, error) {
	return b.categoryRepo.FindByID(ctx, id)
}
