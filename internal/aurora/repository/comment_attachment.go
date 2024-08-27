package repository

import (
	"context"
	"fmt"
	"strings"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type CommentAttachmentRepo interface {
	Insert(ctx context.Context, e *model.CommentAttachment) error
	SoftDelete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*CommentAttachmentData, error)
	Search(ctx context.Context, s *SearchCommentAttachmentOpts) ([]*CommentAttachmentData, error)
	Count(ctx context.Context, s *SearchCommentAttachmentOpts) (*CountResult, error)
}

type sCommentAttachmentRepo struct {
}

func NewCommentAttachmentRepo() CommentAttachmentRepo {
	return &sCommentAttachmentRepo{}
}

func (r *sCommentAttachmentRepo) Insert(ctx context.Context, e *model.CommentAttachment) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *sCommentAttachmentRepo) FindByID(ctx context.Context, id string) (*CommentAttachmentData, error) {
	e := &CommentAttachmentData{}
	err := cockroach.FindOne(ctx, e, "id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("sCommentAttachmentRepo.cockroach.FindOne: %w", err)
	}

	return e, nil
}

func (r *sCommentAttachmentRepo) SoftDelete(ctx context.Context, id string) error {
	sql := "UPDATE comment_attachments SET deleted_at = NOW() WHERE id = $1;"

	cmd, err := cockroach.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("comment_attachments cockroach.Exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("*sCommentAttachmentRepo not found any records to delete")
	}

	return nil
}

// SearchCommentAttachmentOpts all params is options
type SearchCommentAttachmentOpts struct {
	IDs []string
	// todo add more search options
	Limit  int64
	Offset int64
	Sort   *Sort
}

func (s *SearchCommentAttachmentOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.CommentAttachmentFieldID)
	}
	// todo add more search options example:
	//if s.Name != "" {
	//	args = append(args, "%"+s.Name+"%")
	//	conds += fmt.Sprintf(" AND (b.%[2]s ILIKE $%[1]d OR b.%[3]s ILIKE $%[1]d)",
	//		len(args), model.CommentAttachmentFieldName, model.CommentAttachmentFieldCode)
	//}
	//if s.Code != "" {
	//	args = append(args, s.Code)
	//	conds += fmt.Sprintf(" AND b.%s ILIKE $%d", model.CommentAttachmentFieldCode, len(args))
	//}

	b := &model.CommentAttachment{}
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

type CommentAttachmentData struct {
	*model.CommentAttachment
}

func (r *sCommentAttachmentRepo) Search(ctx context.Context, s *SearchCommentAttachmentOpts) ([]*CommentAttachmentData, error) {
	CommentAttachment := make([]*CommentAttachmentData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&CommentAttachment)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return CommentAttachment, nil
}

func (r *sCommentAttachmentRepo) Count(ctx context.Context, s *SearchCommentAttachmentOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("sCommentAttachmentRepo.Count: %w", err)
	}

	return countResult, nil
}
