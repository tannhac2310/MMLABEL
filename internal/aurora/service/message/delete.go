package message

import (
	"context"
)

func (m *messageService) Delete(ctx context.Context, id, userID string) error {
	// todo verify owner message
	return m.EditMessage(ctx, &EditMessageOpts{
		ID:      id,
		UserID: userID,
		Message: "This message has been deleted.",
	})
}
