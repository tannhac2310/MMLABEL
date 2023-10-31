package group

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
)

type FindGroupsOpts struct {
	IDs  []string
	Name string
}

func (g *groupService) FindGroups(ctx context.Context, opts *FindGroupsOpts, limit, offset int64) ([]*Group, error) {
	groups, err := g.groupRepo.Search(ctx, &repository.SearchGroupOpts{
		IDs:    opts.IDs,
		Name:   opts.Name,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("g.groupRepo.Search: %w", err)
	}

	results := make([]*Group, 0, len(groups))
	for _, group := range groups {
		roles, err := g.endforcer.GetRolesForUser(group.ID)
		if err != nil {
			return nil, err
		}

		results = append(results, &Group{
			Group: group,
			Roles: roles,
		})
	}

	return results, nil
}
