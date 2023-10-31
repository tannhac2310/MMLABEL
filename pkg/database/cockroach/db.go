package cockroach

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type Entity interface {
	TableName() string
	FieldMap() (fields []string, values []interface{})
}

// queryer is an interface for Query
type queryer interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

// execer is an interface for Exec
type execer interface {
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
}

type queryExecer interface {
	queryer
	execer
}

// TxStarter is an interface to deal with transaction
type TxStarter interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

// TxController is an interface to deal with transaction
type TxController interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

// Ext is a union interface which can bind, query, and exec
type Ext interface {
	queryer
	execer
	TxStarter
}

// Txer is a interface for Tx
type Txer interface {
	queryer
	execer
	TxController
}
