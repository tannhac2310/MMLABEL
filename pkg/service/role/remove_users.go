package role

import (
	"context"
)

func (r *roleService) RemoveUsers(ctx context.Context, userIDs []string, roleID string) error {
	return r.userRoleRepo.DeleteByUserIDsAndRoleID(ctx, userIDs, roleID)
}
