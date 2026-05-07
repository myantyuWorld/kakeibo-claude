---
name: impl-front
description: フロントエンドの画面実装を一気通貫で実行する（設計確認 → ルート定義 → composable → Vue コンポーネント → テスト → Story → レビュー → PR）。「画面を実装して」「ページを作って」など、フロントエンドの機能実装を依頼されたときに使用する。
argument-hint: "{ページ名}（例: 顧客一覧）"
---

# impl-front

指定されたページの実装を、設計確認からPR作成まで自動で実行する。

## When to Activate

- 「画面を実装して」「ページを作って」とフロントエンドの機能実装を依頼されたとき
- API が実装済みで、対応するフロントエンド画面を作成するとき

## 設計思想

シーケンシャルパイプラインパターンを採用。FSD アーキテクチャに従い、composable（ロジック）→ コンポーネント（UI）→ テスト → Story の順で積み上げる。Phase 6 で複数の Agent を順に呼び出し、FSD 依存方向・アクセシビリティ・テスト網羅性を検証する。

## アーキテクチャ

```
Phase 1: 設計確認
  ↓
Phase 2: ルート定義
  ↓
Phase 3: composable 実装（features/{ドメイン}）
  ↓
Phase 4: Vue コンポーネント実装（pages/{ページ}）
  ↓
Phase 5: テスト生成・実行 ←→ 修正ループ（最大3回）
  ↓
Phase 6: ブラウザ確認（/browser-qa）
  ↓
Phase 7: AIレビュー・自動修正
  vue-reviewer → a11y-reviewer → vue-test-reviewer
  ↓
Phase 8: PR 作成
  ↓
Phase 9: 実装意図コメント投稿
```

## 入力

- ページ名（例: 顧客一覧）

## 実行フロー

### Phase 1: 設計確認

1. 対応する設計ドキュメントがあれば確認する（`docs/design/` 配下）
2. 対応する API 仕様書があれば確認する（`docs/api/` 配下）
3. ページの要件（表示項目、操作、画面遷移）を把握する

### Phase 2: ルート定義

1. `frontend/src/app/routes/` に該当ルートを追加する
2. component は barrel export 経由の lazy loading にする
3. meta（layout、認証要否）を設定する

### Phase 3: composable 実装

1. `frontend/src/features/{ドメイン名}/` にスライスを作成する（既存ならスキップ）
2. API 呼び出し用の composable を実装する（openapi-fetch 使用）
3. フォームがある場合は vee-validate + zod で バリデーション用 composable を実装する
4. 型定義を `model/types.ts` に分離する
5. `index.ts` で Public API を export する

### Phase 4: Vue コンポーネント実装

1. `frontend/src/pages/{ページ名}/` にページコンポーネントを作成する
2. `.claude/rules/frontend/vue-coding.md` に従い実装する
3. `index.ts` で barrel export する

### Phase 5: テスト生成・実行

1. `/generate-frontend-tests` スキルを実行する
2. Storybook の Story を `.claude/rules/frontend/storybook.md` に従い作成する
3. `cd frontend && mise run test-run` を実行する
4. 失敗があれば修正して再実行する（最大3回）

### Phase 6: ブラウザ確認

`/browser-qa` スキルを実行し、実装したページの表示確認を行う。

1. dev サーバーが起動していなければ `cd frontend && mise run dev` でバックグラウンド起動する
2. 実装したページに `browser_navigate` で遷移する
3. 初期表示のスクリーンショットを取得する
4. コンソールエラー・ネットワークエラーがないか確認する
5. 主要な操作（フォーム入力、ボタンクリック等）があれば動作確認する
6. 問題があれば修正し、Phase 5 のテストを再実行する

### Phase 7: AIレビュー・自動修正

1. vue-reviewer エージェントを実行する
2. CRITICAL / HIGH の指摘があれば自動修正する
3. 修正後、再度 `cd frontend && mise run lint && mise run test-run` を実行する
4. a11y-reviewer エージェントを実行する
5. CRITICAL の指摘があれば自動修正する
6. vue-test-reviewer エージェントを実行する
7. CRITICAL の指摘（テストケースの漏れ等）があればテストを追加する

### Phase 8: PR 作成

`.claude/rules/common/git-workflow.md` に従い、以下を実行する。

1. feature ブランチを作成する
2. 変更をコミットする
3. push する
4. Draft PR を作成する（テンプレートに従う）

### Phase 9: 実装意図コメント

`.claude/rules/common/git-workflow.md` の「PR のインラインコメント」に従い、差分の該当行に実装意図を投稿する。

すべてのコメント投稿が完了したら、PR の URL を出力して終了する。

## エラーハンドリング

- **Phase 1（設計確認）が失敗**: 設計ドキュメントや API 仕様書が見つからない場合、AskUserQuestion で確認する。API が未実装なら先にバックエンドを実装するよう案内する
- **Phase 3（composable）でビルドエラー**: API の型定義が古い可能性がある。`cd frontend && mise run generate-api` を実行して型を再生成してからリトライする
- **Phase 5（テスト）が3回失敗**: テスト修正ループを打ち切り、失敗しているテストと原因を報告して Phase 6 に進まない
- **Phase 6（ブラウザ確認）で問題検出**: 修正して Phase 5 のテストを再実行する。修正後も解消しない場合は問題内容を報告して Phase 7 に進む
- **Phase 7（AIレビュー）で CRITICAL が解消できない**: 最大2回の自動修正を試み、解消できなければ指摘内容を報告して Phase 8 に進まない
- **Phase 8（PR 作成）が失敗**: push 失敗やブランチ競合の場合、エラー内容を報告してユーザーに判断を委ねる

## アンチパターン

- API 仕様書を確認せずに composable を実装する（型の不一致が起きる）
- FSD の依存方向に違反する import を書く（pages → features → shared の一方向のみ）
- テストが失敗した状態でコミットする

## 成功指標

- lint とテストが通る状態で PR が作成されている
- AIレビュー（vue-reviewer, a11y-reviewer, vue-test-reviewer）で CRITICAL 指摘が0件
- Storybook の Story が作成されている
- PR に実装意図コメントが投稿されている

---

**Remember**: FSD の依存方向を守り、composable でロジックを分離してからコンポーネントを作る。
