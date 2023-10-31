package cockroach

import (
	"fmt"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
)

type Rows struct {
	err     error
	pgxRows pgx.Rows
}

func (r *Rows) ScanOne(dst interface{}) error {
	if r.err != nil {
		return r.err
	}

	return pgxscan.ScanOne(dst, r.pgxRows)
}

func (r *Rows) ScanAll(dst interface{}) error {
	if r.err != nil {
		return r.err
	}

	return pgxscan.ScanAll(dst, r.pgxRows)
}

func (r *Rows) ScanFields(dst ...interface{}) error {
	if r.err != nil {
		return r.err
	}

	defer r.pgxRows.Close()

	if !r.pgxRows.Next() {
		err := r.pgxRows.Err()
		if err == nil {
			return pgx.ErrNoRows
		}

		return fmt.Errorf("rows.Err: %w", err)
	}

	if err := r.pgxRows.Scan(dst...); err != nil {
		return fmt.Errorf("rows.Scan: %w", err)
	}

	if err := r.pgxRows.Err(); err != nil {
		return fmt.Errorf("rows.Err: %w", err)
	}

	return nil
}
