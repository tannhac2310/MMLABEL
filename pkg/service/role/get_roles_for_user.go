package role

import (
	"context"
	"fmt"
)

func (r *roleService) GetRolesForUser(ctx context.Context, userID string) ([]string, error) {
	userRoles, err := r.userRoleRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("u.userRoleRepo.FindByUserID: %w", err)
	}

	roles := make([]string, 0, len(userRoles))
	for _, r := range userRoles {
		roles = append(roles, r.RoleID)
	}

	return roles, nil
}
