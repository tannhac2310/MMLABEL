package category

import (
	"context"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
)

func (b *categoryService) FindCategories(ctx context.Context, opts *FindCategoriesOpts, limit, offset int64) ([]*model.Category, error) {
	return b.categoryRepo.Search(ctx, &repository.SearchCategoriesOpts{
		IDs:         opts.IDs,
		Name:        opts.Name,
		Search:      opts.Search,
		Description: opts.Description,
		Limit:       limit,
		Offset:      offset,
	})
}

type FindCategoriesOpts struct {
	IDs         []string
	Name        string
	Search      string
	Description string
}
