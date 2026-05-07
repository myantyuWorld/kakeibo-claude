---
name: generate-api-tests
description: 受入基準からバックエンドのテストコードを生成する。「APIのテストを書いて」「バックエンドのテストを生成して」など、バックエンドのテスト作成を依頼されたときに使用する。
---

# generate-api-tests

指定されたドメイン・ユースケースの受入基準からテストコードを生成する。

## When to Activate

- 「APIのテストを書いて」「バックエンドのテストを生成して」と依頼されたとき
- `impl-api` の Phase 4 から呼び出されたとき
- ドメインモデルやサービスの実装後にテストを追加するとき

## 入力

- ドメイン名（例: customer）
- 対象（例: ドメインモデル / ユースケース名 / 全体）

## 手順

### 1. 受入基準の確認

`docs/design/{ドメイン名}/業務ルール.md` の受入基準を確認する。Given/When/Then の各シナリオをテストケースに対応させる。

### 2. テストコードの生成

対象に応じてテストコードを生成する。

#### 単体テスト（internal/{ドメイン名}/ 内）

**ドメインモデルのテスト → `{ドメイン名}_test.go`**

バリデーション、状態遷移、ドメインロジックの検証。

```go
// PASS: 受入基準のシナリオ名をテストケース名に使う
func Test_customer_create(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"正常な名前で顧客を作成できる", "山田太郎", false},
        {"空の名前では作成できない", "", true},
        {"100文字の名前で作成できる", strings.Repeat("あ", 100), false},
        {"101文字の名前では作成できない", strings.Repeat("あ", 101), true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := newCustomer(tt.input)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}

// FAIL: テストケース名が受入基準と対応していない
func Test_customer_create(t *testing.T) {
    tests := []struct {
        name string
        // ...
    }{
        {"test1", ...},     // NG: 意味のない名前
        {"valid", ...},     // NG: 英語で曖昧
    }
}
```

**サービスのテスト → `service_test.go`**

repository を mock してビジネスロジックを検証する。

```go
// PASS: mock は必要最小限に設定する
func Test_service_createCustomer(t *testing.T) {
    repo := new(mockRepository)
    repo.On("Save", mock.Anything, mock.AnythingOfType("*customer")).
        Return(nil)
    svc := newService(repo)

    err := svc.Create(context.Background(), "山田太郎")

    assert.NoError(t, err)
    repo.AssertExpectations(t)
}

// FAIL: 検証しない mock を大量に設定する
func Test_service_createCustomer(t *testing.T) {
    repo := new(mockRepository)
    repo.On("FindByID", mock.Anything, mock.Anything).Return(nil, nil)  // 不要
    repo.On("FindAll", mock.Anything).Return(nil, nil)                   // 不要
    repo.On("Save", mock.Anything, mock.Anything).Return(nil)
    repo.On("Delete", mock.Anything, mock.Anything).Return(nil)          // 不要
}
```

#### 統合テスト（tests/ 内）

**API の正常系テスト → `tests/{ドメイン名}_test.go`**

AAAパターン（Arrange / Act / Assert）で記述。正常系のみ。

```go
// PASS: AAA パターンで正常系を検証
func TestCreateCustomer(t *testing.T) {
    // Arrange
    body := `{"name": "山田太郎"}`

    // Act
    resp := e.POST("/api/customers").
        WithBytes([]byte(body)).
        Expect()

    // Assert
    resp.Status(http.StatusCreated)
}
```

### 3. 受入基準との対応確認

生成後、以下の対応表を出力して漏れがないことを確認する。

```
| 受入基準 | テスト種別 | テストケース |
|---------|----------|------------|
| {シナリオ名} | 単体/統合 | {テスト関数名} |
```

## テストケース命名パターン

| パターン | テストケース名の例 |
|---------|-----------------|
| 正常系 | `正常な名前で顧客を作成できる` |
| バリデーション | `空の名前では作成できない` |
| 境界値 | `100文字の名前で作成できる` / `101文字の名前では作成できない` |
| 状態遷移 | `Active な顧客を削除できる` / `Deleted な顧客は削除できない` |
| 一意性 | `同じメールアドレスでは登録できない` |
| ドメイン間制約 | `未完了の注文がある顧客は削除できない` |

## 境界値テストの書き方

数値制約やフィールド長のテストでは、境界の内側・境界・境界の外側の3パターンを含める。

```go
// 名前が1〜100文字の場合
{"1文字の名前で作成できる", "あ", false},           // 境界（下限）
{"100文字の名前で作成できる", repeat("あ", 100), false}, // 境界（上限）
{"空の名前では作成できない", "", true},               // 境界の外側（下限）
{"101文字の名前では作成できない", repeat("あ", 101), true}, // 境界の外側（上限）
```

## テストコードの規約

- `.claude/rules/backend/testing.md` に従う
- 単体テスト: テーブル駆動テスト + testify
- 統合テスト: AAAパターン、正常系のみ
- テストケース名は受入基準のシナリオ名と対応させる
- `assert` は検証を続行する場合、`require` はテスト続行不可能な場合に使う

---

**Remember**: 受入基準の Given/When/Then がそのままテストケースになる。漏れなく対応させる。
