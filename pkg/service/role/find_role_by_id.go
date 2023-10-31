package role

import (
	"context"
	"fmt"
)

func (r *roleService) FindRoleByID(ctx context.Context, id string) (*Role, error) {
	role, err := r.roleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("r.roleRepo.FindByID: %w", err)
	}

	return &Role{
		Role:        role,
		Permissions: r.getPermissionForRole(role.ID),
	}, nil
}

func (r *roleService) getPermissionForRole(roleID string) []string {
	persWithName := r.endforcer.GetPermissionsForUser(roleID)
	pers := make([]string, 0, len(persWithName))
	for _, p := range persWithName {
		if len(p) < 2 {
			continue
		}

		pers = append(pers, p[1])
	}

	return pers
}
