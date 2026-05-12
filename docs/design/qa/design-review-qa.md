# 設計ドキュメント整合性レビュー: 顧客QA 叩き台

## 概要

- **対象**: 設計ドキュメント全体（要件定義書 v0.8 + 共通レイヤ 6 ドキュメント + 9 ドメイン × 3 ファイル）
- **作成日**: 2026-05-13
- **目的**: 実装着手前の最終整合性チェックで検出された矛盾・抜け・曖昧さについて、顧客に確認すべき事項を整理する
- **参照ドキュメント**:
  - [要件定義書 v0.8](../../要件定義書.md)
  - [集約俯瞰図](../集約俯瞰図.md)
  - [機能一覧](../機能一覧.md)
  - [画面一覧](../画面一覧.md)
  - [ER 図 (DBML)](../er.dbml)
  - [ユースケース図](../usecase.drawio)
  - [業務フロー図](../business-flow.drawio)
  - [対訳表](../対訳表.md)
  - [先行 QA: ユースケース図レビュー](usecase-review-qa.md)
- **想定の使い方**: 顧客レビュー会議のアジェンダとして使用。回答結果に応じて要件定義書 / ドメイン業務ルール / ER / 各俯瞰図を更新する

## 凡例

- 重要度: **Critical**（実装方針が割れる、着手前に必須）/ **High**（要件補記が必要、実装と並行可）/ **Medium**（詳細詰めが必要、実装中に解決可）
- 状態: 未回答 / ✅ 回答済み

---

## Critical（実装着手前に潰す）

実装方針が割れるため、顧客レビューでクローズしないとマイグレーション・トランザクション設計に手戻りが発生する論点。

### Q-C1: ユースケース図の UC 数を 21 → 23 に統一してよいか？

- **重要度**: Critical
- **状況**:
  - [usecase.drawio](../usecase.drawio) の楕円を実数えると **23 UC** (uc_login〜uc_chat_push に加え uc_member_list, uc_notification が Round 1〜2 で追加)
  - [README.md:9](../README.md)・[機能一覧.md](../機能一覧.md)・[画面一覧.md](../画面一覧.md) では「21 UC」と表記
- **問い**: 表記を以下のどちらに合わせるか？
  - **(A)** 23 UC（追加 2 件を含める。各俯瞰図の表記を 23 に更新）
  - **(B)** 21 UC（uc_member_list / uc_notification を *include* 関係に整理して数えない）
- **影響**: 各俯瞰図のヘッダ表記、UC ↔ 機能 # マッピング表の網羅性

### Q-C2: User の物理削除 CASCADE 設計を統一する方針 ✅

- **重要度**: Critical
- **状況**:
  - [er.dbml](../er.dbml) の User 参照 FK で `ON DELETE CASCADE` の指定が不揃い
    - **CASCADE あり**: `chat_reads.user_id`, `user_category_preferences.user_id`, `online_statuses.user_id`
    - **CASCADE なし**: `memberships.user_id`, `invitations.issued_by_user_id`, `invitations.used_by_user_id`, `entries.owner_user_id`, `entries.created_by_user_id`, `chat_messages.sender_user_id`
  - v1 では User 物理削除 API は提供せず（ソフト削除 `Status=Withdrawn` のみ）。「将来用」と書かれているが揃っていない
- **問い**: 以下のどちらで統一するか？
  - **(A)** v1 では全 FK を CASCADE なし（NO ACTION）に統一。物理削除は v2 以降で別途 Cleanup バッチを設計
  - **(B)** 「将来用」を維持し、全 FK に `ON DELETE CASCADE` を一律付与（Membership / Entry / Invitation / ChatMessage を CASCADE 化）
- **影響**: v2 で物理削除を実装する際に「一発 DELETE で全部消える」か「個別 cleanup 必要」かが決まる
- **決定 (2026-05-13)**: **(A) 採用。**v1 は全 User 参照 FK を `RESTRICT` で統一。`user_category_preferences.user_id` / `chat_reads.user_id` / `online_statuses.user_id` の CASCADE 指定を撤去。v2 で User 物理削除を実装する際は専用クリーンアップフローを別途設計する。
  - 反映: [er.dbml](../er.dbml), [集約俯瞰図.md](../集約俯瞰図.md)

