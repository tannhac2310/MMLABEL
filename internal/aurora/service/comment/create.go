package comment

import "context"

type CreateCommentOpts struct{}

func (c *commentService) CreateComment(ctx context.Context, opt *CreateCommentOpts) (string, error) {
	return "", nil
}
