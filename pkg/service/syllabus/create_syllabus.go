package syllabus

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
)

func (b syllabusService) CreateSyllabus(ctx context.Context, opt *CreateSyllabusOpts) (string, error) {
	now := time.Now()

	syllabus := &model.Syllabus{
		ID:          idutil.ULIDNow(),
		Title:       opt.Title,
		Code:        opt.Code,
		TeacherID:   opt.TeacherID,
		CourseID:    opt.CourseID,
		Description: cockroach.String(opt.Description),
		Status:      opt.Status,
		CreatedBy:   opt.CreatedBy,
		UpdatedBy:   cockroach.String(opt.CreatedBy),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	err := b.syllabusRepo.Insert(ctx, syllabus)
	if err != nil {
		return "", fmt.Errorf("p.syllabusRepo.Insert: %w", err)
	}

	return syllabus.ID, nil
}

type CreateSyllabusOpts struct {
	Title       string
	Code        string
	CourseID    string
	TeacherID   string
	Description string
	Status      enum.CommonStatus
	CreatedBy   string
}
