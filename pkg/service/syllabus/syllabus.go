package syllabus

import (
	"context"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
)

type Service interface {
	CreateSyllabus(ctx context.Context, opt *CreateSyllabusOpts) (string, error)
	EditSyllabus(ctx context.Context, opt *EditSyllabusOpts) error
	FindSyllabuses(ctx context.Context, opts *FindSyllabusesOpts, sort *repository.Sort, limit, offset int64) ([]*Data, *repository.CountResult, error)
	SoftDelete(ctx context.Context, id string) error
	FindSyllabusByID(ctx context.Context, id string) (*Data, error)
}

type syllabusService struct {
	syllabusRepo repository.SyllabusRepo
	userRepo     repository.UserRepo
	courseRepo   repository.CourseRepo
}

func NewService(
	syllabusRepo repository.SyllabusRepo,
	userRepo repository.UserRepo,
	courseRepo repository.CourseRepo,
) Service {
	return &syllabusService{
		syllabusRepo: syllabusRepo,
		userRepo:     userRepo,
		courseRepo:   courseRepo,
	}
}

type Data struct {
	*model.Syllabus
	Teacher       *model.UserInfo
	Course        *model.CourseInfo
	CreatedByInfo *model.UserInfo
	UpdatedByInfo *model.UserInfo
}
