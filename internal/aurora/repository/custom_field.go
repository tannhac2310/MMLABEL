package repository

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type CustomFieldRepo interface {
	Insert(ctx context.Context, e *model.CustomField) error
	Update(ctx context.Context, e *model.CustomField) error
	SoftDelete(ctx context.Context, id string) error
	Search(ctx context.Context, s *SearchCustomFieldOpts) ([]*CustomFieldData, error)
	Count(ctx context.Context, s *SearchCustomFieldOpts) (*CountResult, error)
}

type sCustomFieldRepo struct {
}

func NewCustomFieldRepo() CustomFieldRepo {
	return &sCustomFieldRepo{}
}

func (r *sCustomFieldRepo) Insert(ctx context.Context, e *model.CustomField) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}
