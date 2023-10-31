package category

import (
	"context"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
)

type Service interface {
	CreateCategory(ctx context.Context, opt *CreateCategoryOpts) (string, error)
	EditCategory(ctx context.Context, opt *EditCategoryOpts) error
	FindCategories(ctx context.Context, opts *FindCategoriesOpts, limit, offset int64) ([]*model.Category, error)
	SoftDelete(ctx context.Context, id string) error
	FindCategoryByID(ctx context.Context, id string) (*model.Category, error)
}

type categoryService struct {
	categoryRepo repository.CategoryRepo
}

func NewService(
	categoryRepo repository.CategoryRepo,
) Service {
	return &categoryService{
		categoryRepo: categoryRepo,
	}
}
