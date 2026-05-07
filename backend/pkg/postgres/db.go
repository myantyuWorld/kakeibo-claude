package postgres

import (
	"context"
	"database/sql"
	"embed"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
)

type DB struct {
	db *sqlx.DB
}

func NewDB(dsn string) (*DB, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return &DB{db: db}, nil
}

func (d *DB) ext(ctx context.Context) sqlx.ExtContext {
	if tx, ok := ctx.Value(txKey).(*sqlx.Tx); ok {
		return tx
	}
	return d.db
}

func (d *DB) beginTx(ctx context.Context) (*sqlx.Tx, error) {
	return d.db.BeginTxx(ctx, nil)
}

func (d *DB) Close() error {
	return d.db.Close()
}

func (d *DB) Migrate(fs embed.FS) error {
	source := &migrate.EmbedFileSystemMigrationSource{FileSystem: fs, Root: "."}
	_, err := migrate.Exec(d.db.DB, "postgres", source, migrate.Up)
	return err
}

func (d *DB) GetContext(ctx context.Context, dest any, query string, args ...any) error {
	return sqlx.GetContext(ctx, d.ext(ctx), dest, query, args...)
}

func (d *DB) SelectContext(ctx context.Context, dest any, query string, args ...any) error {
	return sqlx.SelectContext(ctx, d.ext(ctx), dest, query, args...)
}

func (d *DB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return d.ext(ctx).ExecContext(ctx, query, args...)
}

func (d *DB) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return d.ext(ctx).QueryContext(ctx, query, args...)
}

func (d *DB) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	if tx, ok := ctx.Value(txKey).(*sqlx.Tx); ok {
		return tx.QueryRowContext(ctx, query, args...)
	}
	return d.db.QueryRowContext(ctx, query, args...)
}

func (d *DB) NamedExecContext(ctx context.Context, query string, arg any) (sql.Result, error) {
	if tx, ok := ctx.Value(txKey).(*sqlx.Tx); ok {
		return tx.NamedExecContext(ctx, query, arg)
	}
	return d.db.NamedExecContext(ctx, query, arg)
}
