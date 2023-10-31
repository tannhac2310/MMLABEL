package role

import (
	"context"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
)

type FindRoleUsersOpts struct {
	RoleIDs []string
	Search  string
}

func (r *roleService) GetUsersForRole(ctx context.Context, opts *FindRoleUsersOpts, limit, offset int64) ([]*repository.RuleUsersData, *repository.CountResult, error) {
	filter := &repository.SearchUserRoleOpts{
		RoleIDs: opts.RoleIDs,
		Search:  opts.Search,
		Limit:   limit,
		Offset:  offset,
	}
	userRole, err := r.userRoleRepo.Search(ctx, filter)

	if err != nil {
		return nil, nil, err
	}
	total, err := r.userRoleRepo.Count(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	return userRole, total, nil
}
