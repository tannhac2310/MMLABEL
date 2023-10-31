package group

import (
	"context"
	"fmt"
)

func (g *groupService) FindGroupByID(ctx context.Context, id string) (*Group, error) {
	group, err := g.groupRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("g.groupRepo.FindByID: %w", err)
	}

	roles, err := g.endforcer.GetRolesForUser(group.ID)
	if err != nil {
		return nil, err
	}

	return &Group{
		Group: group,
		Roles: roles,
	}, nil
}