### Q-C3: UCP の DisplayOrder 一意性を業務ルールから緩めるか？ ✅

- **重要度**: Critical
- **状況**:
  - [user-category-preference 業務ルール](../domains/user-category-preference/業務ルール.md): 「DisplayOrder は同一ユーザー × 同一家計簿でユニーク」
  - [er.dbml:130-148](../er.dbml) では UCP に household_id 列がなく `(user_id, category_id)` UNIQUE のみ。DB レベルで一意性担保不可
  - 並び順一括更新時のレース・部分失敗で重複が発生し得る
- **問い**: 以下のどちらで仕様確定するか？
  - **(A)** 衝突許容（業務ルールを緩める）。表示時は `display_order ASC, created_at ASC` でタイブレーク
  - **(B)** UCP に household_id 列を追加し `(user_id, household_id, display_order)` UNIQUE を貼る（正規化崩しを許容）
  - **(C)** アプリ層で並び順更新 API をシリアル化（行ロック）し、業務ルールどおり一意性を強制
- **影響**: UCP テーブル設計・並び順 API のトランザクション境界・テスト戦略
- **決定 (2026-05-13)**: **(A) 採用。**業務ルールを「`display_order` の重複を潜在的に許容」に緩め、表示時は `(display_order ASC, created_at ASC)` で安定ソート。並び順一括更新 API は単一 Tx 内で全件 UPSERT し、部分失敗はロールバックする。`household_id` 列の追加はしない（正規化を保つ）。
  - 反映: [user-category-preference/業務ルール.md](../domains/user-category-preference/業務ルール.md), [er.dbml](../er.dbml) コメント

### Q-C4: OnlineStatus 更新を middleware に一本化してよいか？ ✅

