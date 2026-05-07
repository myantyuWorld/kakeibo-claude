---
paths:
  - "docs/api/**"
---

# API 仕様書規約

## OpenAPI バージョン

OpenAPI 3.0.3 を使用する。

## ファイル構成

```
docs/api/
├── openapi.yml              # ルート（$ref で各ドメインのパスを参照）
├── paths/
│   └── {domain}.yml         # ドメインごとのパス定義 + ドメイン固有の schemas
└── components/
    └── errors.yml           # 共通エラーレスポンス（responses + schemas）
```

- `openapi.yml` は各ドメインの paths を `$ref` で束ねる
- 新規ドメイン追加時は `openapi.yml` の paths・tags にエントリを追加する

## スキーマの配置ルール

| スキーマの種類 | 配置場所 |
|--------------|---------|
| 複数ドメインで使う共通型（DateTime, ErrorResponse 等） | `components/errors.yml` または共通コンポーネントファイル |
| そのドメインでのみ使うリクエスト/レスポンス型 | `paths/{domain}.yml` 内の `components.schemas` |

## パス規約

- ALWAYS パスセグメントは **kebab-case** のみ（スネークケース・キャメルケースは禁止）
- パスパラメータで操作対象が明らかな場合、リソース名は含めない

```yaml
# OK
/products/{id}
/products/{id}/stock-entries

# NG
/products/{productId}
/products/{product_id}
/products/{id}/stockEntries
```

## パラメータ・フィールド命名

- リクエストボディ、レスポンスボディ、クエリパラメータのフィールド名は **lowerCamelCase** に統一する

```yaml
# OK
properties:
  firstName:
    type: string
  createdAt:
    type: string

# NG
properties:
  first_name:
    type: string
  created_at:
    type: string
```

## operationId

省略する。openapi-fetch はパスベースのアクセスのため不要。

## スキーマ命名規約

| 用途 | パターン | 例 |
|------|---------|-----|
| 作成リクエスト | `Create{Entity}Request` | `CreateProductRequest` |
| 更新リクエスト | `Update{Entity}Request` | `UpdateProductRequest` |
| 詳細レスポンス | `{Entity}Detail` | `ProductDetail` |
| 一覧の各要素 | `{Entity}Summary` | `ProductSummary` |

## 共通エラーレスポンス

エラーの返却は `docs/api/components/errors.yml` の共通スキーマ（responses）を使用すること。

## description

- フィールドの説明は日本語で記載する
- example を積極的に付与する

## 関連

- skill: `generate-api-spec` -- API仕様書生成の手順・テンプレート
