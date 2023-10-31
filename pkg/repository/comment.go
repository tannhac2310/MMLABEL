package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
)

type CommentRepo interface {
	Insert(ctx context.Context, e *model.Comment) error
	Update(ctx context.Context, e *model.Comment) error
	FindByID(ctx context.Context, id string) (*model.Comment, error)
	SoftDelete(ctx context.Context, id string) error
	Search(ctx context.Context, s *SearchCommentsOpts) ([]*model.Comment, error)
}

type commentRepo struct {
}

func NewCommentRepo() CommentRepo {
	return &commentRepo{}
}

func (r *commentRepo) Insert(ctx context.Context, e *model.Comment) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *commentRepo) FindByID(ctx context.Context, id string) (*model.Comment, error) {
	e := &model.Comment{}
	err := cockroach.FindOne(ctx, e, "id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("cockroach.FindOne: %w", err)
	}

	return e, nil
}

func (r *commentRepo) SoftDelete(ctx context.Context, id string) error {
	sql := `UPDATE comments
		SET deleted_at = NOW()
		WHERE id = $1`

	cmd, err := cockroach.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("cockroach.Exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("not found any records to delete")
	}

	return nil
}
func (r *commentRepo) Update(ctx context.Context, e *model.Comment) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

// all params is options
type SearchCommentsOpts struct {
	IDs       []string
	Title     string
	ParentID  string
	ProgramID string
	StageID   string
	ComboID   string
	CourseID  string
	LessonID  string
	Status    enum.CommonStatus
	Limit     int64
	Offset    int64
}

func (s *SearchCommentsOpts) buildQuery() (string, []interface{}) {
	args := []interface{}{}
	conds := ""
	joins := ""

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND c.%s = ANY($1)", model.CommentFieldID)
	}

	if s.Title != "" {
		args = append(args, "%"+s.Title+"%")
		conds += fmt.Sprintf(" AND c.%s ILIKE $%d", model.CommentFieldTitle, len(args))
	}

	if s.ParentID != "" {
		args = append(args, s.ParentID)
		conds += fmt.Sprintf(" AND u.%s = $%d", model.CommentFieldParentID, len(args))
	}

	if s.ProgramID != "" {
		args = append(args, s.ProgramID)
		conds += fmt.Sprintf(" AND u.%s = $%d", model.CommentFieldProgramID, len(args))
	}
	if s.StageID != "" {
		args = append(args, s.StageID)
		conds += fmt.Sprintf(" AND u.%s = $%d", model.CommentFieldStageID, len(args))
	}
	if s.ComboID != "" {
		args = append(args, s.ComboID)
		conds += fmt.Sprintf(" AND u.%s = $%d", model.CommentFieldComboID, len(args))
	}
	if s.CourseID != "" {
		args = append(args, s.CourseID)
		conds += fmt.Sprintf(" AND u.%s = $%d", model.CommentFieldCourseID, len(args))
	}
	if s.LessonID != "" {
		args = append(args, s.LessonID)
		conds += fmt.Sprintf(" AND u.%s = $%d", model.CommentFieldLessonID, len(args))
	}

	if s.Status > 0 {
		args = append(args, s.Status)
		conds += fmt.Sprintf(" AND u.%s = $%d", model.CommentFieldStatus, len(args))
	}

	c := &model.Comment{}
	fields, _ := c.FieldMap()

	return fmt.Sprintf(`SELECT c.%s
		FROM %s AS c %s
		WHERE TRUE %s AND c.deleted_at IS NULL
		ORDER BY c.id DESC
		LIMIT %d
		OFFSET %d`, strings.Join(fields, ", c."), c.TableName(), joins, conds, s.Limit, s.Offset), args
}

func (r *commentRepo) Search(ctx context.Context, s *SearchCommentsOpts) ([]*model.Comment, error) {
	comments := make([]*model.Comment, 0)
	sql, args := s.buildQuery()
	err := cockroach.Select(ctx, sql, args...).ScanAll(&comments)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return comments, nil
}
