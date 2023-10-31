package role

import (
	"context"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
)

func (r *roleService) HighestRole(ctx context.Context, ids []string) (*model.Role, error) {
	return r.roleRepo.HighestRole(ctx, ids)
}
