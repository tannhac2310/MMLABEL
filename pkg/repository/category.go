package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
)

type CategoryRepo interface {
	Insert(ctx context.Context, e *model.Category) error
	Update(ctx context.Context, e *model.Category) error
	FindByID(ctx context.Context, id string) (*model.Category, error)
	SoftDelete(ctx context.Context, id string) error
	Search(ctx context.Context, s *SearchCategoriesOpts) ([]*model.Category, error)
}

type categorieRepo struct {
}

func NewCategoryRepo() CategoryRepo {
	return &categorieRepo{}
}

func (r *categorieRepo) Insert(ctx context.Context, e *model.Category) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *categorieRepo) FindByID(ctx context.Context, id string) (*model.Category, error) {
	e := &model.Category{}
	err := cockroach.FindOne(ctx, e, "id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("CategoryRepo:cockroach.FindOne: %w", err)
	}

	return e, nil
}

func (r *categorieRepo) Update(ctx context.Context, e *model.Category) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

func (r *categorieRepo) SoftDelete(ctx context.Context, id string) error {
	sql := `UPDATE category
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

// all params is options
type SearchCategoriesOpts struct {
	IDs         []string
	Name        string
	Search      string
	Description string
	Limit       int64
	Offset      int64
}

func (s *SearchCategoriesOpts) buildQuery() (string, []interface{}) {
	args := []interface{}{}
	conds := ""
	joins := ""

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND c.%s = ANY($1)", model.CategoryFieldID)
	}

	if s.Name != "" {
		args = append(args, "%"+s.Name+"%")
		conds += fmt.Sprintf(" AND c.%s ILIKE $%d", model.CategoryFieldName, len(args))
	}
	if s.Search != "" {
		args = append(args, "%"+s.Search+"%")
		conds += fmt.Sprintf(" AND c.%s ILIKE $%d", model.CategoryFieldName, len(args))
	}

	if s.Description != "" {
		args = append(args, "%"+s.Description+"%")
		conds += fmt.Sprintf(" AND c.%s ILIKE $%d", model.CategoryFieldDescription, len(args))
	}

	c := &model.Category{}
	fields, _ := c.FieldMap()

	return fmt.Sprintf(`SELECT c.%s
		FROM %s AS c %s
		WHERE TRUE %s AND c.deleted_at IS NULL
		ORDER BY c.id DESC
		LIMIT %d
		OFFSET %d`, strings.Join(fields, ", c."), c.TableName(), joins, conds, s.Limit, s.Offset), args
}

func (r *categorieRepo) Search(ctx context.Context, s *SearchCategoriesOpts) ([]*model.Category, error) {
	categories := make([]*model.Category, 0)
	sql, args := s.buildQuery()
	err := cockroach.Select(ctx, sql, args...).ScanAll(&categories)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return categories, nil
}
