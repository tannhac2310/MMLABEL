package messagerelation

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

func (b *messageRelationService) EditMessageRelation(ctx context.Context, opt *EditMessageRelationOpts) error {
	var err error
	table := model.MessageRelation{}
	updater := cockroach.NewUpdater(table.TableName(), model.MessageRelationFieldID, opt.ID)

	updater.Set(model.MessageRelationFieldUserID, opt.UserID)

	err = cockroach.UpdateFields(ctx, updater)
	if err != nil {
		return fmt.Errorf("err update messageRelation status: %w", err)
	}

	return nil
}

type EditMessageRelationOpts struct {
	ID     string
	UserID string
}
