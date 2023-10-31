package casbinrule

import (
	"context"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
)

func (b *casbinRulesService) FindCasbinRules(ctx context.Context, limit, offset int64) ([]*repository.CasbinRuleData, *repository.CountResult, error) {
	filter := &repository.SearchCasbinRulesOpts{
		Limit:  limit,
		Offset: offset,
	}
	result, err := b.casbinRulesRepo.Search(ctx, filter)
	if err != nil {
		return nil, nil, err
	}
	total, err := b.casbinRulesRepo.Count(ctx, filter)
	if err != nil {
		return nil, nil, err
	}
	return result, total, nil
}

type FindCasbinRulessOpts struct {
	IDs  []string
	Name string
}
