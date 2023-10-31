package attendance

import (
	"context"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
)

func (b *attendanceService) FindAttendances(ctx context.Context, opts *FindAttendancesOpts, sort *repository.Sort, limit, offset int64) ([]*Data, *repository.CountResult, error) {
	filter := &repository.SearchAttendancesOpts{
		IDs:      opts.IDs,
		CourseID: opts.CourseID,
		StageID:  opts.StageID,
		LessonID: opts.LessonID,
		Limit:    limit,
		Offset:   offset,
		Sort:     sort,
	}
	attendances, err := b.attendanceRepo.Search(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	total, err := b.attendanceRepo.Count(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	results := make([]*Data, 0, len(attendances))
	for _, attendance := range attendances {
		sInfo, err := b.convertData(ctx, attendance)
		if err != nil {
			return nil, nil, err
		}
		results = append(results, sInfo)
	}
	return results, total, nil
}

type FindAttendancesOpts struct {
	IDs      []string
	CourseID string
	StageID  string
	LessonID string
	Status   enum.CommonStatus
}

func (b *attendanceService) convertData(ctx context.Context, attendance *repository.AttendanceData) (*Data, error) {
	sInfo := &Data{
		AttendanceData: attendance,
	}
	// find detail
	detail, err := b.attendanceDetailRepo.Search(ctx, &repository.SearchAttendanceDetailsOpts{
		AttendanceID: attendance.ID,
		Limit:        1000,
		Offset:       0,
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
