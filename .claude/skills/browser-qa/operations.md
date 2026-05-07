# ブラウザ操作リファレンス

SKILL.md の Step 4〜5 で参照する、各操作パターンの詳細手順。

## フォーム入力

```
browser_fill_form（フォーム要素のセレクター、入力値）
→ browser_take_screenshot（入力後の状態）
```

## クリック

```
browser_click（ボタン / リンクのセレクター）
→ browser_take_screenshot（クリック後の状態）
```

## ホバー

```
browser_hover（対象セレクター）
→ browser_take_screenshot（hover 状態）
```

## セレクトボックス

```
browser_select_option（セレクター、選択値）
→ browser_take_screenshot
```

## ダイアログ（confirm / alert）

confirm が出ることが分かっている場合は操作前にハンドラを登録する。

```
# PASS: ダイアログが予期される操作の前にハンドラを登録する
browser_handle_dialog（accept: true / false）
→ browser_click（ダイアログを起動する操作）
→ browser_take_screenshot

# FAIL: ダイアログが出た後に登録しようとする（既にブロックされている）
browser_click → browser_handle_dialog
```

## ファイルアップロード

```
browser_file_upload（input[type="file"] のセレクター、ファイルパス）
→ browser_take_screenshot
```

## レスポンシブ確認

依頼された場合に `browser_resize` でビューポートサイズを変更してスクリーンショットを撮る。

```
# モバイル
browser_resize（width: 375, height: 812）→ browser_take_screenshot

# タブレット
browser_resize（width: 768, height: 1024）→ browser_take_screenshot

# デスクトップ
browser_resize（width: 1440, height: 900）→ browser_take_screenshot
```

## セレクター指定の優先順位

`browser_click` 等でセレクターを指定する際、data-testid が付いている要素はそれを優先して使う。

```
# PASS: data-testid を使う（実装変更に強い）
browser_click（セレクター: "[data-testid='submit-button']"）

# FAIL: クラス名や DOM 構造で指定する（実装変更で壊れやすい）
browser_click（セレクター: ".form-container > div:last-child > button"）
```

## Storybook 固有の手順

Storybook では story ごとに URL が変わるため、確認したい story の ID を特定してから遷移する。

```
# トップページで story 一覧を確認する
browser_navigate（http://localhost:6006）
→ browser_snapshot（サイドバーの story ツリーを確認）

# 特定の story に直接遷移する
browser_navigate（http://localhost:6006/?path=/story/features-customer-customerform--default）
→ browser_take_screenshot

# インタラクションテストの結果を確認する
browser_snapshot（"Interactions" パネルの状態を確認）
```

Storybook の story URL の形式はコンポーネントの配置パスに依存する。FSD 構成では以下のようになる。

```
/story/{layer}-{feature}-{componentname}--{storyname}
例: /story/features-customer-customerform--with-validation-error
```
