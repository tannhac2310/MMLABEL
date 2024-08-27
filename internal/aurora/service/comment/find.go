package comment

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

type FindCommentsOpts struct {
	TargetID   string
	TargetType enum.CommentTarget
}

func (c *commentService) FindComments(ctx context.Context, opts FindCommentsOpts, sort *repository.Sort, limit, offset int64) ([]*Data, *repository.CountResult, error) {
	filter := &repository.SearchCommentOpts{
		TargetID:   opts.TargetID,
		TargetType: opts.TargetType,
		Limit:      limit,
		Offset:     offset,
		Sort: &repository.Sort{
			Order: repository.SortOrderDESC,
			By:    fmt.Sprintf("b.%s", model.CommentFieldCreatedAt),
		},
	}

	total, err := c.commentRepo.Count(ctx, filter)
	if err != nil {
		return nil, nil, err
	}
	if total.Count == 0 {
		return make([]*Data, 0), total, nil
	}

	comments, err := c.commentRepo.Search(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	results := make([]*Data, 0, len(comments))
	for _, comment := range comments {
		results = append(results, &Data{comment})
	}
	return results, total, nil
}

func (c *commentService) FindCommentHistories(ctx context.Context, commentId string) ([]*HistoryData, *repository.CountResult, error) {
	filter := &repository.SearchCommentHistoryOpts{
		CommentID: commentId,
		Limit:     10,
		Offset:    0,
		Sort: &repository.Sort{
			Order: repository.SortOrderDESC,
			By:    fmt.Sprintf("b.%s", model.CommentHistoryFieldCreatedAt),
		},
	}

	total, err := c.commentHistoryRepo.Count(ctx, filter)
	if err != nil {
		return nil, nil, err
	}
	if total.Count == 0 {
		return make([]*HistoryData, 0), total, nil
	}

	commentHistories, err := c.commentHistoryRepo.Search(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	results := make([]*HistoryData, 0, len(commentHistories))
	for _, commentHistory := range commentHistories {
		results = append(results, &HistoryData{commentHistory})
	}
	return results, total, nil
}
