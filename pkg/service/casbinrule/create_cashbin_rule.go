package casbinrule

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
)

func (b casbinRulesService) CreateCasbinRule(ctx context.Context, opt *CreateCasbinRulesOpts) (string, error) {
	casbinRules := &model.CasbinRule{
		PType: "p",
		V0:    opt.V0,
		V1:    opt.V1,
		V2:    "",
		V3:    "",
		V4:    "",
		V5:    "",
		Rowid: time.Now().Unix(),
	}
	err := b.casbinRulesRepo.Insert(ctx, casbinRules)
	if err != nil {
		return "", fmt.Errorf("p.casbinRulesRepo.Insert: %w v0: %s, v1: %s ", err, opt.V0, opt.V1)
	}

	return "ok", nil
}

type CreateCasbinRulesOpts struct {
	V0 string
	V1 string
}
