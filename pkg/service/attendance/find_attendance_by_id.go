package attendance

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
)

func (b *attendanceService) FindAttendanceByID(ctx context.Context, id string) (*Data, error) {
	attendances, _, err := b.FindAttendances(ctx, &FindAttendancesOpts{
		IDs: []string{id},
	}, &repository.Sort{
		Order: repository.SortOrderDESC,
		By:    "id",
	}, 1, 0)

	if err != nil {
		return nil, err
	}
	if len(attendances) != 1 {
		return nil, fmt.Errorf("attendance.Search:FindAttendanceByID not found")
	}

	return attendances[0], nil
}
