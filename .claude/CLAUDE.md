# 家計簿アプリ (kakeibo)

複数人（夫婦・パートナーなど）で共有する家計簿アプリ。スプレッドシート運用をシステム化し、入力の手間削減・年単位の比較・視覚化を実現する。

## 配信形態

**LIFF ミニアプリ（LINE 内 Web アプリ）** として提供する。詳細は [要件定義書](../docs/要件定義書.md) §5.5 を参照。

- 認証: LINE Login API
- 通知: LINE Messaging API（公式アカウント）経由の Push
- 外部ブラウザ単独での利用は v1 では想定しない

## 主要ドメイン（暫定）

実装着手前に sudo モデリング（`/review-design`）でドメインを確定させる。現時点での想定:

| ドメイン | 概要 |
|---------|------|
| ユーザー (User) | アプリ利用者。LINE ログインで認証。プロフィール（ニックネーム・サムネ）を持つ |
| 家計簿 (Household) | 収支を共有管理する単位。1 ユーザーは 1 家計簿のみに所属 |
| メンバーシップ (Membership) | ユーザーと家計簿の所属関係。ロール（管理者 / メンバー）を持つ。管理者は 1 家計簿に 1 人（作成者） |
| カテゴリ (Category) | 収支の分類。収入／支出の 2 系統。カテゴリ自体は家計簿レベルで管理（管理者のみ編集可） |
| 個人カテゴリ設定 (UserCategoryPreference) | ユーザーごとのカテゴリ表示・並び順・**固定費フラグ**の設定。固定費フラグはメンバーが自身の使い方に応じて支出カテゴリ単位で個別に設定する |
| 収支明細 (Entry) | 1 件の収入または支出。日付・区分・カテゴリ・金額・メモ・**所有者**（誰の収支か）を持つ。所有者は登録時に選択可能（代理入力対応） |
| 予算 (Budget) | 支出カテゴリごとの月次予算（管理者が設定） |
| 招待 (Invitation) | 招待 URL。有効期限 24 時間・1 回限り |
| チャット (ChatMessage) | 家計簿に紐づくグループチャットのメッセージ。既読・オンライン状態を扱う |

> NOTE: 詳細なエンティティ・値オブジェクト・業務ルールは `docs/design/` 配下のドメイン設計ドキュメントで管理する。

## 重要な業務ルール

- **暦月集計**: 月次集計・固定費自動コピー・予算判定はすべて暦月（1 日〜末日）固定
- **固定費の月初自動コピー**: 各メンバーが固定費フラグを立てた支出カテゴリについて、そのメンバー自身の前月明細（金額・所有者・メモ）を月初に自動コピー
- **代理入力**: 収支明細の所有者は登録時に他メンバーへ変更可能（編集履歴は v1 では保持しない）
- **メンバーの権限**: 所属家計簿の収支明細は所有者を問わず閲覧・編集・削除可能
- **退会時のデータ**: 退会してもユーザーの収支明細は家計簿に残る（家計簿削除時にのみ全データ削除）

## 技術スタック

### バックエンド

- Go 1.26.2 / Echo v4
- PostgreSQL / sqlx
- マイグレーション: sql-migrate
- テスト: testify（テーブル駆動テスト + mock）

### フロントエンド

- Vue 3 / TypeScript / Vite
- vue-router（レイアウト: meta 方式、lazy loading）
- openapi-fetch + openapi-typescript（API クライアント + 型生成）
- vee-validate + zod（フォームバリデーション）
- テスト: Vitest + happy-dom + MSW
- Storybook 10（コンポーネントカタログ + インタラクションテスト）
- アーキテクチャ: FSD（app → pages → features → shared）

## コマンド

直接 `go test` / `pnpm run` 等を実行せず、必ず `mise run` 経由で実行すること。

### 共通（ルートで実行）

| コマンド | 用途 |
|---------|------|
| `mise run setup` | backend/.env と frontend/.env を作成 |
| `mise run up` | Docker コンテナ起動 |
| `mise run down` | Docker コンテナ停止 |

### バックエンド（`backend/` で実行）

| コマンド | 用途 |
|---------|------|
| `mise run dev` | Air によるホットリロード開発サーバー起動 |
| `mise run lint` | golangci-lint 実行 |
| `mise run test` | 単体テスト実行（`./internal/...`） |
| `mise run test-integration` | 統合テスト実行（`./tests/...`） |
| `mise run test-all` | 全テスト実行 |

### フロントエンド（`frontend/` で実行）

| コマンド | 用途 |
|---------|------|
| `mise run dev` | Vite 開発サーバー起動 |
| `mise run build` | プロダクションビルド |
| `mise run lint` | ESLint 実行 |
| `mise run format` | Prettier フォーマット |
| `mise run test` | Vitest ウォッチモード |
| `mise run test-run` | Vitest 単発実行 |
| `mise run generate-api` | OpenAPI 仕様書から型定義を生成 |
| `mise run storybook` | Storybook 開発サーバー起動（port 6006） |

## Git ルール

`.claude/rules/common/git-workflow.md` を参照すること。

## 開発ワークフロー（Loop Contract）

このリポジトリで作業する時の "止め方と検証" の契約。反復の号令・過剰確認・憶測を減らすためのガードレール。**個人開発＝1人が最大の制約**なので、人間が握るのは「何を作るか（優先度）」と「これで欲しいものか（merge）」の2点だけに寄せる。

