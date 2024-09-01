package comment

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
)

type CreateCommentOpts struct {
	UserID     string
	TargetID   string
	TargetType enum.CommentTarget
	Content    string

	Attachments []*CreateCommentAttachment
}

type CreateCommentAttachment struct {
	FileURL string
}

func (c *commentService) CreateComment(ctx context.Context, opt *CreateCommentOpts) (string, error) {
	user, err := c.userRepo.FindByID(ctx, opt.UserID)
	if err != nil {
		return "", fmt.Errorf("user is not found")
	}

	_, ok := enum.CommentTargetName[opt.TargetType]
	if !ok {
		return "", fmt.Errorf("target is invalid")
	}

	id := idutil.ULIDNow()
	now := time.Now()

	err = cockroach.ExecInTx(ctx, func(ctx2 context.Context) error {
		err := c.commentRepo.Insert(ctx, &model.Comment{
			ID:         id,
			UserID:     user.ID,
			TargetID:   opt.TargetID,
			TargetType: opt.TargetType,
			Content:    opt.Content,
			CreatedAt:  now,
			UpdatedAt:  now,
		})
		if err != nil {
			return fmt.Errorf("c.commentRepo.Insert: %w", err)
		}

		for _, attachment := range opt.Attachments {
			err = c.commentAttachmentRepo.Insert(ctx, &model.CommentAttachment{
				ID:        idutil.ULIDNow(),
				CommentID: id,
				Url:       attachment.FileURL,
				CreatedAt: now,
			})
			if err != nil {
				return fmt.Errorf("c.commentAttachmentRepo.Insert: %w", err)
			}
		}

		return nil
	})

	return id, err
}
