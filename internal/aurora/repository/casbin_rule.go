package repository

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type CasbinRuleRepo interface {
	Insert(ctx context.Context, e *model.CasbinRule) error
	Update(ctx context.Context, e *model.CasbinRule) error
	SoftDelete(ctx context.Context, id string) error
	Search(ctx context.Context, s *SearchCasbinRuleOpts) ([]*CasbinRuleData, error)
	Count(ctx context.Context, s *SearchCasbinRuleOpts) (*CountResult, error)
}

type sCasbinRuleRepo struct {
}

func NewCasbinRuleRepo() CasbinRuleRepo {
	return &sCasbinRuleRepo{}
}

func (r *sCasbinRuleRepo) Insert(ctx context.Context, e *model.CasbinRule) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}
