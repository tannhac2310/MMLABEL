package repository

import (
	"context"
	"fmt"
	"strings"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type CommentHistoryRepo interface {
	Insert(ctx context.Context, e *model.CommentHistory) error
	FindByID(ctx context.Context, id string) (*CommentHistoryData, error)
	Search(ctx context.Context, s *SearchCommentHistoryOpts) ([]*CommentHistoryData, error)
	Count(ctx context.Context, s *SearchCommentHistoryOpts) (*CountResult, error)
}

type sCommentHistoryRepo struct {
}

func NewCommentHistoryRepo() CommentHistoryRepo {
	return &sCommentHistoryRepo{}
}

func (r *sCommentHistoryRepo) Insert(ctx context.Context, e *model.CommentHistory) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *sCommentHistoryRepo) FindByID(ctx context.Context, id string) (*CommentHistoryData, error) {
	e := &CommentHistoryData{}
	err := cockroach.FindOne(ctx, e, "id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("sCommentHistoryRepo.cockroach.FindOne: %w", err)
	}

	return e, nil
}

// SearchCommentHistoryOpts all params is options
type SearchCommentHistoryOpts struct {
	IDs []string
	// todo add more search options
	Limit  int64
	Offset int64
	Sort   *Sort
}

func (s *SearchCommentHistoryOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.CommentHistoryFieldID)
	}
	// todo add more search options example:
	//if s.Name != "" {
	//	args = append(args, "%"+s.Name+"%")
	//	conds += fmt.Sprintf(" AND (b.%[2]s ILIKE $%[1]d OR b.%[3]s ILIKE $%[1]d)",
	//		len(args), model.CommentHistoryFieldName, model.CommentHistoryFieldCode)
	//}
	//if s.Code != "" {
	//	args = append(args, s.Code)
	//	conds += fmt.Sprintf(" AND b.%s ILIKE $%d", model.CommentHistoryFieldCode, len(args))
	//}

	b := &model.CommentHistory{}
	fields, _ := b.FieldMap()
	if isCount {
		return fmt.Sprintf("SELECT count(*) as cnt FROM %s AS b %s WHERE TRUE %s AND b.deleted_at IS NULL", b.TableName(), joins, conds), args
	}

	order := " ORDER BY b.id DESC "
	if s.Sort != nil {
		order = fmt.Sprintf(" ORDER BY b.%s %s", s.Sort.By, s.Sort.Order)
	}
	return fmt.Sprintf("SELECT b.%s FROM %s AS b %s WHERE TRUE %s AND b.deleted_at IS NULL %s LIMIT %d OFFSET %d", strings.Join(fields, ", b."), b.TableName(), joins, conds, order, s.Limit, s.Offset), args
}

type CommentHistoryData struct {
	*model.CommentHistory
}

func (r *sCommentHistoryRepo) Search(ctx context.Context, s *SearchCommentHistoryOpts) ([]*CommentHistoryData, error) {
	CommentHistory := make([]*CommentHistoryData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&CommentHistory)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return CommentHistory, nil
}

func (r *sCommentHistoryRepo) Count(ctx context.Context, s *SearchCommentHistoryOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("sCommentHistoryRepo.Count: %w", err)
	}

	return countResult, nil
}
