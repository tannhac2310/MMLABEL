package casbinrule

import (
	"context"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
)

type Service interface {
	CreateCasbinRule(ctx context.Context, opt *CreateCasbinRulesOpts) (string, error)
	FindCasbinRules(ctx context.Context, limit, offset int64) ([]*repository.CasbinRuleData, *repository.CountResult, error)
	Delete(ctx context.Context, v0, v1 string) error
}

type casbinRulesService struct {
	casbinRulesRepo repository.CasbinRuleRepo
}

func NewService(
	casbinRulesRepo repository.CasbinRuleRepo,
) Service {
	return &casbinRulesService{
		casbinRulesRepo: casbinRulesRepo,
	}
}
