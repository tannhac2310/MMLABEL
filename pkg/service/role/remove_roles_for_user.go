package role

import (
	"context"
)

func (r *roleService) RemoveRolesForUser(ctx context.Context, userID string, roleIDs []string) error {
	return r.userRoleRepo.DeleteByUserIDAndRoleIDs(ctx, userID, roleIDs)
}
