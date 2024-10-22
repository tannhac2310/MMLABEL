package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type MKhuonBeRepo interface {
	Insert(ctx context.Context, e *model.MKhuonBe) error
	Update(ctx context.Context, e *model.MKhuonBe) error
	SoftDelete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*MKhuonBeData, error)
	Search(ctx context.Context, s *SearchMKhuonBeOpts) ([]*MKhuonBeData, error)
	Count(ctx context.Context, s *SearchMKhuonBeOpts) (*CountResult, error)
}

type sMKhuonBeRepo struct {
}

func NewMKhuonBeRepo() MKhuonBeRepo {
	return &sMKhuonBeRepo{}
}

func (r *sMKhuonBeRepo) Insert(ctx context.Context, e *model.MKhuonBe) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *sMKhuonBeRepo) FindByID(ctx context.Context, id string) (*MKhuonBeData, error) {
	e := &MKhuonBeData{}
	err := cockroach.FindOne(ctx, e, "id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("sMKhuonBeRepo.cockroach.FindOne: %w", err)
	}

	return e, nil
}
func (r *sMKhuonBeRepo) Update(ctx context.Context, e *model.MKhuonBe) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

func (r *sMKhuonBeRepo) SoftDelete(ctx context.Context, id string) error {
	sql := "UPDATE m_khuon_be SET deleted_at = NOW() WHERE id = $1;"

	cmd, err := cockroach.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("m_khuon_be cockroach.Exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("*sMKhuonBeRepo not found any records to delete")
	}

	return nil
}

// SearchMKhuonBeOpts all params is options
type SearchMKhuonBeOpts struct {
	IDs []string
	// todo add more search options
	Limit  int64
	Offset int64
	Sort   *Sort
}

func (s *SearchMKhuonBeOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.MKhuonBeFieldID)
	}
	// todo add more search options example:
	//if s.Name != "" {
	//	args = append(args, "%"+s.Name+"%")
	//	conds += fmt.Sprintf(" AND (b.%[2]s ILIKE $%[1]d OR b.%[3]s ILIKE $%[1]d)",
	//		len(args), model.MKhuonBeFieldName, model.MKhuonBeFieldCode)
	//}
	//if s.Code != "" {
	//	args = append(args, s.Code)
	//	conds += fmt.Sprintf(" AND b.%s ILIKE $%d", model.MKhuonBeFieldCode, len(args))
	//}

	b := &model.MKhuonBe{}
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

type MKhuonBeData struct {
	*model.MKhuonBe
}

func (r *sMKhuonBeRepo) Search(ctx context.Context, s *SearchMKhuonBeOpts) ([]*MKhuonBeData, error) {
	MKhuonBe := make([]*MKhuonBeData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&MKhuonBe)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return MKhuonBe, nil
}

func (r *sMKhuonBeRepo) Count(ctx context.Context, s *SearchMKhuonBeOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("sMKhuonBeRepo.Count: %w", err)
	}

	return countResult, nil
}
