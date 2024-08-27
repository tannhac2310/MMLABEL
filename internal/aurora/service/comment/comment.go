package comment

import (
	"context"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	repository2 "mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
)

type Service interface {
	FindComments(ctx context.Context, opts FindCommentsOpts, sort *repository.Sort, limit, offset int64) ([]*Data, *repository.CountResult, error)
	FindCommentHistories(ctx context.Context, commentId string) ([]*HistoryData, *repository.CountResult, error)
	CreateComment(ctx context.Context, opt *CreateCommentOpts) (string, error)
	EditComment(ctx context.Context, opt *EditCommentOpts) error
	DeleteComment(ctx context.Context, id string) error
}

type commentService struct {
	commentRepo           repository.CommentRepo
	commentHistoryRepo    repository.CommentHistoryRepo
	commentAttachmentRepo repository.CommentAttachmentRepo
	userRepo              repository2.UserRepo
}

func NewService(
	commentRepo repository.CommentRepo,
	commentHistoryRepo repository.CommentHistoryRepo,
	commentAttachmentRepo repository.CommentAttachmentRepo,
	userRepo repository2.UserRepo,
) Service {
	return &commentService{
		commentRepo:           commentRepo,
		commentHistoryRepo:    commentHistoryRepo,
		commentAttachmentRepo: commentAttachmentRepo,
		userRepo:              userRepo,
	}
}

type Data struct {
	*repository.CommentData
}

type HistoryData struct {
	*repository.CommentHistoryData
}
