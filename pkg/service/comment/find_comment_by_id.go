package comment

import (
	"context"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
)

func (b *commentService) FindCommentByID(ctx context.Context, id string) (*model.Comment, error) {
	return b.commentRepo.FindByID(ctx, id)
}
