package cockroach

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func NewConnectionPool(ctx context.Context, logger *zap.Logger, connectionURI string, debugMode bool) Ext {
	connPoolConfig, err := pgxpool.ParseConfig(connectionURI)
	if err != nil {
		logger.Panic("cannot read uri", zap.Error(err))
	}

	connPoolConfig.ConnConfig.Logger = zapadapter.NewLogger(logger)

	if debugMode {
		connPoolConfig.ConnConfig.LogLevel = pgx.LogLevelInfo
	}

	pool, err := pgxpool.ConnectConfig(ctx, connPoolConfig)
	if err != nil {
		logger.Panic("cannot create new connection pool to postgres", zap.Error(err))
	}

	return pool
}

type (
	dbKey          int
	queryExecerKey int
)

// ctx must inject db pool connection before by func ContextWithDB
// if not, it will panic
func Conn(ctx context.Context) queryExecer {
	if c, ok := ctx.(*gin.Context); ok {
		ctx = c.Request.Context()
	}

	v := ctx.Value(queryExecerKey(0))
	s, ok := v.(queryExecer)
	if !ok {
		panic("ctx not have db pool")
	}

	return s
}

func ContextWithDB(ctx context.Context, db Ext) context.Context {
	ctx = context.WithValue(ctx, dbKey(0), db)
	return context.WithValue(ctx, queryExecerKey(0), db)
}

func dbFromContext(ctx context.Context) Ext {
	if c, ok := ctx.(*gin.Context); ok {
		ctx = c.Request.Context()
	}

	v := ctx.Value(dbKey(0))
	s, _ := v.(Ext)
	return s
}

func Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return Conn(ctx).Query(ctx, sql, args...)
}

func QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return Conn(ctx).QueryRow(ctx, sql, args...)
}

func Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	return Conn(ctx).Exec(ctx, sql, args...)
}

type TxHandler = func(ctx context.Context) error

func ExecInTx(ctx context.Context, fn TxHandler) error {
	tx, err := dbFromContext(ctx).Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "db.Begin")
	}

	ctx = context.WithValue(ctx, queryExecerKey(0), tx)
	defer func() {
		if err == nil {
			// Ignore commit errors. The tx has already been committed by RELEASE.
			_ = tx.Commit(ctx)
		} else {
			// We always need to execute a Rollback() so sql.DB releases the
			// connection.
			_ = tx.Rollback(ctx)
		}
	}()
	// Specify that we intend to retry this txn in case of CockroachDB retryable
	// errors.
	if _, err = tx.Exec(ctx, "SAVEPOINT cockroach_restart"); err != nil {
		return err
	}

	for {
		released := false
		err = fn(ctx)
		if err == nil {
			// RELEASE acts like COMMIT in CockroachDB. We use it since it gives us an
			// opportunity to react to retryable errors, whereas tx.Commit() doesn't.
			released = true
			if _, err = tx.Exec(ctx, "RELEASE SAVEPOINT cockroach_restart"); err == nil {
				return nil
			}
		}
		// We got an error; let's see if it's a retryable one and, if so, restart.
		// We look for either:
		//  - the standard PG errcode SerializationFailureError:40001 or
		//  - the Cockroach extension errcode RetriableError:CR000. This extension
		//    has been removed server-side, but support for it has been left here for
		//    now to maintain backwards compatibility.
		code := errCode(err)
		if retryable := code == "CR000" || code == "40001"; !retryable {
			if released {
				err = newAmbiguousCommitError(err)
			}
			return err
		}

		if _, retryErr := tx.Exec(ctx, "ROLLBACK TO SAVEPOINT cockroach_restart"); retryErr != nil {
			return newTxnRestartError(retryErr, err)
		}
	}
}

func errCode(err error) string {
	switch t := errorCause(err).(type) {
	case *pq.Error:
		fmt.Println("---------------------------------------------------------------------------------------------", t.Code)
		return string(t.Code)

	case *pgconn.PgError:
		return t.Code

	default:
		return ""
	}
}

// ErrorCauser is the type implemented by an error that remembers its cause.
//
// ErrorCauser is intentionally equivalent to the causer interface used by
// the github.com/pkg/errors package.
type ErrorCauser interface {
	// Cause returns the proximate cause of this error.
	Cause() error
}

// UnwrappableError describes errors compatible with errors.Unwrap.
type UnwrappableError interface {
	// Unwrap returns the proximate cause of this error.
	Unwrap() error
}

// Unwrap is equivalent to errors.Unwrap. It's implemented here to maintain
// compatibility with Go versions before 1.13 (when the errors package was
// introduced).
// It returns the result of calling the Unwrap method on err, if err's type
// implements UnwrappableError.
// Otherwise, Unwrap returns nil.
func Unwrap(err error) error {
	u, ok := err.(UnwrappableError)
	if !ok {
		return nil
	}
	return u.Unwrap()
}

// errorCause returns the original cause of the error, if possible. An error has
// a proximate cause if it's type is compatible with Go's errors.Unwrap() (and
// also, for legacy reasons, if it implements ErrorCauser); the original cause
// is the bottom of the causal chain.
func errorCause(err error) error {
	// First handle errors implementing ErrorCauser.
	for err != nil {
		cause, ok := err.(ErrorCauser)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	// Then handle go1.13+ error wrapping.
	for {
		cause := Unwrap(err)
		if cause == nil {
			break
		}
		err = cause
	}
	return err
}

type txError struct {
	cause error
}

// Error implements the error interface.
func (e *txError) Error() string { return e.cause.Error() }

// Cause implements the ErrorCauser interface.
func (e *txError) Cause() error { return e.cause }

// AmbiguousCommitError represents an error that left a transaction in an
// ambiguous state: unclear if it committed or not.
type AmbiguousCommitError struct {
	txError
}

func newAmbiguousCommitError(err error) *AmbiguousCommitError {
	return &AmbiguousCommitError{txError{cause: err}}
}

// TxnRestartError represents an error when restarting a transaction. `cause` is
// the error from restarting the txn and `retryCause` is the original error which
// triggered the restart.
type TxnRestartError struct {
	txError
	retryCause error
	msg        string
}

func newTxnRestartError(err error, retryErr error) *TxnRestartError {
	const msgPattern = "restarting txn failed. ROLLBACK TO SAVEPOINT " +
		"encountered error: %s. Original error: %s."
	return &TxnRestartError{
		txError:    txError{cause: err},
		retryCause: retryErr,
		msg:        fmt.Sprintf(msgPattern, err, retryErr),
	}
}

// Error implements the error interface.
func (e *TxnRestartError) Error() string { return e.msg + "==========" }

// RetryCause returns the error that caused the transaction to be restarted.
func (e *TxnRestartError) RetryCause() error { return e.retryCause }

func String(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  s != "",
	}
}

func Time(s time.Time) sql.NullTime {
	return sql.NullTime{
		Time:  s,
		Valid: !s.IsZero(),
	}
}
func Float64(f float64) sql.NullFloat64 {
	return sql.NullFloat64{
		Float64: f,
		Valid:   f != 0,
	}
}

func IsErrNoRows(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}
