---
paths:
  - "frontend/src/**"
---

# フロントエンドアーキテクチャ

詳細は `docs/adr/frontend/0001-architecture.md` を参照すること。

## レイヤー構成（FSD 簡略版）

```
app → pages → features → shared
```

- ALWAYS 依存は上から下への一方向のみ
- NEVER 同一レイヤー内のスライス間で依存しない

## pages

- フラット構成。ドメイン単位のネストはしない
- 各ページは barrel export（`index.ts`）を持つ

```
pages/
├── home/
│   ├── home-page.vue
│   └── index.ts
├── not-found/
│   ├── not-found-page.vue
│   └── index.ts
└── internal-error/
    ├── internal-error-page.vue
    └── index.ts
```

## features

- ドメイン単位でスライス。1スライスに複数の composable を持てる

```
features/
└── auth/
    ├── model/
    │   ├── use-auth.ts
    │   └── types.ts
    └── index.ts
```

## routes

- 概念単位でファイルを分割する（1ルート1ファイルではない）
- component は barrel export 経由の lazy loading にする

```typescript
// OK — barrel 経由
component: () => import('@/pages/home')

// NG — 内部ファイルへの直接参照
component: () => import('@/pages/home/home-page.vue')
```

## Public API（index.ts）

- pages / features の各スライスは必ず `index.ts` を持つ
- 外部から参照できるのは `index.ts` で export されたものだけ
- NEVER スライス内部のファイルを直接 import しない

## 関連

- **vue-reviewer** agent -- FSD依存方向のレビュー時
- skill: `impl-front` -- 画面実装時
