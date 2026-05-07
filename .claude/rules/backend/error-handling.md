---
paths:
  - "backend/**/*.go"
---

# エラーハンドリング

## 基本方針

- NEVER 標準の `errors` パッケージを使わない
- ALWAYS `pkg/errors` を使う（スタックトレースのため）
- NEVER エラーを `_` で握りつぶさない

## pkg/errors の使い方

### エラー生成

| 関数 | 用途 |
|------|------|
| `errors.New(err)` | 内部エラー（予期しないエラー） |
| `errors.NewInvalid(err)` | バリデーションエラー |
| `errors.NewNotFound(err)` | リソースが見つからない |
| `errors.NewUnauthorized(err)` | 認証エラー |
| `errors.NewForbidden(err)` | 権限エラー |

### WithMessage

デフォルトメッセージで十分な場合はそのまま返す。追加情報が必要な場合のみ `WithMessage` を使う。

- ALWAYS メッセージは日本語とする

```go
// デフォルトメッセージで十分な場合
return errors.NewNotFound(err)

// 具体的な情報が必要な場合
return errors.NewInvalid().WithMessage("名前は100文字以内で入力してください")
```

### メッセージのセキュリティ

NEVER 攻撃者の材料になるメッセージを返さない。システムの内部状態やルールの詳細を推測できる情報は避ける。

```go
// NG: アカウントの存在が推測できる
errors.NewInvalid().WithMessage("このメールアドレスは既に使用されています")

// NG: パスワードポリシーの詳細が漏れる
errors.NewInvalid().WithMessage("パスワードは0-9、A-Z、a-z、記号を組み合わせてください")

// OK: デフォルトメッセージ（「入力内容に誤りがあります」）で返す
errors.NewInvalid()
```

### エラー判定

```go
if errors.IsNotFound(err) { ... }
if errors.IsInvalid(err) { ... }
```

## ErrCode と HTTP ステータスの対応

`pkg/api/error_handler.go` で自動変換される。

| ErrCode | HTTP Status | デフォルトメッセージ |
|---------|-------------|-------------------|
| `ErrInvalid` | 400 Bad Request | 入力内容に誤りがあります |
| `ErrUnauthorized` | 401 Unauthorized | 認証が必要です |
| `ErrForbidden` | 403 Forbidden | アクセス権限がありません |
| `ErrNotFound` | 404 Not Found | リソースが見つかりません |
| `ErrInternal` | 500 Internal Server Error | 予期しないエラーが発生しました |

必要に応じて ErrCode を追加してよい。追加時は `pkg/errors/errors.go` の定数・デフォルトメッセージ・判定関数と、`pkg/api/error_handler.go` の HTTP ステータスマッピングを併せて更新すること。

## レイヤーごとの使い分け

- **ドメイン層**: ビジネスルール違反は `NewInvalid().WithMessage(...)` で具体的なメッセージを返す
- **リポジトリ層**: `sql.ErrNoRows` は `NewNotFound(err)` に変換、その他の DB エラーは `New(err)` でラップ
- **ハンドラ層**: エラーをそのまま return する（error_handler が自動変換）

## 関連

- **go-reviewer** agent -- コードレビュー時
- **security-reviewer** agent -- エラーレスポンスの情報漏洩チェック時
