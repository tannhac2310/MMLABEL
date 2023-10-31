package role

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
)

func (r *roleService) CreateRole(ctx context.Context, opts *CreateRoleOpts) error {
	now := time.Now()

	role := &model.Role{
		ID:        opts.ID,
		Name:      opts.Name,
		Priority:  opts.Priority,
		CreatedAt: now,
		UpdatedAt: now,
	}
	err := r.roleRepo.Insert(ctx, role)
	if err != nil {
		return fmt.Errorf("g.roleRepo.Insert: %w", err)
	}

	if len(opts.Permissions) == 0 {
		_, err = r.endforcer.RemovePolicy(opts.ID)
		if err != nil {
			return fmt.Errorf("r.endforcer.RemovePolicy: %w", err)
		}
	}

	// insert casbin rules
	for _, p := range opts.Permissions {
		_, err := r.endforcer.AddPolicy(opts.ID, p)
		if err != nil {
			return fmt.Errorf("g.endforcer.AddPolicy(%s, %s): %w", opts.ID, r, err)
		}
	}

	return nil
}
