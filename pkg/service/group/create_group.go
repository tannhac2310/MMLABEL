package group

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
)

func (g *groupService) CreateGroup(ctx context.Context, opts *CreateGroupOpts) error {
	now := time.Now()

	group := &model.Group{
		ID:        opts.ID,
		Name:      opts.Name,
		CreatedAt: now,
		UpdatedAt: now,
	}
	err := g.groupRepo.Insert(ctx, group)
	if err != nil {
		return fmt.Errorf("g.groupRepo.Insert: %w", err)
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
