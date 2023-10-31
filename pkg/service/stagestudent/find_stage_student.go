package stagestudent

import (
	"context"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
)

func (b *stageStudentService) FindStageStudents(ctx context.Context, opts *FindStageStudentsOpts, sort *repository.Sort, limit, offset int64) ([]*Data, *repository.CountResult, error) {
	filter := &repository.SearchStageStudentsOpts{
		IDs:       opts.IDs,
		CourseID:  opts.CourseID,
		StageID:   opts.StageID,
		StudentID: opts.StudentID,
		Status:    opts.Status,
		LessonID:  opts.LessonID,
		Limit:     limit,
		Offset:    offset,
		Sort:      sort,
	}
	stageStudents, err := b.stageStudentRepo.Search(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	total, err := b.stageStudentRepo.Count(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	results := make([]*Data, 0, len(stageStudents))
	for _, stageStudent := range stageStudents {
		sInfo, err := b.convertData(ctx, stageStudent)
		if err != nil {
			return nil, nil, err
		}
		results = append(results, sInfo)
	}
	return results, total, nil
}

type FindStageStudentsOpts struct {
	IDs       []string
	StudentID string
	CourseID  string
	StageID   string
	LessonID  string
	Status    enum.StageStudentStatus
}

func (b *stageStudentService) convertData(ctx context.Context, stageStudent *repository.StageStudentData) (*Data, error) {
	sInfo := &Data{
		StageStudentData: stageStudent,
	}
	// find detail
	detail, err := b.attendanceDetailRepo.Search(ctx, &repository.SearchAttendanceDetailsOpts{
		StudentID: stageStudent.StudentID,
		Limit:     1000,
		Offset:    0,
		Sort: &repository.Sort{
			Order: repository.SortOrderASC,
			By:    "id",
		},
	})
	if err != nil {
		return nil, err
	}
	sInfo.Detail = detail

	return sInfo, nil
}
