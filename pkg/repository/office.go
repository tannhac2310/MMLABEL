package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
)

type OfficeRepo interface {
	Insert(ctx context.Context, e *model.Office) error
	Update(ctx context.Context, e *model.Office) error
	SoftDelete(ctx context.Context, id string) error
	Search(ctx context.Context, s *SearchOfficesOpts) ([]*OfficeData, error)
	Count(ctx context.Context, s *SearchOfficesOpts) (*CountResult, error)
}
type OfficeData struct {
	*model.Office
	ProvinceName  string `db:"province_name"`
	DistrictName  string `db:"district_name"`
	CreatedByName string `db:"created_by_name"`
}
type officeRepo struct {
}

func (r *officeRepo) Count(ctx context.Context, s *SearchOfficesOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	fmt.Println(sql)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("lesson.Count: %w", err)
	}

	return countResult, nil
}

func NewOfficeRepo() OfficeRepo {
	return &officeRepo{}
}

func (r *officeRepo) Insert(ctx context.Context, e *model.Office) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *officeRepo) SoftDelete(ctx context.Context, id string) error {
	sql := `UPDATE offices
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

func (r *officeRepo) Update(ctx context.Context, e *model.Office) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

// SearchOfficesOpts all params is options
type SearchOfficesOpts struct {
	IDs    []string
	Name   string
	Search string
	Limit  int64
	Offset int64
}

func (s *SearchOfficesOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	o := &model.Office{}
	fields, _ := o.FieldMap()
	var moreFields []string
	conds := ""
	joins := ` JOIN users cb on cb.id = o.created_by `
	moreFields = append(moreFields, "cb.name as created_by_name")

	joins += fmt.Sprintf(" JOIN regions rp on rp.%s = o.%s ", model.RegionFieldID, model.OfficeFieldProvinceID)
	moreFields = append(moreFields, "rp.name as province_name")

	joins += fmt.Sprintf(" JOIN regions rd on rd.%s = o.%s ", model.RegionFieldID, model.OfficeFieldDistrictID)
	moreFields = append(moreFields, "rd.name as district_name")
	//u.name as student_name, u.email as student_email, u.avatar as student_avatar
	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND o.%s = ANY($1)", model.OfficeFieldID)
	}

	if s.Name != "" {
		args = append(args, "%"+s.Name+"%")
		conds += fmt.Sprintf(" AND o.%s ILIKE $%d", model.OfficeFieldName, len(args))
	}
	if s.Search != "" {
		args = append(args, "%"+s.Search+"%")
		conds += fmt.Sprintf(" AND o.%s ILIKE $%d", model.OfficeFieldName, len(args))
	}

	if isCount {
		return fmt.Sprintf(`SELECT count(*) as cnt
		FROM %s AS o %s
		WHERE TRUE %s AND o.deleted_at IS NULL`, o.TableName(), joins, conds), args
	}
	return fmt.Sprintf(`SELECT o.%s,%s
		FROM %s AS o %s
		WHERE TRUE %s AND o.deleted_at IS NULL
		ORDER BY o.id DESC
		LIMIT %d
		OFFSET %d`, strings.Join(fields, ", o."), strings.Join(moreFields, ", "), o.TableName(), joins, conds, s.Limit, s.Offset), args
}

func (r *officeRepo) Search(ctx context.Context, s *SearchOfficesOpts) ([]*OfficeData, error) {
	offices := make([]*OfficeData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&offices)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return offices, nil
}
