package stage

import (
	"context"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/syllabus"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
)

type Service interface {
	CreateStage(ctx context.Context, opt *CreateStageOpts) (string, error)
	EditStage(ctx context.Context, opt *EditStageOpts) error
	FindStages(ctx context.Context, opts *FindStagesOpts, sort *repository.Sort, limit, offset int64) ([]*Stage, *repository.CountResult, error)
	SoftDelete(ctx context.Context, id string) error
	FindStageByID(ctx context.Context, id string) (*Stage, error)
}

type stageService struct {
	stageRepo        repository.StageRepo
	courseRepo       repository.CourseRepo
	userRepo         repository.UserRepo
	stageStudentRepo repository.StageStudentRepo
	syllabusService  syllabus.Service
}
type Stage struct {
	*model.Stage
	CountStageStudent int64
	Course            *model.Course
	Teacher           []*model.User
	Assistant         []*model.User
	Syllabus          *model.SyllabusInfo
}

func NewService(
	stageRepo repository.StageRepo,
	courseRepo repository.CourseRepo,
	userRepo repository.UserRepo,
	syllabusService syllabus.Service,
	stageStudentRep repository.StageStudentRepo,
) Service {
	return &stageService{
		stageRepo:        stageRepo,
		courseRepo:       courseRepo,
		userRepo:         userRepo,
		stageStudentRepo: stageStudentRep,
		syllabusService:  syllabusService,
	}
}
