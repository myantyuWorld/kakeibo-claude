---
name: generate-er
description: ドメインモデルからER図（DBML）を生成する。「ER図を作って」「テーブル設計して」「DBMLを生成して」など、データベーススキーマの設計を依頼されたときに使用する。ドメイン名を引数で受け取る（例: /generate-er customer）。
argument-hint: "{ドメイン名}"
---

# ER図（DBML）生成 Skill

ドメインモデルの設計ドキュメントを入力に、DBML 形式の ER 図を生成する。

## When to Activate

- 「ER図を作って」「テーブル設計して」「DBMLを生成して」と依頼されたとき
- 新規ドメインの設計ドキュメントが揃い、DB スキーマ定義が必要になったとき
- 既存ドメインにテーブルを追加・変更するとき

## 入力

- ドメイン名（例: customer）
- 指定がない場合: AskUserQuestion で対象を選ばせる

## 手順

### 1. 設計ドキュメントの読み込み

以下を読み込む:

- `docs/design/{domain}/ドメインモデル.md` — エンティティ・値オブジェクト・集約構造
- `docs/design/{domain}/業務ルール.md` — バリデーション・ライフサイクル・文字数制約
- `.claude/rules/backend/database-design.md` — DDD→RDB マッピングルール・アンチパターンチェックリスト

### 2. DDD → RDB マッピング

ドメインモデルの各要素を RDB テーブル・カラムに変換する。変換は `database-design.md` のルールに従う。

#### マッピング判定フロー

```
ドメインモデルの要素を1つずつ確認:

Entity（集約ルート）か？
  → YES → テーブルを作成

ValueObject か？
  → 親 Entity と 1:1 か？
    → YES → 親テーブルにプレフィックス付きカラムとして埋め込む
    → NO（1:N）→ 子テーブルを作成、親への FK を付与

Enum か？
  → カラム（varchar）で表現。テーブルにしない

他集約への参照か？
  → UUID の FK カラムのみ追加
```

```dbml
// PASS: ValueObject（1:1）を親テーブルに埋め込む
Table customers {
  id uuid [pk]
  name varchar [not null]
  address_prefecture varchar [not null]
  address_city varchar [not null]
  address_street varchar [not null]
  status varchar [not null, note: 'Active / Deleted']
}

// FAIL: ValueObject を別テーブルに分離する
Table customers {
  id uuid [pk]
  name varchar [not null]
  address_id uuid [ref: - addresses.id]
}
Table addresses {
  id uuid [pk]
  prefecture varchar [not null]
  city varchar [not null]
  street varchar [not null]
}
```

```dbml
// PASS: Enum はカラムで表現
Table customers {
  id uuid [pk]
  status varchar [not null, note: 'Active / Deleted']
}

// FAIL: Enum をマスタテーブルにする
Table customer_statuses {
  id integer [pk]
  name varchar [not null]
}
```

```dbml
// PASS: 他集約への参照は UUID の FK のみ
Table orders {
  id uuid [pk]
  customer_id uuid [not null, ref: > customers.id]
}

// FAIL: 他集約のデータを非正規化して持つ
Table orders {
  id uuid [pk]
  customer_id uuid [not null]
  customer_name varchar [not null]
}
```

### 3. カラム定義の詳細化

業務ルール.md の制約をカラム定義に反映する。

#### 型マッピング

| ドメインモデルの型 | DBML の型 | 備考 |
|-----------------|----------|------|
| UUID | uuid | PK・FK に使用 |
| string | varchar | 文字数制約がある場合は note に記載 |
| int / integer | integer | |
| 金額・率 | numeric | FLOAT は使わない |
| 日時 | timestamptz | タイムゾーン付き |
| bool | boolean | |
| Enum | varchar | note に取りうる値を列挙 |

#### 制約の反映

```dbml
// PASS: 業務ルールの制約をカラム定義に反映
Table customers {
  id uuid [pk]
  name varchar [not null, note: '1〜100文字']
  email varchar [not null, unique, note: 'メールアドレス形式']
  status varchar [not null, note: 'Active / Deleted', default: 'Active']
  created_at timestamptz [not null, default: `now()`]
  updated_at timestamptz [not null, default: `now()`]
}
```

