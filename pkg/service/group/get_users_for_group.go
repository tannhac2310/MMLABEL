package group

import (
	"context"
	"fmt"
)

func (g *groupService) GetUsersForGroup(ctx context.Context, groupID string) ([]string, error) {
	userGroups, err := g.userGroupRepo.FindByGroupID(ctx, groupID)
	if err != nil {
		return nil, fmt.Errorf("u.userGroupRepo.FindByGroupID: %w", err)
	}

	userIDs := make([]string, 0, len(userGroups))
	for _, r := range userGroups {
		userIDs = append(userIDs, r.UserID)
	}

	return userIDs, nil
}
