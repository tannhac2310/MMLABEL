package messagerelation

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
)

func (b messageRelationService) CreateMessageRelation(ctx context.Context, opt *UpsertMessageRelationOpts) (string, error) {
	now := time.Now()

	messageRelation := &model.MessageRelation{
		ID:        idutil.ULIDNow(),
		ChatID:    opt.ChatID,
		UserID:    opt.UserID,
		CreatedAt: now,
		UpdatedAt: now,
	}
	err := b.messageRelationRepo.Insert(ctx, messageRelation)
	if err != nil {
		return "", fmt.Errorf("p.messageRelationRepo.Insert: %w", err)
	}

	return messageRelation.ID, nil
}

type UpsertMessageRelationOpts struct {
	UserID string
	ChatID string
}
