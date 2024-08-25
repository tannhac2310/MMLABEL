package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type CommentRepo interface {
	Insert(ctx context.Context, e *model.Comment) error
	Update(ctx context.Context, e *model.Comment) error
	SoftDelete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*CommentData, error)
	Search(ctx context.Context, s *SearchCommentOpts) ([]*CommentData, error)
	Count(ctx context.Context, s *SearchCommentOpts) (*CountResult, error)
}

type sCommentRepo struct {
}

func NewCommentRepo() CommentRepo {
	return &sCommentRepo{}
}

func (r *sCommentRepo) Insert(ctx context.Context, e *model.Comment) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *sCommentRepo) FindByID(ctx context.Context, id string) (*CommentData, error) {
	e := &CommentData{}
	err := cockroach.FindOne(ctx, e, "id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("sCommentRepo.cockroach.FindOne: %w", err)
	}

	return e, nil
}
func (r *sCommentRepo) Update(ctx context.Context, e *model.Comment) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

func (r *sCommentRepo) SoftDelete(ctx context.Context, id string) error {
	sql := "UPDATE comments SET deleted_at = NOW() WHERE id = $1;"

	cmd, err := cockroach.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("comments cockroach.Exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("*sCommentRepo not found any records to delete")
	}

	return nil
}

// SearchCommentOpts all params is options
type SearchCommentOpts struct {
	IDs      []string
	TargetID string
	Limit    int64
	Offset   int64
	Sort     *Sort
}

func (s *SearchCommentOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := " LEFT JOIN users AS cu ON b.user_id = cu.id "

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.CommentFieldID)
	}
	if s.TargetID != "" {
		args = append(args, s.TargetID)
		conds += fmt.Sprintf(" AND b.%[2]s = $%[1]d", len(args), model.CommentFieldTargetID)
	}

	b := &model.Comment{}
	fields, _ := b.FieldMap()
	if isCount {
		return fmt.Sprintf("SELECT count(*) as cnt FROM %s AS b %s WHERE TRUE %s AND b.deleted_at IS NULL", b.TableName(), joins, conds), args
	}

	order := " ORDER BY b.id DESC "
	if s.Sort != nil {
		order = fmt.Sprintf(" ORDER BY b.%s %s", s.Sort.By, s.Sort.Order)
	}
	return fmt.Sprintf("SELECT b.%s, cu.name as user_name FROM %s AS b %s WHERE TRUE %s AND b.deleted_at IS NULL %s LIMIT %d OFFSET %d", strings.Join(fields, ", b."), b.TableName(), joins, conds, order, s.Limit, s.Offset), args
}

type CommentData struct {
	*model.Comment
	UserName string `db:"user_name"`
}

func (r *sCommentRepo) Search(ctx context.Context, s *SearchCommentOpts) ([]*CommentData, error) {
	Comment := make([]*CommentData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&Comment)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return Comment, nil
}

func (r *sCommentRepo) Count(ctx context.Context, s *SearchCommentOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("sCommentRepo.Count: %w", err)
	}

	return countResult, nil
}
