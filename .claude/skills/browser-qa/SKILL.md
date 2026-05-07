---
name: browser-qa
description: Playwright MCP を使って dev サーバーや Storybook の画面をブラウザで確認し、スクリーンショット取得・コンソールエラー検出・表示確認を行う。
---

# browser-qa

Playwright MCP（headless モード）を通じてブラウザを操作し、画面の表示確認・エラー検出・インタラクション確認を行う手順書。

## When to Activate

- 「画面を確認して」「表示を見て」「スクショを撮って」と依頼されたとき
- 「コンソールエラーが出ていないか確認して」と依頼されたとき
- 「表示崩れがないか見て」「このフォームが動くか試して」と依頼されたとき
- 「Storybook でこのコンポーネントを確認して」と依頼されたとき

## 前提条件

dev サーバーまたは Storybook がすでに起動済みであること。

```bash
# frontend/ ディレクトリで実行
mise run dev          # Vite dev server → http://localhost:5173
mise run storybook    # Storybook      → http://localhost:6006
```

## headless モードの制約

Playwright MCP は `--headless` モードで動作する。リアルタイムのブラウザウィンドウは表示されないため、スクリーンショット（`browser_take_screenshot`）とアクセシビリティツリー（`browser_snapshot`）が主な確認手段。

## 対象 URL の判断基準

| 確認対象 | デフォルト URL |
|---------|--------------|
| Vite dev server | `http://localhost:5173` |
| Storybook | `http://localhost:6006` |
| ユーザー指定 | 依頼に含まれる URL をそのまま使う |

## ワークフロー

### Step 1: 対象 URL を特定する

- 「Storybook で〜」→ `http://localhost:6006`
- 「〜ページを確認して」→ `http://localhost:5173/{path}`
- URL が明示されていない場合は dev server をデフォルトとする

### Step 2: ナビゲーションと初期スクリーンショット

`browser_navigate` で対象 URL に遷移し、即座にスクリーンショットを取得する。

```
# PASS: 遷移直後にスクリーンショットを撮る
browser_navigate → browser_take_screenshot

# FAIL: スクリーンショットを撮らずにいきなり操作を始める
browser_navigate → browser_click
```

ALWAYS ページロードに時間がかかる場合は `browser_wait_for` で要素の出現を待つ。`setTimeout` 等の固定待機は使わない。

### Step 3: コンソールエラーとネットワークエラーの検出

`browser_console_messages` と `browser_network_requests` でエラーを洗い出す。

| 種類 | 対応 |
|------|------|
| `console.error` / Vue warn | 必ず報告 |
| `console.warn` | 報告（重要度を添える） |
| 4xx / 5xx レスポンス | 必ず報告 |

```
# PASS: エラーの内容を引用して報告する
"[error] TypeError: Cannot read properties of undefined (reading 'name') at CustomerList.vue:42"

# FAIL: エラーの有無だけ伝えてメッセージを省略する
"コンソールエラーがありました"
```

### Step 4: 指定された操作の実行

ユーザーが操作の確認を依頼している場合、操作前後それぞれでスクリーンショットを撮る。各操作パターン（フォーム入力・クリック・ホバー・ダイアログ等）の詳細は [operations.md](operations.md) を参照。

### Step 5: DOM 構造の確認（表示崩れの深掘り）

スクリーンショットだけでは判断が難しい場合、`browser_snapshot` でアクセシビリティツリーを取得する。`browser_evaluate` で JavaScript を直接実行して要素の状態を調べることもできる。

```
# PASS: スクリーンショットで崩れが疑われる箇所を snapshot で深掘りする
browser_take_screenshot →（崩れを発見）→ browser_snapshot

# FAIL: 最初から snapshot だけで済ませる（視覚的な崩れを見落とす）
```

### Step 6: 結果レポートの作成

```
## ブラウザ確認レポート

### 確認 URL
- {URL}

### スクリーンショット
- 初期表示: {説明}
- 操作後: {説明}（操作を行った場合）

### コンソールエラー
- なし / あり（内容を列挙）

### ネットワークエラー
- なし / あり（ステータスコード + URL を列挙）

### 確認事項
- {確認項目ごとに PASS / FAIL / 要確認}

### 気になる点
- {依頼外でも気づいた問題があれば記載}
```

## エラーハンドリング

- **サーバー未起動**: `browser_navigate` が接続拒否 → ユーザーに `mise run dev` の実行を依頼して終了
- **要素が見つからない**: `browser_wait_for` 失敗 → `browser_snapshot` で DOM 状態を確認し、ローディング中かエラー状態かを判断
- **スクリーンショットが真っ白**: Vue のマウント失敗の可能性 → `browser_console_messages` でエラーを確認

## 成功指標

- 依頼されたすべてのページ / story のスクリーンショットが取得できている
- コンソールエラー・ネットワークエラーの有無が内容付きで明示されている
- 操作を依頼された場合、操作前後のスクリーンショットがある
- レポートに確認項目ごとの PASS / FAIL が記載されている

---

**Remember**: headless ブラウザは「見えない目」で操作する。スクリーンショットを惜しまず撮り、操作の前後を記録することで、ユーザーが「その場で見ていた」と同じ情報量を再現する。