- **重要度**: Critical
- **状況**:
  - [機能一覧 #25](../機能一覧.md): 副作用に「OnlineStatus.LastSeenAt UPSERT（同一 Tx）」と記載
  - 同 [機能一覧 派生操作](../機能一覧.md): 「OnlineStatus 更新 | 認証付き全 API | middleware で全リクエストで実行」と記載
  - [chat-message 業務ルール](../domains/chat-message/業務ルール.md): 「認証 middleware で全 API リクエストに共通適用するのが望ましい」
  - 2 つの場所で実装方針が二重表現になっており、実装者が両方書く / 片方見落とすリスク
- **問い**:
  - **(A)** middleware に一本化。機能一覧 #25 副作用記述から OnlineStatus を削除し、Tx 境界は「ChatMessage SELECT + ChatRead INSERT」のみに揃える
  - **(B)** API 個別に書く（middleware を使わず、各 API の Tx 内で UPSERT）
- **影響**: ChatMessage 取得 API のトランザクション境界、middleware 実装の必要性
- **決定 (2026-05-13)**: **(A) 採用。**認証 middleware で全 API リクエストに対し `online_statuses` を `now()` UPSERT。機能一覧 #25 の副作用から OnlineStatus を外し、ChatMessage 取得 API の Tx 境界は `SELECT ChatMessage` + `INSERT ChatRead` のみ。OnlineStatus の正確性は middleware のテストで担保する。
  - 反映: [機能一覧.md](../機能一覧.md), [chat-message/業務ルール.md](../domains/chat-message/業務ルール.md)

### Q-C5: 退会済みユーザーを新規明細の所有者として選べない仕様で確定か？ ✅

- **重要度**: Critical
- **状況**:
  - [entry 業務ルール:50](../domains/entry/業務ルール.md): 「OwnerUserID は Active メンバーシップを持つ User のみ」「退会済みユーザーへの代理入力は不可」
  - [要件定義書 §4.1](../../要件定義書.md): 退会済みユーザーの所有明細は「家計簿に保持され、所有者として参照される」と記載されているのみで、**新規明細の所有者として選べない仕様は要件側にない**
  - ユースケース: 例「夫が退会した後に妻が亡夫名義の固定費（住宅ローン等）を妻が代理で入力したい」が暗黙に封じられている
- **問い**:
  - **(A)** 現状維持（業務ルール通り、Active メンバーのみ所有者に指定可能）。要件定義書 §4.1 にも明記
  - **(B)** 退会済みユーザーも所有者として選択可能にする（プルダウンに「退会済みユーザー」として表示）。業務ルールを修正
- **影響**: 代理入力 API の認可ロジック、固定費月初コピーバッチで退会者の前月明細を扱うか（Q-M2 連動）
- **決定 (2026-05-13)**: **(A) 採用。**Active メンバーのみ新規明細の所有者に指定可能。代替手段（現メンバーが自分名義で代理入力する）で十分とユーザー確認済み。Q-M2 も連動して **(A)**（退会者の前月明細はバッチ対象外）として自然解決。要件定義書 §4.1 への明記は最終整流フェーズで対応。
  - 反映: [entry/業務ルール.md](../domains/entry/業務ルール.md)（月初バッチ §1 に明示）

---

## High（要件側補記または用語整理。実装と並行可）

### Q-H1: 「Active な Membership」の定義を対訳表に追記してよいか？ ✅

- **重要度**: High
- **状況**: Membership 集約に Status 列はなく（[er.dbml](../er.dbml)）、実体としては「Membership レコードが存在し、参照 User.Status=Active」を指す。一方、[機能一覧 凡例](../機能一覧.md)・[entry 業務ルール](../domains/entry/業務ルール.md) で「Active メンバー」「Active な Membership」が頻出するが、[対訳表](../対訳表.md) に定義がない
- **問い**: 対訳表に「Active メンバー = Membership レコード存在 + 参照 User.Status=Active」を追記してよいか？
- **影響**: 用語の二重解釈防止、新規実装者のオンボーディング
- **決定 (2026-05-13)**: **追記する。**Membership に `status` 列は追加せず、対訳表のメンバーシップ節に「Active メンバー = Membership レコード存在 + 参照 User.Status=Active」を明記。SQL では `memberships JOIN users ON ... WHERE users.status = 'Active'` を共通形とする。
  - 反映: [対訳表.md](../対訳表.md)

### Q-H2: CSV インポートの所有者識別を Nickname から UUID に変えてよいか？

- **重要度**: High
- **状況**: [entry 業務ルール](../domains/entry/業務ルール.md): CSV の OwnerNickname で家計簿内 Active メンバーを識別。Nickname に UNIQUE 制約はなく、家計簿内に同名メンバーがいた場合の挙動が未定義
- **問い**:
  - **(A)** CSV を UUID または LineUserID ベースに変更（ユーザー視点で扱いづらくなる）
  - **(B)** Nickname のまま、同名検出時は「行を保留し UI でユーザーに選択させる」フローを追加
  - **(C)** Nickname のまま、同名検出時はエラーで全体を rollback（運用は同名を作らないことで対処）
- **影響**: CSV フォーマット仕様、インポート UI のフロー、受入基準の追加

### Q-H3: メンバー一覧 API の取得項目を要件定義書に追記してよいか？

- **重要度**: High
- **状況**:
  - [要件定義書 §4.1](../../要件定義書.md): 「ニックネーム・サムネ・ロール」
  - [membership 業務ルール](../domains/membership/業務ルール.md): 「ニックネーム・サムネ・ロール・**参加日時**」「退会済みユーザーも『退会済みユーザー』として一覧に含める」
- **問い**: 要件定義書 §4.1 に「参加日時を含む」「退会済みユーザーも一覧に表示（退会済み表示）」を追記してよいか？
- **影響**: 要件と業務ルールの整合、UI 設計の根拠

### Q-H4: 業務フロー図に管理者退会の動線を追加してよいか？

- **重要度**: High
- **状況**: [要件定義書 §11](../../要件定義書.md) と [user 業務ルール](../domains/user/業務ルール.md) で「管理者退会は警告のみで通過、家計簿は管理者不在状態で残る」と確定済みだが、[業務フロー図](../business-flow.drawio) には管理者レーンに退会ステップがない
- **問い**: 業務フロー図の管理者レーンに「アカウント退会（任意、警告のみ）」ノードを追加してよいか？
- **影響**: 業務フロー図の網羅性、フロー図を見る関係者の認識ギャップ防止

### Q-H5: 招待 URL の能動的取消を要件 §1.3 Out of Scope に明記してよいか？

- **重要度**: High
- **状況**: [invitation 業務ルール](../domains/invitation/業務ルール.md) に「招待 URL の手動取り消し（v1 スコープ外、QA-C1 で確定）」とあるが、[要件定義書 §1.3 Out of Scope](../../要件定義書.md) には記載なし
- **問い**: 要件 §1.3 に「招待 URL の能動的取り消し機能（再発行による間接無効化のみ提供）」を追記してよいか？
- **影響**: 要件側の網羅性、顧客への期待値調整

### Q-H6: 過去退会レコードからの Nickname / サムネ引き継ぎは行わない方針で確定か？

- **重要度**: High
- **状況**: [er.dbml](../er.dbml) の部分一意 `idx_users_line_user_id_active`（`WHERE status=Active`）により、退会済み LineUserID の再利用は可能。[user 業務ルール](../domains/user/業務ルール.md) の受入基準では「新規 UUID で新規 User 作成」とあるが、Nickname / ThumbnailURL 引き継ぎについて触れていない
- **問い**:
  - **(A)** 引き継がない（プロフィール登録画面に誘導）— 受入基準に明記
  - **(B)** 引き継ぐ（過去退会レコードから Nickname / ThumbnailURL をコピー）
- **影響**: 再ログインユーザーの UX、User INSERT 時の処理ロジック

### Q-H7: UCP は退会時に削除せず、家計簿削除でのみ消える方針で確定か？

- **重要度**: High
- **状況**:
  - [er.dbml](../er.dbml): UCP は `user_id ON DELETE CASCADE` 指定だが v1 では User 物理削除は発生しない
  - 退会＝ソフト削除のため UCP は孤児化（参照する User が Withdrawn 状態で残存）
  - 家計簿削除時は `Household → Category CASCADE → UCP CASCADE` で消える
- **問い**:
  - **(A)** 現状維持（UCP は退会で削除せず、家計簿削除でのみ消える）。業務ルールに「孤児 UCP を許容」を明記
  - **(B)** 退会時に UCP も削除する（user_id 経由でアプリ層で削除を追加）
- **影響**: 退会処理のトランザクション範囲、孤児 UCP の運用説明

### Q-H8: UC 数と機能数のマッピング表を機能一覧に追加してよいか？

- **重要度**: High
- **状況**: [usecase.drawio](../usecase.drawio) の **23 UC** と [機能一覧](../機能一覧.md) の **26 操作** の対応関係が明文化されていない。UC「家計簿を管理する（作成・編集・削除）」= 機能 #5, #6, #7 のような対応を読み手が推測する必要がある
- **問い**: 機能一覧の末尾に「UC ↔ 機能 # マッピング」表を追加してよいか？（システム自動・派生操作も含めて 26 機能を完全網羅）
- **影響**: 機能の網羅性検証、API 設計時の取りこぼし防止

### Q-H9: 集約俯瞰図のカスケード表記を ER と整合させてよいか？

- **重要度**: High
- **状況**:
  - [集約俯瞰図 mermaid 図](../集約俯瞰図.md): `Household -->|CASCADE| ChatRead` を**直接 CASCADE** として描画
  - 実体は ER 上 `Household → ChatMessage → ChatRead` の二段カスケード（chat_reads に household_id FK なし）
  - 同俯瞰図のカスケード表に **UserCategoryPreference が Household 行に未記載**（Category 経由で間接 CASCADE される事実が読み取りづらい）
- **問い**: 以下を集約俯瞰図に反映してよいか？
  - **(A)** ChatRead は ChatMessage 経由として点線で表記
  - **(B)** UCP を Household 行に「（Category 経由）」と注釈付きで追記
- **影響**: ER と俯瞰図の整合、家計簿削除時のカスケード範囲の正確な把握

### Q-H10: UCP の Lazy 作成時に Household 整合性をアプリ層で必ず検査することを業務ルールに明文化してよいか？

- **重要度**: High
- **状況**: [er.dbml](../er.dbml): UCP は `(user_id, category_id)` のみ。User が現在所属する家計簿のカテゴリ以外への INSERT が DB レベルで弾けない（アプリ層チェックだけ）
- **問い**: [user-category-preference 業務ルール](../domains/user-category-preference/業務ルール.md) に「UCP 作成時は user_id の所属 Household と category_id の所属 Household の一致をアプリ層で検査する」を追記してよいか？
- **影響**: 業務ルールの網羅性、実装時の検査漏れ防止

---

## Medium（詳細詰め。実装中に解決可）

### Q-M1: 退会済みユーザーの ONLINE 表示の扱い

- **状況**: [chat-message 業務ルール](../domains/chat-message/業務ルール.md): OnlineStatus は退会で削除されない（v1 ではそのまま残す）。LastSeenAt が直前で「5 分以内」だと一瞬 ONLINE 表示され得る
- **問い**: 退会済みユーザーをメンバー一覧の ONLINE 判定から除外するか？それともそのまま判定するか？

### Q-M2: 固定費月初コピーで所有者が退会済みになっていた場合の挙動 ✅

- **状況**: [entry 業務ルール](../domains/entry/業務ルール.md): 固定費月初コピーバッチで前月明細の所有者をそのままコピー。退会済み所有者の前月明細をコピーすると、Active メンバーシップ判定で弾かれる可能性
- **問い**:
  - **(A)** 退会済み所有者の前月明細はスキップ（コピーしない）
  - **(B)** バッチでは Active 判定を緩めて、退会済み所有者のままコピー継続
- **決定 (2026-05-13)**: **(A) 採用（Q-C5 と連動）。**バッチの「Active な Membership 1 件ごとにループ」の段階で退会済みユーザーは元から含まれないため、明示スキップ処理は不要。業務ルール本文にも明示済み。

### Q-M3: Push 通知本文テンプレートの確定

- **状況**: [chat-message 業務ルール](../domains/chat-message/業務ルール.md): 「送信者のニックネーム + メッセージ本文の冒頭（実装で確定）」と曖昧
- **問い**: 通知本文を「{ニックネーム}: {本文先頭 N 文字}」とする場合の N（例: 60 文字）、改行処理（スペース置換 / 省略）、絵文字の扱いを実装着手前に確定してよいか？

### Q-M4: カテゴリ管理画面 S13 の Member ユーザー閲覧権限

- **状況**: [画面一覧 S13](../画面一覧.md): カテゴリ管理画面を「Member / Admin」認可とし、Member は「家計簿カテゴリタブ」を**閲覧のみ可**と読める。明示なし
- **問い**: Member ユーザーは「家計簿カテゴリ一覧」を閲覧できるか？それとも管理者のみ閲覧可とするか？

### Q-M5: CSV インポートの 1 ファイル最大件数・API タイムアウト目安

- **状況**: [entry 業務ルール](../domains/entry/業務ルール.md): 件数制限なしと記載。家計簿移行で数千行流す想定なら API タイムアウト・プレビュー性能の前提が必要
- **問い**: v1 想定で 1 回あたり最大何件まで想定するか？（例: 5000 行）API タイムアウト超過時の挙動は？

### Q-M6: ChatMessage 編集・削除「v1 Out of Scope」を要件 §1.3 にも明記してよいか？

- **状況**: [chat-message 業務ルール](../domains/chat-message/業務ルール.md) に明記されているが、[要件定義書 §1.3](../../要件定義書.md) には記載なし
- **問い**: 要件 §1.3 Out of Scope に「チャットメッセージの編集・削除（v1 では送信のみ）」を追記してよいか？

### Q-M7: Entry の所有者「編集不可、削除して再登録」を要件側に明記してよいか？

- **状況**: [entry 業務ルール](../domains/entry/業務ルール.md): OwnerUserID は作成後変更不可と明確。[要件定義書 §4.1 収支登録](../../要件定義書.md) には記載なし
- **問い**: 要件 §4.1 に「所有者は登録時のみ指定可、編集での変更は不可（削除して再登録）」を追記してよいか？

### Q-M8: 退会画面 S12 の認可レベル

- **状況**: [画面一覧 S12](../画面一覧.md): 「Member」認可。`/onboarding/household` で家計簿に未参加のユーザーが退会したい場合の動線が無い
- **問い**:
  - **(A)** Member（家計簿所属者）のみ退会可。未所属ユーザーは退会不可
  - **(B)** Authed（認証済み）に変更し、未所属ユーザーも退会可能にする

### Q-M9: Entry の Memo「0 文字 / NULL / 空文字」の扱い

- **状況**: [er.dbml](../er.dbml): `memo varchar` で NULL 許容。[entry ドメインモデル](../domains/entry/ドメインモデル.md): 「0〜200 文字」と整数範囲表現
- **問い**: 「Memo なし」を NULL のみとするか、空文字も許容するか？API 側で正規化（空文字 → NULL）するか？

### Q-M10: 視覚化 API（分析ダッシュボード 5 グラフ）の認可・件数制限

- **状況**: [entry ドメインモデル](../domains/entry/ドメインモデル.md): 視覚化 5 クエリは記載があるが、[機能一覧](../機能一覧.md) では参照のみで API 仕様（認可・取得期間上限）の業務ルール記載なし
- **問い**: 視覚化 API の認可（Active メンバー限定）と取得期間上限（過去 12 ヶ月 / 24 ヶ月など）を業務ルール化してよいか？

---

## Minor サマリ（参考）

以下は記述揺れレベル。実装にはほぼ影響しないが、ドキュメント整流時に修正候補:

- 対訳表で「チャット | Chat」と「チャットメッセージ | ChatMessage」が併存。README 表記「チャット (ChatMessage)」が ChatMessage 集約のみを指すように読めるが、実際は 3 集約（ChatMessage / ChatRead / OnlineStatus）含む
- 業務フロー図 M2「LINE ログイン・プロフィール登録」が一塊化されているが、要件・UC では別操作（#1, #2）
- 機能一覧の合計「26」と実行者別サマリ「4 + 2 + 9 + 10 + 1 = 26」は一致 ✓
- 集約俯瞰図の凡例 `o..`（ID 参照） / `*--`（コンポジション）が各ドメインモデル mermaid 図でも一貫使用されている ✓

---

## 総評

全体として **9 割方の整合性が取れており、実装着手前の最終整流レベル**に達している。Critical 5 件のうち以下 2 件は実装着手前に必ずクローズする:

1. **Q-C2（CASCADE 一貫性）** — マイグレーション設計が割れる
2. **Q-C3（UCP DisplayOrder 一意性）** — トランザクション設計が割れる

残る Critical 3 件（Q-C1 / Q-C4 / Q-C5）は実装と並行して合意可能だが、できれば最初のスプリント開始前に確定したい。

High / Medium は要件定義書・対訳表・業務ルールへの補記が中心で、実装イテレーション中の解決で問題ない。

---

## バージョン履歴

| 日付 | 版 | 内容 |
|------|---|------|
| 2026-05-13 | 0.1 | 初版作成。document-reviewer による設計ドキュメント全体整合性レビューに基づく顧客 QA 叩き台（Critical 5 / High 10 / Medium 10）|
