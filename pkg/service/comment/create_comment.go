package comment

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
)

func (b commentService) CreateComment(ctx context.Context, opt *CreateCommentOpts) (string, error) {
	now := time.Now()

	comment := &model.Comment{
		ID:        idutil.ULIDNow(),
		Title:     opt.Title,
		Content:   opt.Content,
		Status:    opt.Status,
		ParentID:  opt.ParentID,
		ProgramID: opt.ProgramID,
		StageID:   opt.StageID,
		ComboID:   opt.ComboID,
		CourseID:  opt.CourseID,
		LessonID:  opt.LessonID,
		CreatedAt: now,
		CreatedBy: opt.CreatedBy,
		UpdatedBy: opt.CreatedBy,
		UpdatedAt: now,
	}
	err := b.commentRepo.Insert(ctx, comment)
	if err != nil {
		return "", fmt.Errorf("p.commentRepo.Insert %s $s %w", opt.CreatedBy, err)
	}

	return comment.ID, nil
}

type CreateCommentOpts struct {
	Title     string
	Content   string
	ParentID  string
	ProgramID string
	StageID   string
	ComboID   string
	CourseID  string
	Status    enum.CommonStatus
	LessonID  string
	CreatedBy string
}
