package template

import (
	"backend/pkg/errors"

	"github.com/google/uuid"
)

// ステータス列挙
// - ドメインモデル.md の状態定義に対応させること
type status string

const (
	active  status = "active"
	deleted status = "deleted"
)

// エンティティ
// - ドメインモデル.md の集約ルートに対応する
// - フィールドは全て unexported とする
// - 特定のユースケースや UI の都合でフィールドを追加しない
type entity struct {
	id     uuid.UUID
	name   string
	vo     valueObject
	status status
}

// コンストラクタ
// - バリデーションを含める
// - 不正な状態のエンティティが生成されないことを保証する
func newEntity(name string, vo valueObject) (*entity, error) {
	if name == "" {
		return nil, errors.NewInvalid().WithMessage("名前は必須です")
	}
	return &entity{
		id:     uuid.New(),
		name:   name,
		vo:     vo,
		status: active,
	}, nil
}

// ドメインロジック
// - ビジネスルールの実装はエンティティのメソッドとして定義する
// - 状態遷移の制約もここで守る
func (e *entity) delete() error {
	if e.status == deleted {
		return errors.NewInvalid().WithMessage("削除済みのリソースは削除できません")
	}
	e.status = deleted
	return nil
}

// 値オブジェクト
// - 不変。値レシーバを使う
// - 単独で永続化・取得しない
type valueObject struct {
	field1 string
	field2 string
}
