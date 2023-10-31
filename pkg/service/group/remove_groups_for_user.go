package group

import (
	"context"
)

func (g *groupService) RemoveGroupsForUser(ctx context.Context, userID string, groupIDs []string) error {
	return g.userGroupRepo.DeleteByUserIDAndGroupIDs(ctx, userID, groupIDs)
}
