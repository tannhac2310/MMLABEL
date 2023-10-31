package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type MessageRelationRepo interface {
	Insert(ctx context.Context, e *model.MessageRelation) error
	Upsert(ctx context.Context, chatID, userID string) error
	Update(ctx context.Context, e *model.MessageRelation) error
	SoftDelete(ctx context.Context, id string) error
	Search(ctx context.Context, s *SearchMessageRelationOpts) ([]*MessageRelationData, error)
	Count(ctx context.Context, s *SearchMessageRelationOpts) (*CountResult, error)
	SearchOne(ctx context.Context, s *SearchMessageRelationOpts) (*MessageRelationData, error)
}

type messageRelationsRepo struct {
}

func NewMessageRelationRepo() MessageRelationRepo {
	return &messageRelationsRepo{}
}

func (r *messageRelationsRepo) Insert(ctx context.Context, e *model.MessageRelation) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}
func (r *messageRelationsRepo) Upsert(ctx context.Context, chatID, userID string) error {
	sql := `INSERT INTO message_relations (chat_id, user_id, deleted_at)
    VALUES ($1, $2, $3)
    ON CONFLICT (chat_id, user_id, deleted_at)
    DO UPDATE SET chat_id = excluded.chat_id, user_id = excluded.user_id`

	cmd, err := cockroach.Exec(ctx, sql, chatID, userID)
	if err != nil {
		return fmt.Errorf("cockroach.Exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("not found any records to delete")
	}

	return nil
}

func (r *messageRelationsRepo) Update(ctx context.Context, e *model.MessageRelation) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}
func (r *messageRelationsRepo) SearchOne(ctx context.Context, s *SearchMessageRelationOpts) (*MessageRelationData, error) {
	message := &MessageRelationData{}
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanOne(&message)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return message, nil
}
func (r *messageRelationsRepo) SoftDelete(ctx context.Context, id string) error {
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

// SearchMessageRelationOpts all params is options
type SearchMessageRelationOpts struct {
	IDs    []string
	ChatID string
	UserID string
	Limit  int64
	Offset int64
	Sort   *Sort
}

func (s *SearchMessageRelationOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.MessageRelationFieldID)
	}

	if s.ChatID != "" {
		args = append(args, s.ChatID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.MessageRelationFieldChatID, len(args))
	}

	if s.UserID != "" {
		args = append(args, s.UserID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.MessageRelationFieldUserID, len(args))
	}

	b := &model.MessageRelation{}
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

type MessageRelationData struct {
	*model.MessageRelation
}

func (r *messageRelationsRepo) Search(ctx context.Context, s *SearchMessageRelationOpts) ([]*MessageRelationData, error) {
	message := make([]*MessageRelationData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&message)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return message, nil
}

func (r *messageRelationsRepo) Count(ctx context.Context, s *SearchMessageRelationOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("messages_relation.Count: %w", err)
	}

	return countResult, nil
}
