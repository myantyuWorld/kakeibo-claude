package template

import (
	"context"

	"github.com/google/uuid"
)

// 他ドメインから参照される読み取り専用のインターフェース実装を行う
// 利用側が自パッケージ内にインターフェースを定義し、この構造体が暗黙的に満たす
type Resolver struct {
	repo Repository
}

func NewResolver(repo Repository) *Resolver {
	return &Resolver{repo: repo}
}

// 自ドメインの外部公開用の構造体
// 他ドメインに必要最小限の情報だけ公開する
// 値オブジェクトはプリミティブ型に変換すること
type Summary struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func (r *Resolver) GetSummary(ctx context.Context, id uuid.UUID) (Summary, error) {
	e, err := r.repo.findByID(ctx, id)
	if err != nil {
		return Summary{}, err
	}
	return Summary{
		ID:   e.id,
		Name: e.name,
	}, nil
}
