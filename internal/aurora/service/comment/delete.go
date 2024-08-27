package comment

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

func (c *commentService) DeleteComment(ctx context.Context, id string) error {
	// exec in transaction
	return cockroach.ExecInTx(ctx, func(ctx2 context.Context) error {
		err := c.commentRepo.SoftDelete(ctx2, id)
		if err != nil {
			return fmt.Errorf("c.commentRepo.SoftDelete: %w", err)
		}
		return nil
	})
}
