---
paths:
  - "docs/design/**/*.md"
  - "backend/**/*.go"
  - "backend/**/*.sql"
---

# データベース設計ガイドライン

## DDD → RDB マッピングルール

- ALWAYS Entity（集約ルート）はテーブルにする
- ALWAYS ValueObject（1:1）は同一テーブルのカラム群にする。プレフィックス付きで埋め込む（例: `address_city`, `address_zip`）
- NEVER ValueObject（1:1）を別テーブルに分離しない
- ALWAYS ValueObject（1:N）は子テーブルにする。親テーブルへのFK必須
- ALWAYS Enum はカラム（string / CHECK制約）で表現する
- NEVER Enum をマスタテーブルにしない。値の追加はマイグレーションで行う
- ALWAYS 他集約への参照は FK（IDカラムのみ）で表現する
- NEVER 他集約のデータを非正規化して持たない
- ALWAYS フロントエンドに公開するIDはUUIDにする。連番IDを外部に露出しない

## SQLアンチパターン チェックリスト

スキーマ設計時に以下を確認する。

- NEVER カンマ区切りで複数値を1カラムに入れない（ジェイウォーク）→ 交差テーブルを使う
- NEVER 全テーブルに機械的に `id` を付けない（IDリクワイアド）→ 自然キー・複合キーを検討する
- ALWAYS FK制約を設定する（キーレスエントリ）→ 参照整合性はDB側で担保する
- NEVER attr_name / attr_value のような汎用カラムを作らない（EAV）→ 具体的なカラム定義にする
- NEVER type + id で複数テーブルを参照しない（ポリモーフィック関連）→ 交差テーブルか共通親テーブルにする
- NEVER tag1, tag2, tag3 のような横持ちカラムを作らない（マルチカラムアトリビュート）→ 従属テーブルにする
- NEVER 金額・率に FLOAT を使わない（ラウンディングエラー）→ NUMERIC / DECIMAL を使う
- NEVER 許可値を CHECK 制約にベタ書きしない（サーティワンフレーバー）→ 変更頻度が高い場合は参照テーブルを検討する
- ALWAYS 業務上必須のカラムには NOT NULL を付与する（フィア・オブ・ジ・アンノウン）

## 関連

- skill: `generate-er` -- ER図生成時
- **go-reviewer** agent -- スキーマ変更レビュー時
