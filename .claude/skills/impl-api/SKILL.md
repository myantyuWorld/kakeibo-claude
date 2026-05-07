---
name: impl-api
description: 指定ドメイン・ユースケースの API 実装を一気通貫で実行する（ドキュメント確認 → ドメインモデル → API → テスト → レビュー → PR）。「XXを実装して」「XX APIを作って」など、バックエンドの機能実装を依頼されたときに使用する。
argument-hint: "{ドメイン名} {ユースケース名}（例: customer 顧客登録）"
---

# impl-api

指定されたドメイン・ユースケースの実装を、ドキュメント確認からPR作成まで自動で実行する。

## When to Activate

- 「XXを実装して」「XX APIを作って」とバックエンドの機能実装を依頼されたとき
- 設計ドキュメント（業務ルール.md、ドメインモデル.md）が揃った状態で実装フェーズに入るとき

## 設計思想

シーケンシャルパイプラインパターンを採用。各 Phase の成果物が次の Phase の入力になる直列構成で、Phase 5 で複数の Agent を順に呼び出して品質を担保する。1ユースケース単位で完結させることで、差分が小さく レビューしやすい PR を生成する。

## アーキテクチャ

```
Phase 1: ドキュメント確認
  ↓
Phase 2: ドメインモデル実装（未存在時のみ）
  ↓
Phase 3: API 実装（handler → service → repository）
  ↓
Phase 4: テスト生成・実行 ←→ 修正ループ（最大3回）
  ↓
Phase 5: AIレビュー・自動修正
  go-reviewer → security-reviewer → go-test-reviewer
  ↓
Phase 6: PR 作成
  ↓
Phase 7: 実装意図コメント投稿
```

## 入力

- ドメイン名（例: customer）
- ユースケース名（例: 顧客登録）

## 実行フロー

### Phase 1: ドキュメント確認

`docs/design/{ドメイン名}/index.md` を読み、リンクされている関連ドキュメントを全て確認する。

### Phase 2: ドメインモデル実装

`internal/{ドメイン名}/{ドメイン名}.go` が存在しない場合のみ実行する。

`.claude/rules/backend/domain-modeling.md` に従い、ドメインモデルを実装する。

### Phase 3: API 実装

1. `docs/design/{ドメイン名}/業務ルール.md` から指定ユースケースのビジネスルールを確認する
2. `docs/design/{ドメイン名}/ドメインモデル.md` の対訳表で命名を確認する
3. 以下のファイルを実装・更新する

#### handler.go

- リクエスト/レスポンス構造体（exported）を定義する
- Service インターフェースに該当メソッドを追加する
- ハンドラメソッド（exported）を実装する
  - リクエストのバインド
  - service の呼び出し
  - レスポンスの返却

#### service.go

- Repository インターフェースに必要なメソッドを追加する
- service にユースケースメソッドを実装する
  - ビジネスルールの実行
  - トランザクション管理（必要な場合）

#### repository.go

- repository に SQL 発行メソッドを実装する
- DB の行構造体（unexported）と toModel 変換を定義する

4. ルーティングを登録する

### Phase 4: テスト生成・実行

1. `/generate-api-tests {ドメイン名} {ユースケース名}` スキルを実行する
2. `mise run lint` を実行する
3. `mise run test` を実行する
4. 失敗があれば修正して再実行する（最大3回）

### Phase 5: AIレビュー・自動修正

1. go-reviewer エージェントを実行する
2. CRITICAL / HIGH の指摘があれば自動修正する
3. 修正後、再度 `mise run lint && mise run test` を実行する
4. security-reviewer エージェントを実行する
5. CRITICAL の指摘があれば自動修正する
6. go-test-reviewer エージェントを実行する
7. CRITICAL の指摘（受入基準の未カバー等）があればテストを追加する

### Phase 6: PR 作成

`.claude/rules/common/git-workflow.md` に従い、以下を実行する。

1. feature ブランチを作成する
2. 変更をコミットする
3. push する
4. Draft PR を作成する（テンプレートに従う）

### Phase 7: 実装意図コメント

`.claude/rules/common/git-workflow.md` の「PR のインラインコメント」に従い、差分の該当行に実装意図を投稿する。

すべてのコメント投稿が完了したら、PR の URL を出力して終了する。

## 遵守ルール

- `.claude/rules/backend/architecture.md` の各レイヤーサンプルに従う
- `.claude/rules/backend/error-handling.md` のレイヤーごとの使い分けに従う
- 1ユースケース = 1実装単位。複数ユースケースを同時に実装しない

## エラーハンドリング

- **Phase 1（ドキュメント確認）が失敗**: 設計ドキュメントが見つからない場合、AskUserQuestion でパスを確認する。存在しない場合は中断し、先に設計を行うよう案内する
- **Phase 3（API 実装）でビルドエラー**: 即座に修正を試みる。3回失敗したら中断してエラー内容を報告する
- **Phase 4（テスト）が3回失敗**: テスト修正ループを打ち切り、失敗しているテストと原因を報告して Phase 5 に進まない
- **Phase 5（AIレビュー）で CRITICAL が解消できない**: 最大2回の自動修正を試み、解消できなければ指摘内容を報告して Phase 6 に進まない
- **Phase 6（PR 作成）が失敗**: push 失敗やブランチ競合の場合、エラー内容を報告してユーザーに判断を委ねる

## アンチパターン

- 設計ドキュメントを読まずに実装を始める（Phase 1 をスキップする）
- 複数ユースケースを1つの PR にまとめる
- AIレビューの CRITICAL 指摘を無視して PR を作成する
- テストが失敗した状態でコミットする

## 成功指標

- lint と全テストが通る状態で PR が作成されている
- AIレビュー（go-reviewer, security-reviewer, go-test-reviewer）で CRITICAL 指摘が0件
- 受入基準に対応するテストケースが存在する
- PR に実装意図コメントが投稿されている

---

**Remember**: 1ユースケース1PR。設計ドキュメントに書かれていないことは実装しない。
