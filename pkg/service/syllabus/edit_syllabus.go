package syllabus

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

func (b *syllabusService) EditSyllabus(ctx context.Context, opt *EditSyllabusOpts) error {
	syllabus, err := b.syllabusRepo.FindByID(ctx, opt.ID)
	if err != nil {
		return fmt.Errorf("b.pondRepo.FindByID: %w", err)
	}

	syllabus.Title = opt.Title
	syllabus.Code = opt.Code
	syllabus.TeacherID = opt.TeacherID
	syllabus.CourseID = opt.CourseID
	syllabus.Description = cockroach.String(opt.Description)
	syllabus.Status = opt.Status
	syllabus.UpdatedBy = cockroach.String(opt.UserID)

	err = b.syllabusRepo.Update(ctx, syllabus)
	if err != nil {
		return fmt.Errorf("p.syllabusRepo.Update: %w", err)
	}

	return nil
}

type EditSyllabusOpts struct {
	ID          string
	Title       string
	CourseID    string
	Code        string
	TeacherID   string
	Description string
	Status      enum.CommonStatus
	UserID      string
}
