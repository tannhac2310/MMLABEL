package repository

import (
	"context"
	"fmt"
	"strings"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
)

type UserNotificationRepo interface {
	Insert(ctx context.Context, e *model.UserNotification) error
	UpdateNotificationsStatus(ctx context.Context, IDs []string, userID string, status enum.NotificationStatus) error
	FindByID(ctx context.Context, id string) (*model.UserNotification, error)
	Search(ctx context.Context, s *SearchNotificationOtps) ([]*model.UserNotification, error)
	CountUserNotifications(ctx context.Context, s *CountNotificationOtps) (count int64, err error)
	FindByUserID(ctx context.Context, id string) ([]*model.UserNotification, error)
}

type userNotificationRepo struct {
}

func NewUserNotificationRepo() UserNotificationRepo {
	return &userNotificationRepo{}
}

func (r *userNotificationRepo) UpdateNotificationsStatus(ctx context.Context, ids []string, userID string, status enum.NotificationStatus) error {
	sql := `UPDATE user_notifications
		SET updated_at = NOW(),
			status = $3
		WHERE id = ANY($1) AND user_id = $2`

	cmd, err := cockroach.Exec(ctx, sql, ids, userID, status)
	if err != nil {
		return fmt.Errorf("r.c.UpdateOne: %w", err)
	}

	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("not found records to update")
	}

	return nil
}

func (r *userNotificationRepo) CountUserNotifications(ctx context.Context, s *CountNotificationOtps) (count int64, err error) {
	args := []interface{}{}
	conds := ""

	if s.UserID != "" {
		args = append(args, s.UserID)
		conds += "AND user_id = $1 "
	}

	if s.Type > 0 {
		args = append(args, s.Type)
		conds += fmt.Sprintf("AND type = $%d ", len(args))
	}

	if s.Status > 0 {
		args = append(args, s.Status)
		conds += fmt.Sprintf("AND status = $%d ", len(args))
	}

	sql := "SELECT COUNT(id) FROM user_notifications WHERE TRUE " + conds

	var total int64
	err = cockroach.QueryRow(ctx, sql, args...).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("cockroach.QueryRow: %w", err)
	}

	return total, nil
}

type CountNotificationOtps struct {
	UserID string
	Type   enum.NotificationType
	Status enum.NotificationStatus
}

type SearchNotificationOtps struct {
	IDs    []string
	UserID string
	Type   enum.NotificationType
	Status enum.NotificationStatus
	Limit  int64
	Offset int64
}

func (s *SearchNotificationOtps) buildQuery() (string, []interface{}) {
	args := []interface{}{}
	conds := ""

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += "AND id = ANY($1) "
	}

	if s.UserID != "" {
		args = append(args, s.UserID)
		conds += fmt.Sprintf("AND user_id = $%d ", len(args))
	}

	if s.Type > 0 {
		args = append(args, s.Type)
		conds += fmt.Sprintf("AND type = $%d ", len(args))
	}

	if s.Status > 0 {
		args = append(args, s.Status)
		conds += fmt.Sprintf("AND status = $%d ", len(args))
	}

	e := &model.UserNotification{}
	fields, _ := e.FieldMap()

	return fmt.Sprintf(`SELECT %s
		FROM user_notifications
		WHERE TRUE %s AND deleted_at IS NULL
		ORDER BY id DESC
		LIMIT %d
		OFFSET %d`,
		strings.Join(fields, ","),
		conds,
		s.Limit,
		s.Offset,
	), args
}
func (r *userNotificationRepo) Search(ctx context.Context, s *SearchNotificationOtps) ([]*model.UserNotification, error) {
	sql, args := s.buildQuery()
	results := make([]*model.UserNotification, 0)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&results)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return results, nil
}

func (r *userNotificationRepo) Insert(ctx context.Context, e *model.UserNotification) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("cockroach.Create: %w", err)
	}

	return nil
}

func (r *userNotificationRepo) FindByID(ctx context.Context, id string) (*model.UserNotification, error) {
	e := &model.UserNotification{}
	err := cockroach.FindOne(ctx, e, "id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("cockroach.FindOne: %w", err)
	}

	return e, nil
}

func (r *userNotificationRepo) FindByUserID(ctx context.Context, id string) ([]*model.UserNotification, error) {
	e := &model.UserNotification{}
	result := make([]*model.UserNotification, 0)
	err := cockroach.FindMany(
		ctx,
		e,
		&result,
		fmt.Sprintf("%s = $1", model.UserNamePasswordFieldUserID),
		id)
	if err != nil {
		return nil, fmt.Errorf("cockroach.FindMany: %w", err)
	}

	return result, nil
}
