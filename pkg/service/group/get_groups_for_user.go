package group

import (
	"context"
	"fmt"
)

func (g *groupService) GetGroupsForUser(ctx context.Context, userID string) ([]string, error) {
	userGroups, err := g.userGroupRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("u.userGroupRepo.FindByUserID: %w", err)
	}

	groups := make([]string, 0, len(userGroups))
	for _, r := range userGroups {
		groups = append(groups, r.GroupID)
	}

	return groups, nil
}
