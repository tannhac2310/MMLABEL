package stagestudent

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
)

func (b *stageStudentService) FindStageStudentByID(ctx context.Context, id string) (*Data, error) {
	stageStudents, _, err := b.FindStageStudents(ctx, &FindStageStudentsOpts{
		IDs: []string{id},
	}, &repository.Sort{
		Order: repository.SortOrderDESC,
		By:    "id",
	}, 1, 0)

	if err != nil {
		return nil, err
	}
	if len(stageStudents) != 1 {
		return nil, fmt.Errorf("stageStudent.Search:FindStageStudentByID not found")
	}

	return stageStudents[0], nil
}
