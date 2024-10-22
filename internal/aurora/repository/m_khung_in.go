package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type MKhungInRepo interface {
	Insert(ctx context.Context, e *model.MKhungIn) error
	Update(ctx context.Context, e *model.MKhungIn) error
	SoftDelete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*MKhungInData, error)
	Search(ctx context.Context, s *SearchMKhungInOpts) ([]*MKhungInData, error)
	Count(ctx context.Context, s *SearchMKhungInOpts) (*CountResult, error)
}

type sMKhungInRepo struct {
}

func NewMKhungInRepo() MKhungInRepo {
	return &sMKhungInRepo{}
}

func (r *sMKhungInRepo) Insert(ctx context.Context, e *model.MKhungIn) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *sMKhungInRepo) FindByID(ctx context.Context, id string) (*MKhungInData, error) {
	e := &MKhungInData{}
	err := cockroach.FindOne(ctx, e, "id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("sMKhungInRepo.cockroach.FindOne: %w", err)
	}

	return e, nil
}
func (r *sMKhungInRepo) Update(ctx context.Context, e *model.MKhungIn) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

func (r *sMKhungInRepo) SoftDelete(ctx context.Context, id string) error {
	sql := "UPDATE m_khung_in SET deleted_at = NOW() WHERE id = $1;"

	cmd, err := cockroach.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("m_khung_in cockroach.Exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("*sMKhungInRepo not found any records to delete")
	}

	return nil
}

// SearchMKhungInOpts all params is options
type SearchMKhungInOpts struct {
	IDs []string
	// todo add more search options
	Limit  int64
	Offset int64
	Sort   *Sort
}

func (s *SearchMKhungInOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.MKhungInFieldID)
	}
	// todo add more search options example:
	//if s.Name != "" {
	//	args = append(args, "%"+s.Name+"%")
	//	conds += fmt.Sprintf(" AND (b.%[2]s ILIKE $%[1]d OR b.%[3]s ILIKE $%[1]d)",
	//		len(args), model.MKhungInFieldName, model.MKhungInFieldCode)
	//}
	//if s.Code != "" {
	//	args = append(args, s.Code)
	//	conds += fmt.Sprintf(" AND b.%s ILIKE $%d", model.MKhungInFieldCode, len(args))
	//}

	b := &model.MKhungIn{}
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

type MKhungInData struct {
	*model.MKhungIn
}

func (r *sMKhungInRepo) Search(ctx context.Context, s *SearchMKhungInOpts) ([]*MKhungInData, error) {
	MKhungIn := make([]*MKhungInData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&MKhungIn)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return MKhungIn, nil
}

func (r *sMKhungInRepo) Count(ctx context.Context, s *SearchMKhungInOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("sMKhungInRepo.Count: %w", err)
	}

	return countResult, nil
}
