package template

import (
	"backend/pkg/errors"
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// DB行構造体
// - unexported。リポジトリ内でのみ使用する
// - db タグで DB カラム名とマッピングする
type entityRow struct {
	ID     uuid.UUID `db:"id"`
	Name   string    `db:"name"`
	Field1 string    `db:"field1"`
	Field2 string    `db:"field2"`
	Status string    `db:"status"`
}

// toModel: DB行構造体 → ドメインモデルへの変換
func (r *entityRow) toModel() *entity {
	return &entity{
		id:   r.ID,
		name: r.Name,
		vo: valueObject{
			field1: r.Field1,
			field2: r.Field2,
		},
		status: status(r.Status),
	}
}

// リポジトリ実装
// - Repository インターフェース（service.go で定義）を暗黙的に満たす
// - sql.ErrNoRows は NewNotFound に変換、その他の DB エラーは New でラップする
type repository struct {
	db *sqlx.DB
}

func newRepository(db *sqlx.DB) *repository {
	return &repository{db: db}
}

func (r *repository) save(ctx context.Context, e *entity) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO entities (id, name, field1, field2, status) VALUES ($1, $2, $3, $4, $5)",
		e.id, e.name, e.vo.field1, e.vo.field2, e.status,
	)
	if err != nil {
		return errors.New(err)
	}
	return nil
}

func (r *repository) findByID(ctx context.Context, id uuid.UUID) (*entity, error) {
	var row entityRow
	if err := r.db.GetContext(ctx, &row, "SELECT * FROM entities WHERE id = $1", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFound(err)
		}
		return nil, errors.New(err)
	}
	return row.toModel(), nil
}
