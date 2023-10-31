package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

type ChatRepo interface {
	Insert(ctx context.Context, e *model.Chat) error
	Update(ctx context.Context, e *model.Chat) error
	SoftDelete(ctx context.Context, id string) error
	Search(ctx context.Context, s *SearchChatsOpts) ([]*ChatData, error)
	Count(ctx context.Context, s *SearchChatsOpts) (*CountResult, error)
	SearchOne(ctx context.Context, s *SearchChatsOpts) (*ChatData, error)
}

type chatsRepo struct {
}

func NewChatRepo() ChatRepo {
	return &chatsRepo{}
}

func (r *chatsRepo) Insert(ctx context.Context, e *model.Chat) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *chatsRepo) Update(ctx context.Context, e *model.Chat) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

func (r *chatsRepo) SoftDelete(ctx context.Context, id string) error {
	sql := `UPDATE message
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

// SearchChatsOpts all params is options
type SearchChatsOpts struct {
	IDs        []string
	ID         string
	AuthorID   string
	EntityType enum.ChatEntityType
	EntityID   string
	Limit      int64
	Offset     int64
	Sort       *Sort
}

func (s *SearchChatsOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.ChatFieldID)
	}

	if s.ID != "" {
		args = append(args, s.ID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.ChatFieldID, len(args))
	}
	if s.AuthorID != "" {
		args = append(args, s.AuthorID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.ChatFieldAuthorID, len(args))
	}

	if s.EntityType > 0 {
		args = append(args, s.EntityType)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.ChatFieldEntityType, len(args))
	}
	if s.EntityID != "" {
		args = append(args, s.EntityID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.ChatFieldEntityID, len(args))
	}

	b := &model.Chat{}
	fields, _ := b.FieldMap()
	if isCount {
		return fmt.Sprintf(`SELECT count(*) as cnt
		FROM %s AS b %s
		WHERE TRUE %s AND b.deleted_at IS NULL`, b.TableName(), joins, conds), args
	}

	order := " ORDER BY b.id DESC "
	if s.Sort != nil {
		order = fmt.Sprintf(" ORDER BY b.%s %s", s.Sort.By, s.Sort.Order)
	}
	return fmt.Sprintf(`SELECT b.%s
		FROM %s AS b %s
		WHERE TRUE %s AND b.deleted_at IS NULL
		%s
		LIMIT %d
		OFFSET %d`, strings.Join(fields, ", b."), b.TableName(), joins, conds, order, s.Limit, s.Offset), args
}

type ChatData struct {
	*model.Chat
}

func (r *chatsRepo) Search(ctx context.Context, s *SearchChatsOpts) ([]*ChatData, error) {
	message := make([]*ChatData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&message)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return message, nil
}
func (r *chatsRepo) SearchOne(ctx context.Context, s *SearchChatsOpts) (*ChatData, error) {
	message := &ChatData{}
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanOne(message)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select1: %w", err)
	}

	return message, nil
}

func (r *chatsRepo) Count(ctx context.Context, s *SearchChatsOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("chat.Count: %w", err)
	}

	return countResult, nil
}
