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

type EditCommentOpts struct {
	ID         string
	UserID     string
	TargetID   string
	TargetType enum.CommentTarget
	Content    string

	Attachments []*EditCommentAttachment
}

type EditCommentAttachment struct {
	FileURL string
}

func (c *commentService) EditComment(ctx context.Context, opt *EditCommentOpts) error {
	commentData, err := c.commentRepo.FindByID(ctx, opt.ID)
	if err != nil {
		return fmt.Errorf("comment is not found")
	}
	comment := commentData.Comment
	if comment.UserID != opt.UserID {
		return fmt.Errorf("user is forbidden")
	}
	if comment.TargetID != opt.TargetID || comment.TargetType != int16(opt.TargetType) {
		return fmt.Errorf("comment is invalid")
	}
	if comment.Content == opt.Content {
		return nil
	}

	now := time.Now()

	err = cockroach.ExecInTx(ctx, func(ctx2 context.Context) error {
		err = c.commentHistoryRepo.Insert(ctx, &model.CommentHistory{
			ID:        idutil.ULIDNow(),
			CommentID: comment.ID,
			Content:   comment.Content,
			CreatedAt: now,
		})
		if err != nil {
			return fmt.Errorf("c.commentHistoryRepo.Insert: %w", err)
		}

		comment.Content = opt.Content
		comment.UpdatedAt = now
		err := c.commentRepo.Update(ctx, comment)
		if err != nil {
			return fmt.Errorf("c.commentRepo.Update: %w", err)
		}

		for _, attachment := range opt.Attachments {
			err = c.commentAttachmentRepo.Insert(ctx, &model.CommentAttachment{
				ID:        idutil.ULIDNow(),
				CommentID: comment.ID,
				Url:       attachment.FileURL,
				CreatedAt: now,
			})
			if err != nil {
				return fmt.Errorf("c.commentAttachmentRepo.Insert: %w", err)
			}
		}

		return nil
	})

	return err
}
