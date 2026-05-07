package template

import "context"

// Creator: 他ドメインからの書き込み操作用のインターフェース実装
// Resolver（読み取り専用）とは異なり、他ドメインのイベントに応じてデータを作成する
// 例: 注文ドメインが注文作成時に、配送ドメインの配送レコードを作成する
// 利用側が自パッケージ内にインターフェースを定義し、この構造体が暗黙的に満たす
//
// Creator は service を経由してドメインロジック（バリデーション等）を共有する
// API と Creator で同じ service メソッドを使うため、入力は request/response ではなく DTO（CreateInput 等）を定義する
type Creator struct {
	svc Service
}

func NewCreator(svc Service) *Creator {
	return &Creator{svc: svc}
}

// Creator の利用方針
// - ほとんどのケースでは発生しないはず。大量に作成されるなら設計を見直すこと
// - 同期処理である必要がある場合のみ作成する
// - メール通知やプッシュ通知といった副作用は Creator ではなくイベント駆動で対応する
func (c *Creator) Create(ctx context.Context, input CreateInput) error {
	return c.svc.create(ctx, input)
}
