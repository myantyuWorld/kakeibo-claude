---
name: review-design
description: sudoモデリング成果物（drawio）をレビューし、仕様の曖昧さを指摘・付箋追記する。「設計をレビューして」「図をチェックして」「sudoの成果物を確認して」など、設計段階のdrawio図のレビューを依頼されたときに使用する。
argument-hint: "{図名}（例: ドメインモデル図）。省略時は全図を対象"
---

# review-design

sudoモデリングの drawio 成果物をレビューし、仕様の曖昧さや不整合を洗い出す。指摘箇所には drawio 上に付箋を追記する。

## When to Activate

- 「設計をレビューして」「図をチェックして」と依頼されたとき
- sudoモデリングの成果物が作成・更新された直後
- 設計フェーズで仕様の曖昧さを洗い出したいとき

## 設計思想

設計フェーズでの AI 介入は、実装フェーズでの手戻りコストを大幅に削減する。drawio は捨て資料（コンテキストマップとシステム関連図を除く）なので、付箋で汚すことに遠慮は不要。

## アーキテクチャ

```
Phase 1: 対象ファイルの特定・読み込み
  ↓
Phase 2: document-reviewer エージェントによるレビュー
  ↓
Phase 3: drawio に付箋を追記
  ↓
Phase 4: レビュー結果の出力
```

## 入力

- 図名（省略時は `docs/diagram/` 配下の全 drawio ファイルを対象）

## 実行フロー

### Phase 1: 対象ファイルの特定・読み込み

1. `docs/diagram/` 配下の drawio ファイルを検索する
2. 図名が指定された場合、該当する図に絞り込む
3. 以下の図を読み込む（存在するもののみ）

| 図 | 想定パス |
|----|---------|
| コンテキストマップ | `docs/diagram/コンテキストマップ.drawio` |
| システム関連図 | `docs/diagram/システム関連図.drawio` |
| ユースケース図 | `docs/diagram/ユースケース図.drawio` |
| ドメインモデル図 | `docs/diagram/ドメインモデル図.drawio` |
| オブジェクト図 | `docs/diagram/オブジェクト図.drawio` |

ファイルが見つからない場合は AskUserQuestion でパスを確認する。

### Phase 2: document-reviewer エージェントによるレビュー

document-reviewer エージェントを実行し、全図を横断的にレビューする。

エージェントへの入力:
- Phase 1 で読み込んだ drawio ファイルの内容（XML）
- 関連する設計ドキュメントの内容（あれば）

### Phase 3: drawio に付箋を追記

レビュー結果の CRITICAL / HIGH の指摘について、該当する drawio ファイルに付箋（sticky note）を追記する。

#### 付箋の XML テンプレート

```xml
<mxCell id="review-{連番}" value="{指摘内容}" style="shape=note;whiteSpace=wrap;html=1;backgroundOutline=1;fillColor=#FFF2CC;strokeColor=#D6B656;fontSize=11;align=left;verticalAlign=top;spacingTop=5;spacingLeft=5;spacingRight=5;" vertex="1" parent="1">
  <mxGeometry x="{x座標}" y="{y座標}" width="200" height="80" as="geometry" />
</mxCell>
```

#### 付箋の配置ルール

- 指摘対象の要素の近くに配置する
- 既存の要素と重ならないように座標を調整する
- 付箋の id は `review-` プレフィックスで統一する（再実行時に既存の付箋を識別・削除するため）

#### 付箋の内容フォーマット

```
[{重要度}] {指摘の要約}
```

### Phase 4: レビュー結果の出力

レビュー結果を以下の形式でテキスト出力する。

```markdown
## レビューサマリー

- 対象: {レビューした図の一覧}
- CRITICAL: {件数}件
- HIGH: {件数}件
- MEDIUM: {件数}件
- LOW: {件数}件
- 付箋追記: {追記した図のファイル名一覧}

## 指摘一覧

### CRITICAL
- [{図名}] {指摘内容}

### HIGH
- [{図名}] {指摘内容}

### MEDIUM
- [{図名}] {指摘内容}

### LOW
- [{図名}] {改善提案}
```

## エラーハンドリング

- **Phase 1 で drawio ファイルが見つからない**: AskUserQuestion でパスを確認。存在しない場合は中断し、先に図を作成するよう案内する
- **Phase 2 で drawio の XML が解析できない**: 対象ファイルをスキップし、解析できたファイルのみでレビューを続行する
- **Phase 3 で付箋の座標が決定できない**: 図の右下に固定配置する（x=既存要素の最大x + 50, y=0 から順に並べる）

## アンチパターン

- 実装の詳細（コード構造、API設計）に踏み込んだ指摘をする（設計フェーズの範囲を超える）
- 全ての指摘を付箋にする（MEDIUM / LOW はテキスト出力のみで十分）
- 図の構造や見た目（レイアウト、色使い）に対する指摘をする

## 成功指標

- CRITICAL / HIGH の指摘に対応する付箋が drawio に追記されている
- 既存の図の要素を破壊していない

---

**Remember**: 設計の曖昧さは実装前に潰す。付箋は遠慮なく貼る。
