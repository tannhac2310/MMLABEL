package comment

import (
	"context"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
)

type Service interface {
	CreateComment(ctx context.Context, opt *CreateCommentOpts) (string, error)
	EditComment(ctx context.Context, opt *EditCommentOpts) error
	FindComments(ctx context.Context, opts *FindCommentsOpts, limit, offset int64) ([]*model.Comment, error)
	SoftDelete(ctx context.Context, id string) error
	FindCommentByID(ctx context.Context, id string) (*model.Comment, error)
}

type commentService struct {
	commentRepo repository.CommentRepo
}

func NewService(
	commentRepo repository.CommentRepo,
) Service {
	return &commentService{
		commentRepo: commentRepo,
	}
}
