package repository

import (
	"context"
	"fmt"
	"strings"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
)

type RegionRepo interface {
	Insert(ctx context.Context, e *model.Region) error
	Update(ctx context.Context, e *model.Region) error
	FindByID(ctx context.Context, id string) (*model.Region, error)
	SoftDelete(ctx context.Context, id string) error
	Search(ctx context.Context, s *SearchRegionsOpts) ([]*model.Region, error)
}

type regionRepo struct {
}

func NewRegionRepo() RegionRepo {
	return &regionRepo{}
}

func (r *regionRepo) Insert(ctx context.Context, e *model.Region) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *regionRepo) FindByID(ctx context.Context, id string) (*model.Region, error) {
	e := &model.Region{}
	err := cockroach.FindOne(ctx, e, "id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("cockroach.FindOne: %w", err)
	}

	return e, nil
}

func (r *regionRepo) Update(ctx context.Context, e *model.Region) error {
	return cockroach.Update(ctx, e)
}

func (r *regionRepo) SoftDelete(ctx context.Context, id string) error {
	sql := `UPDATE regions
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
type SearchRegionsOpts struct {
	IDs      []int64
	Name     string
	Search   string
	ParentID int64
	Limit    int64
	Offset   int64
}

func (s *SearchRegionsOpts) buildQuery() (string, []interface{}) {
	args := []interface{}{}
	conds := ""
	joins := ""

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND r.%s = ANY($1)", model.RegionFieldID)
	}

	if s.Name != "" {
		args = append(args, "%"+s.Name+"%")
		conds += fmt.Sprintf(" AND r.%s ILIKE $%d", model.RegionFieldName, len(args))
	}
	if s.Search != "" {
		args = append(args, "%"+s.Search+"%")
		conds += fmt.Sprintf(" AND r.%s ILIKE $%d", model.RegionFieldName, len(args))
	}

	if s.ParentID >= 0 {
		args = append(args, s.ParentID)
		conds += fmt.Sprintf(" AND r.%s = $%d", model.RegionFieldParentID, len(args))
	}
	r := &model.Region{}
	fields, _ := r.FieldMap()

	return fmt.Sprintf(`SELECT r.%s
		FROM %s AS r %s
		WHERE TRUE %s
		ORDER BY r.display_order ASC
		LIMIT %d
		OFFSET %d`, strings.Join(fields, ", r."), r.TableName(), joins, conds, s.Limit, s.Offset), args
}

func (r *regionRepo) Search(ctx context.Context, s *SearchRegionsOpts) ([]*model.Region, error) {
	regions := make([]*model.Region, 0)
	sql, args := s.buildQuery()
	err := cockroach.Select(ctx, sql, args...).ScanAll(&regions)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return regions, nil
}
