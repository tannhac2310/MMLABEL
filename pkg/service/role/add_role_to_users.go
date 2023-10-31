package role

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
)

func (r *roleService) AddRoleToUsers(ctx context.Context, userIDs []string, roleID, createdBy string) error {
	now := time.Now()
	for _, userID := range userIDs {
		userRole := &model.UserRole{
			ID:        idutil.ULIDNow(),
			UserID:    userID,
			RoleID:    roleID,
			CreatedBy: cockroach.String(createdBy),
			CreatedAt: now,
			UpdatedAt: now,
		}

		if err := r.userRoleRepo.Insert(ctx, userRole); err != nil {
			return fmt.Errorf("r.userRoleRepo.Insert: %w", err)
		}
	}

	return nil
}
