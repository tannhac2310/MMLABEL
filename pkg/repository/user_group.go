package repository

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
)

type UserGroupRepo interface {
	Insert(ctx context.Context, e *model.UserGroup) error
	FindByID(ctx context.Context, id string) (*model.UserGroup, error)
	FindByUserID(ctx context.Context, id string) ([]*model.UserGroup, error)
	FindByGroupID(ctx context.Context, id string) ([]*model.UserGroup, error)
	DeleteByUserIDAndGroupIDs(ctx context.Context, userID string, groupIDs []string) error
}

type userGroupRepo struct {
}

func NewUserGroupRepo() UserGroupRepo {
	return &userGroupRepo{}
}

func (r *userGroupRepo) Insert(ctx context.Context, e *model.UserGroup) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("cockroach.Insert: %w", err)
	}

	return nil
}

func (r *userGroupRepo) FindByID(ctx context.Context, id string) (*model.UserGroup, error) {
	e := &model.UserGroup{}
	err := cockroach.FindOne(ctx, e, "id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("cockroach.FindOne: %w", err)
	}

	return e, nil
}

func (r *userGroupRepo) FindByUserID(ctx context.Context, id string) ([]*model.UserGroup, error) {
	result := make([]*model.UserGroup, 0)
	err := cockroach.FindMany(
		ctx,
		&model.UserGroup{},
		&result,
		fmt.Sprintf("%s = $1", model.UserGroupFieldUserID),
		id)
	if err != nil {
		return nil, fmt.Errorf("cockroach.FindMany: %w", err)
	}

	return result, nil
}

func (r *userGroupRepo) DeleteByUserIDAndGroupIDs(ctx context.Context, userID string, groupIDs []string) error {
	sql := `UPDATE user_group
		SET deleted_at = NOW()
		WHERE user_id = $1 AND group_id = ANY($2)`

	cmd, err := cockroach.Exec(ctx, sql, userID, groupIDs)
	if err != nil {
		return fmt.Errorf("cockroach.Exec: %w", err)
	}

	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("not found any records to delete")
	}

	return nil
}

func (r *userGroupRepo) FindByGroupID(ctx context.Context, id string) ([]*model.UserGroup, error) {
	result := make([]*model.UserGroup, 0)
	err := cockroach.FindMany(
		ctx,
		&model.UserGroup{},
		&result,
		fmt.Sprintf("%s = $1", model.UserGroupFieldGroupID),
		id)
	if err != nil {
		return nil, fmt.Errorf("cockroach.FindMany: %w", err)
	}

	return result, nil
}
