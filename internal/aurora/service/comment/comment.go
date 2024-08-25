package comment

import (
	"context"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
)

type Service interface {
	FindComments(ctx context.Context, opts FindCommentsOpts, sort *repository.Sort, limit, offset int64) ([]*Data, *repository.CountResult, error)
	CreateComment(ctx context.Context, opt *CreateCommentOpts) (string, error)
	EditComment(ctx context.Context, opt *EditCommentOpts) error
	DeleteComment(ctx context.Context, id string) error
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

type Data struct {
	*repository.CommentData
}
