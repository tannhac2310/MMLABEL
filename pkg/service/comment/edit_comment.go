package comment

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

func (b *commentService) EditComment(ctx context.Context, opt *EditCommentOpts) error {
	comment, err := b.commentRepo.FindByID(ctx, opt.ID)
	if err != nil {
		return fmt.Errorf("b.pondRepo.FindByID: %w", err)
	}

	comment.Title = opt.Title
	comment.Content = opt.Content
	comment.Status = opt.Status

	err = b.commentRepo.Update(ctx, comment)
	if err != nil {
		return fmt.Errorf("p.pondRepo.Update: %w", err)
	}

	return nil
}

type EditCommentOpts struct {
	ID      string
	Title   string
	Content string
	Status  enum.CommonStatus
}
