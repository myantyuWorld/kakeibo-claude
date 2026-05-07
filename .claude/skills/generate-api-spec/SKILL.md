---
name: generate-api-spec
description: 設計ドキュメントからAPI仕様書（OpenAPI YAML）を生成する。「API仕様書を作って」「APIスペックを生成して」など、API仕様書の作成を依頼されたときに使用する。ドメイン名とユースケース名を引数で受け取る（例: /generate-api-spec product 商品登録）。ドメイン名のみの場合は全ユースケースを対象とする。
---

# API仕様書生成 Skill

設計ドキュメントと画面キャプチャから OpenAPI 仕様（YAML）を生成する。

## When to Activate

- 「API仕様書を作って」「APIスペックを生成して」と依頼されたとき
- 新規ユースケースの設計ドキュメントが揃い、API定義が必要になったとき
- 既存ドメインに新しいエンドポイントを追加するとき

## 実行手順

### 1. 対象の特定

- 引数: `{ドメイン名}` `{ユースケース名}`（例: `product 商品登録`）
- ユースケース名は業務ルール.md のセクション名に対応する
- ユースケース名が指定された場合: そのユースケースのAPIのみ生成する
- ドメイン名のみの場合: 全ユースケースを対象とする
- 指定がない場合: AskUserQuestion で対象を選ばせる

### 2. 設計ドキュメントの読み込み

以下を読み込む:

- `docs/design/{domain}/index.md` — ドメイン概要・ユースケース一覧
- `docs/design/{domain}/ドメインモデル.md` — エンティティ・値オブジェクト構造
- `docs/design/{domain}/業務ルール.md` — バリデーションルール・受入基準
- `.claude/rules/architecture.md` — ルーティング規約・エラーハンドリング

### 3. ユースケースの分類と生成方針

ユースケースを「更新系」「取得系」に分類し、それぞれ異なるアプローチで生成する。

#### 更新系（POST / PUT / DELETE）

ドメインモデルと業務ルールから自動生成できる。

- **リクエストボディ**: エンティティのフィールドから導出
  - コンストラクタの引数 → POST のリクエストボディ
  - Update メソッドの引数 → PUT のリクエストボディ
  - ID のみ → DELETE のパスパラメータ
- **バリデーション**: 業務ルールの制約をそのまま記載
- **エラーレスポンス**: 業務ルールの受入基準の異常系から導出
- **レスポンスボディ**: 作成・更新されたリソースを返す

#### 取得系（GET）

UIに依存するため、画面キャプチャまたはユーザーへのヒアリングで補完する。

**画面キャプチャが提供された場合:**

1. 画面キャプチャを読み取り、表示されている情報を列挙する
2. 各情報がどのエンティティ・フィールドに対応するか推論する
3. 関連ドメインのデータが含まれているか確認する（例: 商品一覧に店舗名が表示されている）
4. 推論結果を AskUserQuestion でユーザーに確認する

推論できる情報:
- 表示フィールド（何が画面に見えているか）
- 関連データの有無（他ドメインの情報が表示されているか）
- 一覧/詳細の区別

推論が難しい情報（AskUserQuestion で確認する）:
- ページネーション方式（オフセット / カーソル / 無限スクロール）
- ソート条件・デフォルト順
- フィルタ条件
- 検索機能の有無

**画面キャプチャがない場合:**

AskUserQuestion で以下をヒアリングする:

1. レスポンスに含めるフィールド
2. 関連ドメインのデータを含めるか
3. ページネーションの要否と方式
4. フィルタ・ソートの要否と条件

### 4. 既存のAPI仕様書構造の確認

生成前に既存の構造を読み込む:

- `docs/api/openapi.yml` — ルートファイル（全ドメインを束ねる）
- `docs/api/paths/{domain}.yml` — 対象ドメインのパス定義（あれば）
- `docs/api/components/errors.yml` — 共通エラー定義

### 5. OpenAPI YAML の生成

`docs/api/paths/{domain}.yml` に出力する。

- 既にファイルが存在する場合は、既存のパスに追記する（既存のエンドポイントは上書きしない）
- 新規ドメインの場合は `docs/api/openapi.yml` の paths・tags にも `$ref` エントリを追加する
- `.claude/rules/common/api-spec.md` の規約に従うこと
- 以下の YAML テンプレートをベースに生成する
- スキーマ名はドメインモデル.md の対訳表（ユビキタス言語）に従う

#### paths/{domain}.yml テンプレート

```yaml
# POST（作成）
/{domain}:
  post:
    summary: "{ユースケース名}"
    tags:
      - {Domain}
    requestBody:
      required: true
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/Create{Entity}Request'
    responses:
      "201":
        description: Created
      "400":
        $ref: '../components/errors.yml#/responses/InvalidError'
      "500":
        $ref: '../components/errors.yml#/responses/InternalServerError'

# GET（一覧）
/{domain}:
  get:
    summary: "{ユースケース名}"
    tags:
      - {Domain}
    responses:
      "200":
        description: OK
        content:
          application/json:
            schema:
              type: array
              items:
                $ref: '#/components/schemas/{Entity}Summary'
      "500":
        $ref: '../components/errors.yml#/responses/InternalServerError'

# GET（詳細）
/{domain}/{id}:
  get:
    summary: "{ユースケース名}"
    tags:
      - {Domain}
    parameters:
      - name: id
        in: path
        required: true
        schema:
          type: string
          format: uuid
    responses:
      "200":
        description: OK
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/{Entity}Detail'
      "404":
        $ref: '../components/errors.yml#/responses/NotFoundError'
      "500":
        $ref: '../components/errors.yml#/responses/InternalServerError'

components:
  schemas:
    Create{Entity}Request:
      type: object
      required:
        - name
      properties:
        name:
          type: string
          description: "{フィールド説明}"
          example: "キーボード"
    {Entity}Detail:
      type: object
      required:
        - id
        - name
      properties:
        id:
          type: string
          format: uuid
        name:
          type: string
    {Entity}Summary:
      type: object
      required:
        - id
        - name
      properties:
        id:
          type: string
          format: uuid
        name:
          type: string
```
- 既存の `docs/api/paths/` にあるファイルのスタイルに合わせる

### 6. 生成結果の確認

生成後、以下の対応表を出力する:

```
## ユースケース → API 対応表
| ユースケース | メソッド | パス | 生成方法 |
|------------|--------|------|---------|
| 商品登録 | POST | /api/products | 自動生成 |
| 商品一覧 | GET | /api/products | キャプチャ + 確認 |
| 商品詳細 | GET | /api/products/:id | キャプチャ + 確認 |
| 商品更新 | PUT | /api/products/:id | 自動生成 |
| 商品削除 | DELETE | /api/products/:id | 自動生成 |
| 在庫補充 | POST | /api/products/:id/restock | 自動生成 |
```

### 7. index.md の更新

生成完了後、`docs/design/{domain}/index.md` の関連ドキュメントに API 仕様書へのリンクを追加する。

```markdown
## 関連ドキュメント

- [ドメインモデル](ドメインモデル.md) - エンティティ・集約の構造と関係
- [業務ルール](業務ルール.md) - ビジネスルール・制約・受入基準
- [API仕様書](../../api/paths/{domain}.yml) - OpenAPI 仕様
```

---

**Remember**: 更新系はドメインモデルから自動導出、取得系はUIから逆算。
