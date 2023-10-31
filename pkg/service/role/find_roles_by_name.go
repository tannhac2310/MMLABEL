package role

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
)

type FindRolesOpts struct {
	IDs  []string
	Name string
}

func (r *roleService) FindRoles(ctx context.Context, opts *FindRolesOpts, limit, offset int64) ([]*Role, error) {
	roles, err := r.roleRepo.Search(ctx, &repository.SearchRoleOpts{
		IDs:    opts.IDs,
		Name:   opts.Name,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("r.roleRepo.Search: %w", err)
	}

	roleWithPers := make([]*Role, 0, len(roles))
	for _, role := range roles {
		userCount, err := r.userRoleRepo.Count(ctx, &repository.SearchUserRoleOpts{
			RoleIDs: []string{"-1", role.ID},
		})
		if err != nil {
			return nil, err
		}
		roleWithPers = append(roleWithPers, &Role{
			Role:        role,
			Permissions: r.getPermissionForRole(role.ID),
			UserCount:   userCount.Count,
		})
	}

	return roleWithPers, nil
}
