package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
)

type UserRepo interface {
	Insert(ctx context.Context, e *model.User) error
	Update(ctx context.Context, e *model.User) error
	SoftDelete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*model.User, error)
	FindByPhoneOrEmail(ctx context.Context, k string) (*UserData, error)
	Search(ctx context.Context, s *SearchUsersOpts) ([]*UserData, error)
	Count(ctx context.Context, s *SearchUsersOpts) (*CountResult, error)
}

type userRepo struct {
}
type UserData struct {
	*model.User
}

func NewUserRepo() UserRepo {
	return &userRepo{}
}

func (r *userRepo) SoftDelete(ctx context.Context, id string) error {
	sql := `UPDATE users
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

func (r *userRepo) Insert(ctx context.Context, e *model.User) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *userRepo) FindByID(ctx context.Context, id string) (*model.User, error) {
	e := &model.User{}
	err := cockroach.FindOne(ctx, e, "id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("UserRepo:cockroach.FindOne: %w, %s", err, id)
	}
	return e, nil
}

func (r *userRepo) Update(ctx context.Context, e *model.User) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

// SearchUsersOpts all params is options
type SearchUsersOpts struct {
	IDs         []string
	NotIDs      []string
	NotRoleIDs  []string
	Name        string
	Department  string
	Search      string
	PhoneNumber string
	Email       string
	GroupID     string
	RoleID      string
	Type        enum.UserType
	Limit       int64
	Offset      int64
}

func (s *SearchUsersOpts) buildQuery(isCount bool) (string, []interface{}) {
	args := []interface{}{}
	conds := ""
	joins := ""

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND u.%s = ANY($1)", model.UserFieldID)
	}
	if len(s.NotIDs) > 0 {
		args = append(args, s.NotIDs)
		conds += fmt.Sprintf(" AND NOT u.%s = ANY($1)", model.UserFieldID)
	}

	if len(s.NotRoleIDs) > 0 {
		userRoleTable := model.UserRole{}
		args = append(args, s.NotRoleIDs)
		joins += fmt.Sprintf(`LEFT JOIN %[1]s 
																				ON %[1]s.%[2]s = u.id
																				AND %[1]s.deleted_at IS NULL 
																				AND %[1]s.%[3]s = ANY($%d) `,
			userRoleTable.TableName(), model.UserRoleFieldUserID, model.UserRoleFieldRoleID, len(args),
		)

		conds += fmt.Sprintf(" AND %s.%s IS NULL", userRoleTable.TableName(), model.UserRoleFieldUserID)
	}

	if s.Name != "" {
		args = append(args, "%"+s.Name+"%")
		conds += fmt.Sprintf(" AND u.%s ILIKE $%d", model.UserFieldName, len(args))
	}
	if s.Department != "" {
		var departs = []string{s.Department}
		if s.Department == "CBD" {
			departs = []string{
				"CBD", "CA", "NLC", "NLT", "CLD", "CLM", "CTP",
				"CTAY", "CVT", "BE", "BCD", "BCM","PCN", "BTD", "BTM",
				"BKE", "DA", "DD", "DM", "DNN", "DNKT", "UDE"}
		}
		args = append(args, departs)
		conds += fmt.Sprintf(" AND NOT u.%s = ANY($1)", model.UserFieldDepartments)
		// args = append(args, "%"+s.Department+"%")
		// conds += fmt.Sprintf(" AND u.%s ILIKE $%d", model.UserFieldDepartments, len(args))
	}
	fmt.Println("Minh:", s.Department)
	if s.Search != "" {
		args = append(args, "%"+s.Search+"%")
		conds += fmt.Sprintf(" AND (u.%s ILIKE $%d", model.UserFieldName, len(args))
		conds += fmt.Sprintf(" OR u.%s ILIKE $%d", model.UserFieldPhoneNumber, len(args))
		conds += fmt.Sprintf(" OR u.%s ILIKE $%d) ", model.UserFieldEmail, len(args))
	}

	if s.PhoneNumber != "" {
		args = append(args, "%"+s.PhoneNumber+"%")
		conds += fmt.Sprintf(" AND u.%s ILIKE $%d", model.UserFieldPhoneNumber, len(args))
	}

	if s.Email != "" {
		args = append(args, s.Email)
		conds += fmt.Sprintf(" AND u.%s = $%d", model.UserFieldEmail, len(args))
	}
	if s.Type > 0 {
		args = append(args, s.Type)
		conds += fmt.Sprintf(" AND u.%s = $%d", model.UserFieldType, len(args))
	}

	if s.RoleID != "" {
		userRole := &model.UserRole{}
		tableName := userRole.TableName()

		args = append(args, s.RoleID)
		conds += fmt.Sprintf(" AND %s.%s = $%d", tableName, model.UserRoleFieldRoleID, len(args))

		joins += fmt.Sprintf(" INNER JOIN %[1]s AS %[1]s ON %[1]s.%[2]s = u.id AND %[1]s.deleted_at IS NULL", tableName, model.UserRoleFieldUserID)
	}

	if s.GroupID != "" {
		userGroup := &model.UserGroup{}
		tableName := userGroup.TableName()

		args = append(args, s.RoleID)
		conds += fmt.Sprintf(" AND %s.%s = $%d", tableName, model.UserGroupFieldGroupID, len(args))

		joins += fmt.Sprintf(" INNER JOIN %[1]s AS %[1]s ON %[1]s.%[2]s = u.id AND %[1]s.deleted_at IS NULL", tableName, model.UserGroupFieldUserID)
	}
	
	u := &model.User{}
	fields, _ := u.FieldMap()
	if isCount {
		return fmt.Sprintf(`SELECT count(1) as cnt
		FROM %s AS u %s
		WHERE TRUE %s AND u.deleted_at IS NULL`, u.TableName(), joins, conds), args
	}
	return fmt.Sprintf(`SELECT u.%s
		FROM %s AS u %s
		WHERE TRUE %s AND u.deleted_at IS NULL
		ORDER BY u.id DESC
		LIMIT %d
		OFFSET %d`, strings.Join(fields, ", u."), u.TableName(), joins, conds, s.Limit, s.Offset), args
}

func (r *userRepo) Search(ctx context.Context, s *SearchUsersOpts) ([]*UserData, error) {
	users := make([]*UserData, 0)
	sql, args := s.buildQuery(false)
	fmt.Println(sql)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&users)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return users, nil
}

func (r *userRepo) FindByPhoneOrEmail(ctx context.Context, k string) (*UserData, error) {
	e := &UserData{}
	err := cockroach.FindOne(
		ctx,
		e,
		fmt.Sprintf("%s = $1 OR %s = $1", model.UserFieldEmail, model.UserFieldPhoneNumber),
		k)
	if err != nil {
		return nil, fmt.Errorf("cockroach.FindOne: %w", err)
	}

	return e, nil
}

func (r *userRepo) Count(ctx context.Context, s *SearchUsersOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("syllabus.Count: %w", err)
	}

	return countResult, nil
}
