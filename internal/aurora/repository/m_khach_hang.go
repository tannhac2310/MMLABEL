package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type MKhachHangRepo interface {
	Insert(ctx context.Context, e *model.MKhachHang) error
	Update(ctx context.Context, e *model.MKhachHang) error
	SoftDelete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*MKhachHangData, error)
	Search(ctx context.Context, s *SearchMKhachHangOpts) ([]*MKhachHangData, error)
	Count(ctx context.Context, s *SearchMKhachHangOpts) (*CountResult, error)
}

type sMKhachHangRepo struct {
}

func NewMKhachHangRepo() MKhachHangRepo {
	return &sMKhachHangRepo{}
}

func (r *sMKhachHangRepo) Insert(ctx context.Context, e *model.MKhachHang) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *sMKhachHangRepo) FindByID(ctx context.Context, id string) (*MKhachHangData, error) {
	e := &MKhachHangData{}
	err := cockroach.FindOne(ctx, e, "id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("sMKhachHangRepo.cockroach.FindOne: %w", err)
	}

	return e, nil
}
func (r *sMKhachHangRepo) Update(ctx context.Context, e *model.MKhachHang) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

func (r *sMKhachHangRepo) SoftDelete(ctx context.Context, id string) error {
	sql := "UPDATE m_khach_hang SET deleted_at = NOW() WHERE id = $1;"

	cmd, err := cockroach.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("m_khach_hang cockroach.Exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("*sMKhachHangRepo not found any records to delete")
	}

	return nil
}

// SearchMKhachHangOpts all params is options
type SearchMKhachHangOpts struct {
	IDs []string
	// todo add more search options
	Limit  int64
	Offset int64
	Sort   *Sort
}

func (s *SearchMKhachHangOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.MKhachHangFieldID)
	}
	// todo add more search options example:
	//if s.Name != "" {
	//	args = append(args, "%"+s.Name+"%")
	//	conds += fmt.Sprintf(" AND (b.%[2]s ILIKE $%[1]d OR b.%[3]s ILIKE $%[1]d)",
	//		len(args), model.MKhachHangFieldName, model.MKhachHangFieldCode)
	//}
	//if s.Code != "" {
	//	args = append(args, s.Code)
	//	conds += fmt.Sprintf(" AND b.%s ILIKE $%d", model.MKhachHangFieldCode, len(args))
	//}

	b := &model.MKhachHang{}
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

type MKhachHangData struct {
	*model.MKhachHang
}

func (r *sMKhachHangRepo) Search(ctx context.Context, s *SearchMKhachHangOpts) ([]*MKhachHangData, error) {
	MKhachHang := make([]*MKhachHangData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&MKhachHang)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return MKhachHang, nil
}

func (r *sMKhachHangRepo) Count(ctx context.Context, s *SearchMKhachHangOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("sMKhachHangRepo.Count: %w", err)
	}

	return countResult, nil
}
