package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type MPhimRepo interface {
	Insert(ctx context.Context, e *model.MPhim) error
	Update(ctx context.Context, e *model.MPhim) error
	SoftDelete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*MPhimData, error)
	Search(ctx context.Context, s *SearchMPhimOpts) ([]*MPhimData, error)
	Count(ctx context.Context, s *SearchMPhimOpts) (*CountResult, error)
}

type sMPhimRepo struct {
}

func NewMPhimRepo() MPhimRepo {
	return &sMPhimRepo{}
}

func (r *sMPhimRepo) Insert(ctx context.Context, e *model.MPhim) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *sMPhimRepo) FindByID(ctx context.Context, id string) (*MPhimData, error) {
	e := &MPhimData{}
	err := cockroach.FindOne(ctx, e, "id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("sMPhimRepo.cockroach.FindOne: %w", err)
	}

	return e, nil
}
func (r *sMPhimRepo) Update(ctx context.Context, e *model.MPhim) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

func (r *sMPhimRepo) SoftDelete(ctx context.Context, id string) error {
	sql := "UPDATE m_phim SET deleted_at = NOW() WHERE id = $1;"

	cmd, err := cockroach.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("m_phim cockroach.Exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("*sMPhimRepo not found any records to delete")
	}

	return nil
}

// SearchMPhimOpts all params is options
type SearchMPhimOpts struct {
	IDs []string
	// todo add more search options
	Limit  int64
	Offset int64
	Sort   *Sort
}

func (s *SearchMPhimOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.MPhimFieldID)
	}
	// todo add more search options example:
	//if s.Name != "" {
	//	args = append(args, "%"+s.Name+"%")
	//	conds += fmt.Sprintf(" AND (b.%[2]s ILIKE $%[1]d OR b.%[3]s ILIKE $%[1]d)",
	//		len(args), model.MPhimFieldName, model.MPhimFieldCode)
	//}
	//if s.Code != "" {
	//	args = append(args, s.Code)
	//	conds += fmt.Sprintf(" AND b.%s ILIKE $%d", model.MPhimFieldCode, len(args))
	//}

	b := &model.MPhim{}
	fields, _ := b.FieldMap()
	if isCount {
		return fmt.Sprintf("SELECT count(*) as cnt FROM %s AS b %s WHERE TRUE %s AND b.deleted_at IS NULL", b.TableName(), joins, conds), args
	}

	order := " ORDER BY b.id DESC "
	if s.Sort != nil {
		order = fmt.Sprintf(" ORDER BY b.%s %s", s.Sort.By, s.Sort.Order)
	}
	return fmt.Sprintf("SELECT b.%s FROM %s AS b %s WHERE TRUE %s AND b.deleted_at IS NULL %s LIMIT %d OFFSET %d", strings.Join(fields, ", b."), b.TableName(), joins, conds, order, s.Limit, s.Offset), args
}

type MPhimData struct {
	*model.MPhim
}

func (r *sMPhimRepo) Search(ctx context.Context, s *SearchMPhimOpts) ([]*MPhimData, error) {
	MPhim := make([]*MPhimData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&MPhim)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return MPhim, nil
}

func (r *sMPhimRepo) Count(ctx context.Context, s *SearchMPhimOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("sMPhimRepo.Count: %w", err)
	}

	return countResult, nil
}
