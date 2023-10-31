package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
)

type UserRoleRepo interface {
	Insert(ctx context.Context, e *model.UserRole) error
	Search(ctx context.Context, s *SearchUserRoleOpts) ([]*RuleUsersData, error)
	Count(ctx context.Context, s *SearchUserRoleOpts) (*CountResult, error)
	FindByUserID(ctx context.Context, id string) ([]*model.UserRole, error)
	FindByRoleID(ctx context.Context, id string) ([]*model.UserRole, error)
	DeleteByUserIDAndRoleIDs(ctx context.Context, userID string, roleIDs []string) error
	DeleteByUserIDsAndRoleID(ctx context.Context, userIDs []string, roleID string) error
}
type SearchUserRoleOpts struct {
	RoleIDs   []string
	UserIDs   []string
	UserTypes []enum.UserType
	Search    string
	Limit     int64
	Offset    int64
}
type userRoleRepo struct {
}

func NewUserRoleRepo() UserRoleRepo {
	return &userRoleRepo{}
}

func (r *userRoleRepo) Insert(ctx context.Context, e *model.UserRole) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("cockroach.Insert: %w", err)
	}

	return nil
}

func (r *userRoleRepo) FindByID(ctx context.Context, id string) (*model.UserRole, error) {
	e := &model.UserRole{}
	err := cockroach.FindOne(ctx, e, "id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("cockroach.FindOne: %w", err)
	}

	return e, nil
}

func (r *userRoleRepo) FindByUserID(ctx context.Context, id string) ([]*model.UserRole, error) {
	result := make([]*model.UserRole, 0)
	err := cockroach.FindMany(
		ctx,
		&model.UserRole{},
		&result,
		fmt.Sprintf("%s = $1", model.UserRoleFieldUserID),
		id)
	if err != nil {
		return nil, fmt.Errorf("cockroach.FindMany: %w", err)
	}

	return result, nil
}

func (r *userRoleRepo) DeleteByUserIDAndRoleIDs(ctx context.Context, userID string, roleIDs []string) error {
	sql := `UPDATE user_role
		SET deleted_at = NOW()
		WHERE user_id = $1 AND role_id = ANY($2)`

	cmd, err := cockroach.Exec(ctx, sql, userID, roleIDs)
	if err != nil {
		return fmt.Errorf("cockroach.Exec: %w", err)
	}

	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("not found any records to delete")
	}

	return nil
}
func (r *userRoleRepo) DeleteByUserIDsAndRoleID(ctx context.Context, userIDs []string, roleID string) error {
	sql := `UPDATE user_role
		SET deleted_at = NOW()
		WHERE user_id = ANY($1) AND role_id = $2`

	cmd, err := cockroach.Exec(ctx, sql, userIDs, roleID)
	if err != nil {
		return fmt.Errorf("cockroach.Exec: %w", err)
	}

	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("not found any records to delete")
	}

	return nil
}

func (r *userRoleRepo) FindByRoleID(ctx context.Context, id string) ([]*model.UserRole, error) {
	result := make([]*model.UserRole, 0)
	err := cockroach.FindMany(
		ctx,
		&model.UserRole{},
		&result,
		fmt.Sprintf("%s = $1", model.UserRoleFieldRoleID),
		id)
	if err != nil {
		return nil, fmt.Errorf("cockroach.FindMany: %w", err)
	}

	return result, nil
}

type RuleUsersData struct {
	ID            string          `db:"id"` // user_id
	RoleID        string          `db:"role_id"`
	RoleName      string          `db:"role_name"`
	UserID        string          `db:"user_id"`
	Name          string          `db:"name"`
	Avatar        string          `db:"avatar"`
	PhoneNumber   string          `db:"phone_number"`
	Email         string          `db:"email"`
	Address       string          `db:"address"`
	Type          enum.UserType   `db:"type"`
	Status        enum.UserStatus `db:"status"`
	CreatedBy     sql.NullString  `db:"created_by"`
	CreatedByName sql.NullString  `db:"created_by_name"`
	CreatedAt     time.Time       `db:"created_at"`
}

func (r *userRoleRepo) Search(ctx context.Context, s *SearchUserRoleOpts) ([]*RuleUsersData, error) {
	message := make([]*RuleUsersData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&message)
	if err != nil {
		return nil, fmt.Errorf("userRoleRepo.Search: %w", err)
	}

	return message, nil
}
func (r *userRoleRepo) Count(ctx context.Context, s *SearchUserRoleOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("userRoleRepo.Count: %w", err)
	}

	return countResult, nil
}

func (s *SearchUserRoleOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ` 
JOIN users u on u.id = b.user_id 
JOIN roles r on r.id = b.role_id 
LEFT JOIN users cb on cb.id = b.created_by
`

	if len(s.UserIDs) > 0 {
		args = append(args, s.UserIDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.UserRoleFieldUserID)
	}
	if len(s.RoleIDs) > 0 {
		args = append(args, s.RoleIDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.UserRoleFieldRoleID)
	}
	if len(s.UserTypes) > 0 {
		args = append(args, s.UserTypes)
		conds += fmt.Sprintf(" AND u.%s = ANY($1)", model.UserFieldType)
	}
	if s.Search != "" {
		args = append(args, "%"+s.Search+"%")
		conds += fmt.Sprintf(" AND (u.%s ILIKE $%d", model.UserFieldName, len(args))
		conds += fmt.Sprintf(" OR u.%s ILIKE $%d", model.UserFieldPhoneNumber, len(args))
		conds += fmt.Sprintf(" OR u.%s ILIKE $%d) ", model.UserFieldEmail, len(args))
	}

	b := &model.UserRole{}
	if isCount {
		return fmt.Sprintf(`SELECT count(*) as cnt
		FROM %s AS b %s
		WHERE TRUE %s AND b.deleted_at IS NULL`, b.TableName(), joins, conds), args
	}

	order := " ORDER BY b.id DESC "

	return fmt.Sprintf(`SELECT 
		b.role_id,
		r.name as role_name,
		u.id,
		u.name,
		u.avatar,
		u.phone_number,
		u.email,
		u.address,
		u.type,
		u.status,
		b.created_by,
		cb.name as created_by_name,
		b.created_at
		FROM %s AS b %s
		WHERE TRUE %s AND b.deleted_at IS NULL
		%s
		LIMIT %d
		OFFSET %d`, b.TableName(), joins, conds, order, s.Limit, s.Offset), args
}
