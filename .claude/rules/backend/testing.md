---
paths:
  - "backend/**/*.go"
  - "backend/tests/**"
---

# テスト方針

## 単体テスト

- `internal/{aggregate}/` 内に `_test.go` を配置する（同一パッケージ）
- ALWAYS テーブル駆動テストで記述する
- ALWAYS testify を使用する（assert / require / mock）
- ALWAYS 受入基準の Given/When/Then とテストケースを対応させる
- テスト実行: `mise run test`

```go
func Test_customer_delete(t *testing.T) {
    tests := []struct {
        name    string
        status  status
        wantErr bool
    }{
        {"Active な顧客を削除できる", active, false},
        {"Deleted な顧客は削除できない", deleted, true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            c := &customer{status: tt.status}
            err := c.delete()
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, deleted, c.status)
            }
        })
    }
}
```

## 統合テスト

- `tests/` にドメインごとのファイルを配置する（`tests/customer_test.go`）
- AAAパターン（Arrange / Act / Assert）で記述する
- ALWAYS 正常系のみをテストする（異常系・網羅は単体テストで担保）
- テスト用 DB（compose.yml の test-db）を使用する
- テスト実行: `mise run test-integration`

## 関連

- **go-test-reviewer** agent -- テストレビュー時
- skill: `generate-api-tests` -- テストコード生成時

