package role

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
)

func (r *roleService) AddRolesForUser(ctx context.Context, userID string, roleIDs []string) error {
	now := time.Now()
	for _, roleID := range roleIDs {
		userRole := &model.UserRole{
			ID:        idutil.ULIDNow(),
			UserID:    userID,
			RoleID:    roleID,
			CreatedAt: now,
			UpdatedAt: now,
		}

		if err := r.userRoleRepo.Insert(ctx, userRole); err != nil {
			return fmt.Errorf("r.userRoleRepo.Insert: %w", err)
		}
	}

	return nil
}