#### 共通カラム

すべてのテーブルに以下を付与する:

- `created_at timestamptz [not null, default: \`now()\`]`
- `updated_at timestamptz [not null, default: \`now()\`]`

### 4. リレーション定義

DBML の Ref 構文でリレーションを定義する。

```dbml
// 1:N — 顧客は複数の注文を持つ
Ref: orders.customer_id > customers.id

// N:N — 交差テーブル経由
Table order_items {
  order_id uuid [not null, ref: > orders.id]
  product_id uuid [not null, ref: > products.id]

  indexes {
    (order_id, product_id) [pk]
  }
}
```

### 5. アンチパターンチェック

生成した DBML を `database-design.md` のチェックリストで検証する。
チェック結果をユーザーに報告する。

```
## アンチパターンチェック結果

| チェック項目 | 結果 | 備考 |
|------------|------|------|
| ジェイウォーク | OK | カンマ区切りカラムなし |
| IDリクワイアド | OK | 交差テーブルは複合PK |
| キーレスエントリ | OK | 全FK制約あり |
| EAV | OK | 汎用カラムなし |
| ポリモーフィック関連 | OK | type+id パターンなし |
| マルチカラムアトリビュート | OK | 横持ちカラムなし |
| ラウンディングエラー | OK | 金額は numeric |
| サーティワンフレーバー | OK | Enum はカラムで表現 |
| フィア・オブ・ジ・アンノウン | OK | 必須カラムに NOT NULL |
```

### 6. DBML ファイルの出力

`docs/design/{domain}/er.dbml` に出力する。

```bash
# 出力先
docs/design/{domain}/er.dbml
```

### 7. ドキュメントの更新

`docs/design/{domain}/index.md` の関連ドキュメントに ER 図へのリンクを追加する。

```markdown
## 関連ドキュメント

- [ドメインモデル](ドメインモデル.md) - エンティティ・集約の構造と関係
- [業務ルール](業務ルール.md) - ビジネスルール・制約・受入基準
- [ER図](er.dbml) - データベーススキーマ定義（DBML）
```

## DBML 構文リファレンス

```dbml
// テーブル定義
Table table_name {
  column_name type [constraints]
}

// 制約
// pk       — 主キー
// not null — NOT NULL
// unique   — ユニーク
// default: value — デフォルト値（SQL式は バッククォートで囲む）
// ref: > other_table.id — 外部キー（インライン）
// note: 'メモ' — カラムの説明

// 複合主キー・複合インデックス
Table table_name {
  col_a uuid [not null]
  col_b uuid [not null]

  indexes {
    (col_a, col_b) [pk]
  }
}

// リレーション（別定義）
Ref: table_a.col > table_b.col   // many-to-one
Ref: table_a.col - table_b.col   // one-to-one
Ref: table_a.col <> table_b.col  // many-to-many

// テーブルグループ（集約単位でグルーピング）
TableGroup customer_aggregate {
  customers
  customer_contacts
}

// Enum（DBML上の定義。実テーブルでは varchar カラムに note で列挙する）
// → DBML の Enum 構文は使わない。カラムの note に記載する
```

## 命名規約

| 対象 | 規約 | 例 |
|------|------|-----|
| テーブル名 | snake_case・複数形 | `customers`, `order_items` |
| カラム名 | snake_case | `first_name`, `created_at` |
| FK カラム | `{参照先テーブルの単数形}_id` | `customer_id` |
| 交差テーブル | `{テーブルA}_{テーブルB}`（アルファベット順） | `order_products` |
| ValueObject 埋め込み | `{VO名}_{フィールド名}` | `address_city` |

命名はドメインモデル.md の対訳表（ユビキタス言語）に従う。

## 成功指標

- ドメインモデル.md の全 Entity・ValueObject が DBML に反映されている
- 業務ルール.md の文字数制約・NOT NULL が note / 制約に反映されている
- `database-design.md` のアンチパターンチェックが全項目 OK
- 対訳表の英語名がテーブル名・カラム名に使われている

---

**Remember**: ドメインモデルの構造を忠実にテーブルに変換する。ValueObject は埋め込み、Enum はカラム。テーブルを増やしすぎない。
