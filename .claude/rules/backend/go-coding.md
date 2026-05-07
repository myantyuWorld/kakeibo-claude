---
paths:
  - "backend/**/*.go"
---

# Go コーディング規約

## フォーマット

- ALWAYS gofmt を適用する（goimports でも可）
- ALWAYS golangci-lint でLintを通すこと

## 命名

- MixedCaps（exported）/ mixedCaps（unexported）。アンダースコアは使わない
- インターフェース名は -er 接尾辞（Reader, Writer, Formatter 等）
- パッケージ名は小文字の単一語。util, common, base のような汎用名は避ける
- ゲッターに Get プレフィックスは付けない（`obj.Name()` であって `obj.GetName()` ではない）
- 頭字語は全大文字を保つ（`HTTPServer`, `userID`）
- パッケージ名と型名・定数名の重複を避ける。呼び出し側で `パッケージ名.型名` となるため、パッケージ名で文脈が明らかな情報を繰り返さない

```go
package order

type Status string   // OK: order.Status
type OrderStatus string // NG: order.OrderStatus（order が重複）

const (
    Draft  Status = "draft"   // OK: order.Draft
    OrderStatusDraft Status = "draft" // NG: order.OrderStatusDraft
)
```

## エラーハンドリング

- NEVER エラーを `_` で握りつぶさない
- NEVER 標準の `errors` パッケージを使わない
- ALWAYS `pkg/errors` を使う（スタックトレースのため）
- デフォルトメッセージで十分な場合はそのまま返す
- 追加情報が必要な場合は `WithMessage` を使う。メッセージは必ず日本語とする

## ポインタとレシーバ

- レシーバが状態を変更する場合はポインタレシーバを使う
- レシーバが大きな構造体の場合はポインタレシーバでコピーコストを避ける
- 一貫性: 同じ型のメソッドはポインタ/値レシーバを混在させない
- 値オブジェクト（Address 等の不変型）は値レシーバで十分

## 構造体とインターフェース

- 構造体はゼロ値が有用な状態になるよう設計する
- インターフェースは利用側（consumer）で定義する。実装側で定義しない
- 必要になるまでインターフェースを作らない（YAGNI）
- 同一パッケージ内で具象型を1つラップするだけのインターフェースは作らない
- ただしドメイン境界を跨ぐ場合は、利用側が必要最小限のインターフェースを定義する

## 並行処理

- goroutine の起動元がそのライフサイクルに責任を持つ
- context.Context は関数の第1引数として伝搬する
- channel よりも sync パッケージの方がシンプルな場合はそちらを使う
- 共有メモリの保護には sync.Mutex を使い、スコープを最小に保つ

## import

- 標準ライブラリ / 外部パッケージ / プロジェクト内パッケージ の3グループに分ける
- ブランクインポートは副作用が必要な場合のみ、コメントで理由を明記する

## その他

- NEVER init() を使わない。明示的な初期化関数を使う
- NEVER パニックを使わない。エラーを返す
- 型アサーションでは常に ok パターンを使う（`v, ok := x.(T)`）
- defer は関数の先頭近くに書き、リソースの確保と解放を近接させる

## 関連

- **go-reviewer** agent -- コードレビュー時