| 要素 | 中身 |
|---|---|
| **Trigger** | フィードバック・タスクの入口は `backlog.md`（キュー）。意見も実装タスクもまずここに1行積む。仕分けは `/loop-review`、消化は `/pr-watch` |
| **Scope** | `backend/internal` `backend/pkg` `frontend/src` `docs/` は自動編集可。`backend/migrations` と依存追加（go.mod / package.json）は変更前に一言。`.golangci.yml` `eslint.config.ts` は変更禁止（hook でブロック） |
| **Permission** | 編集・commit ＝ 依頼＋確認OKで実施。**`feat/` ブランチの push ＋ Draft PR 作成は `/pr-watch` ループに委任**。**main 直 push・Draft の Ready 化・merge・マイグレーション実行・依存追加・外部送信は本人トリガー**（コマンドを提示し実行は本人） |
| **Verification** | backend `mise run lint` ＋ `mise run test-all` ／ frontend `mise run lint` ＋ `mise run test-run`（Stop hook で自動実行）＋ **CI green**（`ci-backend` / `ci-frontend`）＋ AI レビューの CRITICAL 0 件。**UI は `/browser-qa` のスクショ**を証拠に添える |
| **Budget** | **1ユースケース／1画面 ＝ 1ブランチ ＝ 1 Draft PR**。同時 open PR は1本（レビューは直列）。段階を小さく刻む |
| **Stop / Escalation** | 同種エラー2回で手を止めて報告。テスト修正は3回で打ち切り。**未確定点は仕分け時に一括確認済みが前提 → `backlog.md` のゲート欄が空なら着手ごとに確認せず走る／埋まっていれば着手せず停止して質問**。優先度が未設定なら着手しない |
| **Reporting** | 設計判断・"なぜ" は `docs/` へ（業務ルール → `docs/design/domains/<domain>/業務ルール.md`、アーキ・技術選定 → `docs/adr/<backend\|frontend>/NNNN-*.md`、未決の論点 → `docs/design/qa/`）。**要望の状態と順番は `backlog.md`**（キューには決定の理由を書かない） |

### 運用ルール

1. **commit はセット。** 機能変更が確認OKになったら規約に沿って commit まで一気に進める（毎回「commit」と号令させない）。ただし main 直 push・merge は本人トリガー。
2. **エラーは診断を任せてよい。** エラーは原因仮説なしで貼ってOK。こちらが切り分け→対処案まで出す（実行＝人／診断＝AI に分離する）。
3. **業務ルールは人間が握る。** 家計簿の運用ルール（暦月集計・固定費コピー・代理入力の扱い等）は立案・推測しない。決定は `docs/design/domains/<domain>/業務ルール.md` に残す。
4. **聞く場所は仕分け、走る場所は着手。** 未確定点は `/loop-review` でまとめて潰す。着手直前に確認で止めない。

### 承認フロー

承認は「実物を見る」ではなく **証拠を添えて判定する**。承認は2種に分けて扱う。

- **動作の可否** → 機械（lint / test / CI / AI レビュー）で保証。見なくても判定できる。
- **これで欲しいものか（UX）** → 人間が判断。ここだけ実物（`/browser-qa` のスクショ、ローカル `mise run dev`）が要る。

承認ゲートは1段: **Draft PR → 人間が Ready にして merge**。`/impl-api` `/impl-front` は必ず Draft で PR を作るので、**Ready 化そのものが人間の意思表示**になる。AI は Draft のまま放置された PR を催促しない。

### タスク送りループ

現タスクの PR が close されたら次タスクに自動着手する heartbeat ループ。手順の実体は `pr-watch` skill。

- **Trigger**: 現タスクの PR が close（merged / closed-unmerged）。`Monitor` でその PR を監視し、30秒粒度・terminal 状態でだけ発火する（空振り起動なし）。
- **動作**: merged → 後片付け（`main` 同期＋merged ブランチ削除）→ backlog を完了ログへ → 次の**優先度順**の未着手に着手（種別に応じて `/impl-api` `/impl-front` へ委譲）→ 証拠を貼って Draft PR → その新 PR を監視開始。closed(未merge) → 見送り／差し戻しとして記録 → 次へ。
- **Stop**: 未着手なし／全項目の優先度未設定／同種エラー2回 → 停止して報告し、`/loop-review` を促す（永遠に走らせない）。

### ループ系スキル

| スキル | 用途 |
|---|---|
| `/loop-review` | 要望の仕分け → ルール反映 → Loop Contract 更新（キューが空／新しい要望が来た時） |
| `/loop-diagnosis` | 操作ログの診断（月1回。繰り返し・過介入・催促を洗い出す） |
| `/pr-watch` | merge 駆動のタスク送り（PR close → 次タスク着手） |

## ドキュメント

実装時は以下のドキュメントを参照すること:

- [要件定義書](../docs/要件定義書.md) - 背景・目的・スコープ・機能要件・非機能要件・制約事項
- [設計ドキュメント](../docs/design/README.md) - ドメインごとの設計ドキュメント入口
- [フロントエンド ADR](../docs/adr/frontend/0001-architecture.md) - FSD アーキテクチャの決定記録
- [API 仕様書](../docs/api/openapi.yml) - OpenAPI 定義
