package group

import (
	"context"
	"fmt"
)

func (g *groupService) EditGroup(ctx context.Context, opts *EditGroupOpts) error {
	group, err := g.groupRepo.FindByID(ctx, opts.ID)
	if err != nil {
		return fmt.Errorf("g.groupRepo.FindByID: %w", err)
	}

	group.Name = opts.Name
	err = g.groupRepo.Update(ctx, group)
	if err != nil {
		return fmt.Errorf("g.groupRepo.Update: %w", err)
	}

	if len(opts.Roles) == 0 {
		_, err = g.endforcer.RemoveGroupingPolicy(opts.ID)
		if err != nil {
			return fmt.Errorf("g.endforcer.RemoveGroupingPolicy: %w", err)
		}
	}

	// insert casbin rules
	for _, r := range opts.Roles {
		_, err := g.endforcer.AddGroupingPolicy(opts.ID, r)
		if err != nil {
			return fmt.Errorf("g.endforcer.AddGroupingPolicy(%s, %s): %w", opts.ID, r, err)
		}
	}

	return nil
}
