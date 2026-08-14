---
name: loop-review
description: ステークホルダーの意見・要望や定期タイミングをトリガーに、ループエンジニアリングの見直しを丸ごと回す。履歴診断→新要望の仕分け→memory/CLAUDE.md/rules/skillへ反映→Loop Contract更新まで一連で行う
---

# ループエンジニアリング 見直しサイクル

「診断して終わり」ではなく、**診断 → 三分法で仕分け → memory / CLAUDE.md / rules / skill に反映 → Loop Contract 更新 → 記録** までを1本で回すための手順。単発の履歴診断だけなら `loop-diagnosis` skill を使う（この skill はそれを内包する上位サイクル）。

## いつ回すか（トリガー）

- **新しいステークホルダー入力が来た時**（kakeibo なら家計簿を一緒に使うパートナーの要望・不満）
- **バックログのキューが空 or 未仕分けが溜まった時**（＝ `/pr-watch` が「着手できる項目なし」で停止した時）
- 前回から1ヶ月経過（定期ハートビート）
- 同じ指示を繰り返している／過剰確認が続いていると気づいた時

## インプットを集める

1. 前回レビュー以降のユーザー発話（`loop-diagnosis` の手順1〜2）。
2. **新しいステークホルダー入力**: `backlog.md` の未仕分け行を読む。誰が・どんな言葉で・何を欲しがっているかを引用付きで確認し、分類（機能/ルール/人間）と優先度を埋める。誰の判断領域か（業務ルール＝人間）を明記する。
3. **実装キューの補充**: 実装フェーズでキューが枯れていたら、`docs/design/機能一覧.md`（26 操作）と `docs/design/画面一覧.md`（19 画面）から**直近で着手する数件だけ**を取り込む。全件を写さない。依存順（User → Household → Membership → Invitation → Category → UserCategoryPreference → Entry → Budget → ChatMessage）と、API が無いと画面が作れない前後関係を考慮する。
4. 現在の Loop Contract（`.claude/CLAUDE.md` の「開発ワークフロー」）と既存の rules / skill / memory を読み、前提を把握する。

## 回す手順

1. **診断** — `loop-diagnosis` skill を実行（繰り返し／過介入／催促の抽出＋HTMLレポート）。
2. **新要望の仕分け** — ステークホルダー要望を三分法にかける:
   - **ループ化すべき**（機能追加・自動化・検証の自動化）
   - **ルール整備で足りる**（CLAUDE.md / rules / skill / memory / hook）
   - **人間判断を残す**（業務ルール・機能の優先順位は本人／利用者が決める）

   要望は「作る／作らない」をこの場で即決しない。Loop Contract の Scope と Verification に落とせる形に噛み砕く。**過剰設計を避ける**（個人開発＝1人が最大の制約）。
3. **設計の疑問を潰す（ゲートを空にする）** — backlog の「ゲート」欄に残っている論点を、**ここでまとめて**質問して決める。決めた内容は docs/ 側へ書く:
   - 業務ルール・受入基準 → `docs/design/domains/<domain>/業務ルール.md`
   - ドメインモデルの変更 → `docs/design/domains/<domain>/ドメインモデル.md`
   - アーキテクチャ・技術選定 → `docs/adr/<backend|frontend>/NNNN-*.md` を新規追加
   - 未決のまま残す論点 → `docs/design/qa/*.md`

   **ここでゲートを空にしておくことが、`/pr-watch` が止まらずに走れる条件**（着手直前に聞かない）。
4. **反映（打ち手を置き場所へ）** — 恒久的な打ち手をそれぞれの家に置く:
   - AI の振る舞いに関わるもの（プロジェクト横断） → memory（type: feedback）
   - プロジェクトのガードレール／ワークフロー → `.claude/CLAUDE.md`
   - コーディング規約・アーキ制約 → `.claude/rules/**`
   - 反復手順 → `.claude/skills/<name>/SKILL.md`
   - 機械的に強制できるもの → `.claude/settings.json` の hook

   **勝手に commit しない**（提案して委ねる。push は `/pr-watch` の消化フェーズか本人トリガー）。
5. **Loop Contract 更新** — 7要素（Trigger / Scope / Permission / Verification gate / Budget / Stop・Escalation / Reporting）に差分があれば `.claude/CLAUDE.md` を更新する。
6. **記録** — 何を変えたか・次に持ち越す論点を残す。設計判断は docs/（ADR または業務ルール）、要望の状態と順番は `backlog.md`。

## 原則

- **事実ベース。** 要望は「誰の・どんな言葉か」を引用で残す（捏造しない）。
- **優先順位は人間が握る。** 機能の要否と順番は本人・利用者が決める。AI は仕分けと反映を担う。
- **一度に全部やらない。** 費用対効果の高い順に1つずつ仕組みへ移す（月1回眺めて一つずつ）。
- **自分側の問題も見る。** AI の過剰確認がユーザーの催促を誘発していないか毎回チェックする。
- **聞く場所はここ、走る場所は着手。** 未確定点はこのフェーズで一括して潰す。
