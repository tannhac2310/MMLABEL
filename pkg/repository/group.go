package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
)

type GroupRepo interface {
	Insert(ctx context.Context, e *model.Group) error
	Update(ctx context.Context, e *model.Group) error
	FindByID(ctx context.Context, id string) (*model.Group, error)
	Search(ctx context.Context, s *SearchGroupOpts) ([]*model.Group, error)
}

type groupRepo struct {
}

func NewGroupRepo() GroupRepo {
	return &groupRepo{}
}

func (r *groupRepo) Insert(ctx context.Context, e *model.Group) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("cockroach.Create: %w", err)
	}

	return nil
}

func (r *groupRepo) FindByID(ctx context.Context, id string) (*model.Group, error) {
	e := &model.Group{}
	err := cockroach.FindOne(ctx, e, "id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("cockroach.FindOne: %w", err)
	}

	return e, nil
}

func (r *groupRepo) Update(ctx context.Context, e *model.Group) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

type SearchGroupOpts struct {
	IDs    []string
	Name   string
	Limit  int64
	Offset int64
}

func (r *groupRepo) Search(ctx context.Context, s *SearchGroupOpts) ([]*model.Group, error) {
	ponds := make([]*model.Group, 0)
	sql, args := s.buildQuery()

	err := cockroach.Select(ctx, sql, args...).ScanAll(&ponds)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return ponds, nil
}

func (s *SearchGroupOpts) buildQuery() (string, []interface{}) {
	args := []interface{}{}
	conds := ""
	joins := ""

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND g.%s = ANY($1)", model.GroupFieldID)
	}

	if s.Name != "" {
		args = append(args, "%"+s.Name+"%")
		conds += fmt.Sprintf(" AND g.%s ILIKE $%d", model.GroupFieldName, len(args))
	}

	e := &model.Group{}
	fields, _ := e.FieldMap()

	return fmt.Sprintf(`SELECT g.%s
		FROM %s AS g %s
		WHERE TRUE %s AND g.deleted_at IS NULL
		ORDER BY g.id DESC
		LIMIT %d
		OFFSET %d`, strings.Join(fields, ", g."), e.TableName(), joins, conds, s.Limit, s.Offset), args
}
