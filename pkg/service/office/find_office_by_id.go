package office

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
)

func (b *officeService) FindOfficeByID(ctx context.Context, id string) (*repository.OfficeData, error) {
	offices, _, err := b.FindOffices(ctx, &FindOfficesOpts{
		IDs: []string{id},
	}, 1, 0)

	if err != nil {
		return nil, err
	}
	if len(offices) != 1 {
		return nil, fmt.Errorf("offices.Search:FindOfficeByID not found")
	}

	return offices[0], nil
}
