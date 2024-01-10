package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type DeviceRepo interface {
	Insert(ctx context.Context, e *model.Device) error
	Update(ctx context.Context, e *model.Device) error
	SoftDelete(ctx context.Context, id string) error
	Search(ctx context.Context, s *SearchDevicesOpts) ([]*DeviceData, error)
	Count(ctx context.Context, s *SearchDevicesOpts) (*CountResult, error)
}

type devicesRepo struct {
}

func NewDeviceRepo() DeviceRepo {
	return &devicesRepo{}
}

func (r *devicesRepo) Insert(ctx context.Context, e *model.Device) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *devicesRepo) Update(ctx context.Context, e *model.Device) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

func (r *devicesRepo) SoftDelete(ctx context.Context, id string) error {
	sql := `UPDATE devices
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

// SearchDevicesOpts all params is options
type SearchDevicesOpts struct {
	IDs    []string
	Name   string
	Step   string
	Steps  []string
	Code   string
	Limit  int64
	Offset int64
	Sort   *Sort
}

func (s *SearchDevicesOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.DeviceFieldID)
	}
	if s.Name != "" {
		args = append(args, "%"+s.Name+"%")
		conds += fmt.Sprintf(" AND (b.%[2]s ILIKE $%[1]d OR b.%[3]s ILIKE $%[1]d)",
			len(args), model.DeviceFieldName, model.DeviceFieldCode)
	}
	if s.Step != "" {
		args = append(args, "%"+s.Step+"%")
		conds += fmt.Sprintf(" AND b.%s ILIKE $%d", model.DeviceFieldStep, len(args))
	}
	if len(s.Steps) > 0 {
		args = append(args, s.Steps)
		// condition like in array steps
		// compare each step in steps like value of model.DeviceFieldStep
		conds += fmt.Sprintf(" AND b.%s = ANY($%d)", model.DeviceFieldStep, len(args))
	}
	if s.Code != "" {
		args = append(args, s.Code)
		conds += fmt.Sprintf(" AND b.%s ILIKE $%d", model.DeviceFieldCode, len(args))
	}

	b := &model.Device{}
	fields, _ := b.FieldMap()
	if isCount {
		return fmt.Sprintf(`SELECT count(*) as cnt
		FROM %s AS b %s
		WHERE TRUE %s AND b.deleted_at IS NULL`, b.TableName(), joins, conds), args
	}

	order := " ORDER BY b.id DESC "
	if s.Sort != nil {
		order = fmt.Sprintf(" ORDER BY b.%s %s", s.Sort.By, s.Sort.Order)
	}
	return fmt.Sprintf(`SELECT b.%s
		FROM %s AS b %s
		WHERE TRUE %s AND b.deleted_at IS NULL
		%s
		LIMIT %d
		OFFSET %d`, strings.Join(fields, ", b."), b.TableName(), joins, conds, order, s.Limit, s.Offset), args
}

type DeviceData struct {
	*model.Device
}

func (r *devicesRepo) Search(ctx context.Context, s *SearchDevicesOpts) ([]*DeviceData, error) {
	devices := make([]*DeviceData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&devices)
	if err != nil {
		return nil, fmt.Errorf("deviceRepo.cockroach.Select: %w", err)
	}

	return devices, nil
}

func (r *devicesRepo) Count(ctx context.Context, s *SearchDevicesOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("chat.Count: %w", err)
	}

	return countResult, nil
}
