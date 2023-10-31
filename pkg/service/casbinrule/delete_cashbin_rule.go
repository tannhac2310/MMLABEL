package casbinrule

import (
	"context"
)

func (b *casbinRulesService) Delete(ctx context.Context, v0, v1 string) error {
	return b.casbinRulesRepo.Delete(ctx, v0, v1)
}
