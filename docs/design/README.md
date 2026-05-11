# 設計ドキュメント

ドメインごとの設計ドキュメントを管理する。

## 共通

- [対訳表（ユビキタス言語）](対訳表.md) - 日本語の業務用語とコード上の英語名の対応（全ドメイン横断）
- [集約俯瞰図](集約俯瞰図.md) - 全 9 ドメイン・11 集約を 1 枚で俯瞰（集約間 ID 参照・カスケード削除マップ）

## ドメイン一覧

| ドメイン | 概要 | ドキュメント |
|---------|------|-------------|
| ユーザー (User) | LINE ログインで認証する利用者 | [user/](user/index.md) |
| 家計簿 (Household) | 収支を共有管理する単位 | [household/](household/index.md) |
| メンバーシップ (Membership) | ユーザーと家計簿の所属関係（ロール: 管理者 / メンバー） | [membership/](membership/index.md) |
| 招待 (Invitation) | 招待 URL（24 時間・1 回限り） | [invitation/](invitation/index.md) |
| カテゴリ (Category) | 収入／支出の分類（家計簿レベル管理） | [category/](category/index.md) |
| 個人カテゴリ設定 (UserCategoryPreference) | ユーザーごとのカテゴリ表示・並び順・固定費フラグ | [user-category-preference/](user-category-preference/index.md) |
| 収支明細 (Entry) | 1 件の収入または支出（所有者・代理入力対応） | [entry/](entry/index.md) |
| 予算 (Budget) | 支出カテゴリの月次予算 | [budget/](budget/index.md) |
| チャット (ChatMessage) | 家計簿のグループチャットメッセージ | [chat-message/](chat-message/index.md) |

> 各ドメインの詳細設計は sudo モデリング（`/review-design`）を経て確定させる。要件の根拠は [要件定義書](../要件定義書.md) を参照。
