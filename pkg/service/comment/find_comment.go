package comment

import (
	"context"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
)

func (b *commentService) FindComments(ctx context.Context, opts *FindCommentsOpts, limit, offset int64) ([]*model.Comment, error) {
	return b.commentRepo.Search(ctx, &repository.SearchCommentsOpts{
		IDs:       opts.IDs,
		Title:     opts.Title,
		ParentID:  opts.ParentID,
		ProgramID: opts.ProgramID,
		StageID:   opts.StageID,
		ComboID:   opts.ComboID,
		CourseID:  opts.CourseID,
		Status:    opts.Status,
		LessonID:  opts.LessonID,
		Limit:     limit,
		Offset:    offset,
	})
}

type FindCommentsOpts struct {
	IDs       []string
	Title     string
	ParentID  string
	ProgramID string
	StageID   string
	ComboID   string
	CourseID  string
	LessonID  string
	Status    enum.CommonStatus
}
