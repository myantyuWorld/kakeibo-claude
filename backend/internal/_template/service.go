package template

import (
	"backend/pkg/errors"
	"backend/pkg/transaction"
	"context"

	"github.com/google/uuid"
)

// Repositoryインターフェース
// - 他ドメインから直接利用は想定していないため、型は exported、メソッドは unexportedとする
// - リポジトリはドメイン単位でデータの取得・更新を行うこと
// - 特定のフィールドだけ取得するようなメソッドの作成は禁止（アンチパターン）
type Repository interface {
	save(ctx context.Context, e *entity) error
	findByID(ctx context.Context, id uuid.UUID) (*entity, error)
}

// 他ドメインから取得する情報の構造体
// 実体は参照先ドメインの resolver.go で定義されている
type ExampleSummary struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// 取得系のAPIで他ドメインの情報が必要な場合に定義するinterface
// このinterfaceは他ドメイン側で暗黙的に満たされる（resolver.go参照）
type ExampleResolver interface {
	GetSummary(ctx context.Context, id uuid.UUID) (ExampleSummary, error)
}

type service struct {
	repo     Repository
	resolver ExampleResolver
	tx       transaction.Transactor
}

func newService(repo Repository, tx transaction.Transactor) *service {
	return &service{repo: repo, tx: tx}
}

// serviceメソッドの入出力ルール
// - API と service が 1対1 の場合は handler の request/response 構造体をそのまま受け渡す
// - Creator 等、複数から呼び出される場合は request/response とは別の DTO（CreateInput 等）を定義する
// - handler にドメインモデルを返却しない。response 構造体に詰め替えて返すこと

// CreateInput: API（handler）と他ドメイン（Creator）の両方から使われる入力 DTO
type CreateInput struct {
	Name   string
	Status string
}

func (s *service) create(ctx context.Context, input CreateInput) error {
	e, err := newEntity(input.Name, valueObject{
		field1: input.Status,
	})
	if err != nil {
		return err
	}
	return s.repo.save(ctx, e)
}

// resolver の使用例: 他ドメインの情報を取得してレスポンスに含める
func (s *service) getByID(ctx context.Context, id uuid.UUID) (*GetResponse, error) {
	e, err := s.repo.findByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if e.status == deleted {
		return nil, errors.NewNotFound()
	}

	// 他ドメインの情報を resolver 経由で取得する
	summary, err := s.resolver.GetSummary(ctx, e.id)
	if err != nil {
		return nil, err
	}

	return &GetResponse{
		ID:      e.id.String(),
		Name:    e.name,
		Status:  string(e.status),
		Example: summary,
	}, nil
}
