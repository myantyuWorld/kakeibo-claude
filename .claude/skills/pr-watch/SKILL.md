---
name: pr-watch
description: 現タスクのPRのclose状態を Monitor で監視し、merge された瞬間に次の backlog 項目へ自動着手する（merge駆動のタスク送りループ）
---

# PR ウォッチ（PR close → 次タスク）

現タスクの PR が close されたら次を動かす。人間は **Draft → Ready → merge（承認ゲート）だけを握り**、AI は検知して次に着手する。**現 PR を `Monitor` で見張り、merge/close された瞬間に通知→着手する**（`gh pr view <n>` を30秒粒度でポーリングし terminal 状態でだけ発火＝空振り起動が無い）。**kakeibo-claude をカレントにして起動**すること。

Monitor コマンド例（PR #N を監視・merge/close で1回だけ通知して終了）:

```sh
cd <repo> && while true; do s=$(gh pr view <N> --json state -q '.state' 2>/dev/null); case "$s" in MERGED|CLOSED) echo "PR #<N> が $s"; break;; esac; sleep 30; done
```

`persistent: true` で張る。通知が来たら下の手順を回す。

## 1 tick の手順

1. **現タスクの PR を確認** — `backlog.md` の「レビュー中」行の PR 番号を読み、状態を取る:
   `gh pr view <n> --json number,state,isDraft,mergedAt,title`
   - `state` が **OPEN** → 何もしない。「まだレビュー中」とだけ報告して tick 終了。**Draft のままでも催促しない**（Ready 化は人間の判断）。

2. **MERGED だった** → 完了として送る:
   - **後片付け**（人に「main 綺麗にして」と言わせない）:
     `git switch main && git pull` → merged ブランチを削除（local `git branch -d`、remote は `gh pr merge --delete-branch` 済みでなければ `git push origin --delete`）→ `git status` が clean なことを確認。
   - `backlog.md`: 該当行を「完了ログ」へ移し、PR 番号と merged 日を記録する。
   - **次の着手項目を選ぶ**: 実装キュー → 要望キューの順に「未着手」を**優先度順**に1件。
   - **着手判定は二値**: その行の**「ゲート」欄**が**「なし」 → 確認せずそのまま着手**（着手のたびに「確認を挟みます」と止めない）。**「あり」 → 着手せず停止**し、その点をまとめて質問する。ゲートは本来 `/loop-review` の仕分けで潰しておくもの。
   - **「ゲートなし」を「着手可」と読まない**: **優先が `未設定` / `—` の行は、ゲートが無くても着手しない**。ステータス `保留` も同じ。設計が済んでいることと、やると決まっていることは別。**大きな新機能は実施判断が別に必要**。
   - **着手**: 種別に応じて既存スキルへ委譲する（手順を再発明しない）:
     | 種別 | 起動するスキル |
     |---|---|
     | API | `/impl-api <domain> <usecase>` |
     | 画面 | `/impl-front <page>` |
     | 設計 | `/review-design` `/generate-api-spec` `/generate-er` `/review-business-rules` |
     | 基盤 | スキルなし。直接実装し、証拠は下記に揃える |

     スキルが Phase 内で **feature ブランチ作成 → commit → push → Draft PR 作成 → 実装意図のインラインコメント**まで実施する（`.claude/rules/common/git-workflow.md` 準拠）。
   - **証拠を PR 本文に貼る**（下記「証拠」）→ `backlog.md` を「レビュー中(Draft PR)」＋PR番号に更新 → **その新 PR を `Monitor` で監視開始**（次の merge を待つ）。

3. **CLOSED（未merge）だった** → 却下として記録:
   - `backlog.md`: 該当行を「見送り」へ（**理由を必ず残す**）or 差し戻して「未着手」に戻す。次項目へ（手順2と同様に着手）。

## 証拠（PR 本文に貼るもの）

承認は「実物を見る」ではなく **証拠を添えて判定する**。動作の可否は機械が保証し、人間は「これで欲しいものか」だけを見る。

| 証拠 | 取り方 | 何を保証するか |
|---|---|---|
| backend lint / test | `cd backend && mise run lint && mise run test-all` | 動作の可否（単体＋統合） |
| frontend lint / test | `cd frontend && mise run lint && mise run test-run` | 動作の可否 |
| CI | PR に対して GitHub Actions が自動実行（`ci-backend` / `ci-frontend`） | 環境依存でないこと |
| AI レビュー結果 | `/impl-api` `/impl-front` が実行（CRITICAL 0 件） | 規約・設計との整合 |
| 画面のスクショ | フロント変更時のみ `/browser-qa` | UX の判断材料（人間が見る） |

lint/test は `.claude/settings.json` の Stop hook でも毎回走る。**hook が通っている＝証拠が取れている**ので、結果を貼るだけでよい。

## Stop / 予算（永遠に走らせない）

- 未着手が無い／全項目の優先度が未設定 → **停止して報告**し、`/loop-review` を提案する。
- **同時 open PR は1本**（solo・レビューは直列）。複数開かない。
- 同種の失敗が2回 → 停止して人間へエスカレーション。テスト修正は既存スキルの上限（3回）に従う。
- **`feat/` ブランチの push ＋ Draft PR 作成はループに委任済み**。**main 直 push・Ready 化・merge・依存追加・マイグレーション実行・外部送信は本人トリガー**（このループは merge 検知〜次タスクの Draft PR 作成まで）。

## 原則

- **merge ＝ 人間の承認ゲート。** AI はそれを検知して次を動かすだけ。**勝手に Ready 化・merge しない。**
- **証拠を必ず添える。** 新 PR には lint / test / CI / AI レビューの結果を貼る（見なくても動作の可否が判定できる状態にする）。
- **業務ルール・優先度は人間が握る。** AI は着手と証拠添付まで。
- **止まる場所は仕分けフェーズに寄せる。** 未確定点を着手直前に聞くと「merge だけ握れば進む」という merge 駆動の利点が消える。聞くのは `/loop-review` の時。
