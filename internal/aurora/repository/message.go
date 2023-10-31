package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type MessageRepo interface {
	Insert(ctx context.Context, e *model.Message) error
	Update(ctx context.Context, e *model.Message) error
	SoftDelete(ctx context.Context, id string) error
	Search(ctx context.Context, s *SearchMessagesOpts) ([]*MessageData, error)
	SearchOne(ctx context.Context, s *SearchMessagesOpts) (*MessageData, error)
	Count(ctx context.Context, s *SearchMessagesOpts) (*CountResult, error)
}

type imMessagesRepo struct {
}

func NewMessageRepo() MessageRepo {
	return &imMessagesRepo{}
}

func (r *imMessagesRepo) Insert(ctx context.Context, e *model.Message) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *imMessagesRepo) Update(ctx context.Context, e *model.Message) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

func (r *imMessagesRepo) SoftDelete(ctx context.Context, id string) error {
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

// SearchMessagesOpts all params is options
type SearchMessagesOpts struct {
	IDs    []string
	ChatID string
	UserID string
	Limit  int64
	Offset int64
	Sort   *Sort
}

func (s *SearchMessagesOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := `
join message_relations mr on mr.chat_id = b.chat_id
join chats c on c.id = b.chat_id
`
	// check permission user can read message
	args = append(args, s.UserID)
	conds += fmt.Sprintf(" AND mr.%s = $%d", model.MessageRelationFieldUserID, len(args))

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($%d)", model.MessageFieldID, len(args))
	}

	if s.ChatID != "" {
		args = append(args, s.ChatID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.MessageFieldChatID, len(args))
	}

	b := &model.Message{}
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
	return fmt.Sprintf(`SELECT b.%s, c.title as chat_title
		FROM %s AS b %s
		WHERE TRUE %s AND b.deleted_at IS NULL
		%s
		LIMIT %d
		OFFSET %d`, strings.Join(fields, ", b."), b.TableName(), joins, conds, order, s.Limit, s.Offset), args
}

type MessageData struct {
	*model.Message
	ChatTitle sql.NullString `db:"chat_title"`
}

func (r *imMessagesRepo) Search(ctx context.Context, s *SearchMessagesOpts) ([]*MessageData, error) {
	message := make([]*MessageData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&message)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return message, nil
}

func (r *imMessagesRepo) SearchOne(ctx context.Context, s *SearchMessagesOpts) (*MessageData, error) {
	message := &MessageData{}
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanOne(&message)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return message, nil
}

func (r *imMessagesRepo) Count(ctx context.Context, s *SearchMessagesOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("messages.Count: %w", err)
	}

	return countResult, nil
}
