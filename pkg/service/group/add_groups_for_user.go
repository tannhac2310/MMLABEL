package group

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
)

func (g *groupService) AddGroupsForUser(ctx context.Context, userID string, groupIDs []string) error {
	// TODO: check groupIDs existed in table groups
	// check duplicated group

	now := time.Now()
	for _, groupID := range groupIDs {
		userGroup := &model.UserGroup{
			ID:        idutil.ULIDNow(),
			UserID:    userID,
			GroupID:   groupID,
			CreatedAt: now,
			UpdatedAt: now,
		}

		if err := g.userGroupRepo.Insert(ctx, userGroup); err != nil {
			return fmt.Errorf("r.userGroupRepo.Insert: %w", err)
		}
	}

	return nil
}
