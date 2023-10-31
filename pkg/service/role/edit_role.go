package role

import (
	"context"
	"fmt"
)

func (r *roleService) EditRole(ctx context.Context, opts *EditRoleOpts) error {
	role, err := r.roleRepo.FindByID(ctx, opts.ID)
	if err != nil {
		return fmt.Errorf("g.roleRepo.FindByID: %w", err)
	}

	role.Name = opts.Name
	role.Priority = opts.Priority
	err = r.roleRepo.Update(ctx, role)
	if err != nil {
		return fmt.Errorf("g.roleRepo.Update: %w", err)
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
