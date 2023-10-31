package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
)

type BannerRepo interface {
	Insert(ctx context.Context, e *model.Banner) error
	Update(ctx context.Context, e *model.Banner) error
	SoftDelete(ctx context.Context, id string) error
	Search(ctx context.Context, s *SearchBannersOpts) ([]*BannerData, error)
	Count(ctx context.Context, s *SearchBannersOpts) (*CountResult, error)
}

type bannerRepo struct {
}
type BannerData struct {
	*model.Banner
	CreatedByName string `db:"created_by_name"`
}

func (r *bannerRepo) Count(ctx context.Context, s *SearchBannersOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("stage.Count: %w", err)
	}

	return countResult, nil
}

func NewBannerRepo() BannerRepo {
	return &bannerRepo{}
}

func (r *bannerRepo) Insert(ctx context.Context, e *model.Banner) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *bannerRepo) Update(ctx context.Context, e *model.Banner) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

func (r *bannerRepo) SoftDelete(ctx context.Context, id string) error {
	sql := `UPDATE banners
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

// SearchBannersOpts all params is options
type SearchBannersOpts struct {
	IDs    []string
	Name   string
	Limit  int64
	Offset int64
}

func (s *SearchBannersOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joinFields := ` cb.name as created_by_name `
	joins := ` JOIN users cb on cb.id = b.created_by `

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.BannerFieldID)
	}

	if s.Name != "" {
		args = append(args, "%"+s.Name+"%")
		conds += fmt.Sprintf(" AND b.%s ILIKE $%d", model.BannerFieldName, len(args))
	}

	b := &model.Banner{}
	fields, _ := b.FieldMap()
	if isCount {
		return fmt.Sprintf(`SELECT count(1) as cnt
		FROM %s AS b
		WHERE TRUE %s AND b.deleted_at IS NULL`, b.TableName(), conds), args
	}
	return fmt.Sprintf(`SELECT b.%s, %s
		FROM %s AS b %s
		WHERE TRUE %s AND b.deleted_at IS NULL
		ORDER BY b.id DESC
		LIMIT %d
		OFFSET %d`, strings.Join(fields, ", b."), joinFields, b.TableName(), joins, conds, s.Limit, s.Offset), args
}

func (r *bannerRepo) Search(ctx context.Context, s *SearchBannersOpts) ([]*BannerData, error) {
	banners := make([]*BannerData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&banners)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return banners, nil
}
