package comment

import (
	"context"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
)

type FindCommentsOpts struct {
}

func (c *commentService) FindComments(ctx context.Context, opts FindCommentsOpts, sort *repository.Sort, limit, offset int64) ([]*Data, *repository.CountResult, error) {
	return nil, nil, nil
}
