# コーディング共通規約

言語固有のルールは `frontend/` または `backend/` のルールを優先すること。

## コメント

NEVER「何をしているか」を説明するコメントを書かない。コードを読めばわかる。
ALWAYS コメントを書くなら「なぜそうしているか（Why）」を書く。

```go
// NG: 何をしているかの説明
// Validator は Echo の echo.Validator インターフェースを満たす。
type validator struct { ... }

// NG: 処理の逐次説明
// ユーザーを取得する
user, err := repo.FindByID(ctx, id)

// OK: 業務上の理由
// 退会済みユーザーは論理削除のため、検索対象から除外する
query = query.Where("deleted_at IS NULL")

// OK: 技術的な制約の説明
// Echo の CustomValidator はポインタレシーバでないと登録できない
func (v *validator) Validate(i any) error { ... }
```

- テンプレート（`internal/_template/`）のコメントは設計判断の説明として許容する
- TODO コメントは理由と対応時期を併記する（例: `// TODO: v2 で認証方式を変更する`）
