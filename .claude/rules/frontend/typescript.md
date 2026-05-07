---
paths:
  - "frontend/src/**/*.ts"
  - "frontend/src/**/*.vue"
---

# TypeScript コーディング規約

## 型

- NEVER `any` を使わない
- NEVER Non-null assertion（`!`）を使わない
- NEVER `null` を使わない。`undefined` を使う
- 戻り値の型アノテーションは不要。TypeScript の型推論に任せる
- 型エイリアスは composable と混ぜず `model/types.ts` に分離する

## import

- ALWAYS パスは `@/` プレフィックスを使用する
- ALWAYS 型のみの import は `import type` を使用する

## API 呼び出し

- openapi-fetch の `{ data, error }` を throw せずそのまま返す
- API ラッパーが1行（`() => client.GET('/path')`）なら api/ セグメントは不要。composable から直接呼ぶ
- API ファイルの命名はエンドポイントに合わせる（`auth-api.ts` ではなく `me.ts`）

## ファイル分割

- 薄い関数を別ファイルに分けない。インライン化できる程度なら呼び出し元にまとめる

## 関連

- **vue-reviewer** agent -- TypeScript規約のレビュー時
