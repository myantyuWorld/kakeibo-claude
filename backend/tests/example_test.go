// 統合テストのテンプレート。新規ドメインのテスト作成時はこのファイルを参考にする。
package tests

import (
	template "backend/internal/_template"
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExample(t *testing.T) {
	t.Skip("サンプルのためスキップ")

	// テスト間のデータ干渉を防ぐため、各テスト関数の先頭で必ず呼ぶ
	truncateTables(t)

	// テストデータは API 経由ではなく SQL で直接投入する。
	// API の実装状況に依存せず、Arrange を安定させるため。
	_, err := db.ExecContext(context.Background(),
		`INSERT INTO customers (id, name, email) VALUES ($1, $2, $3)`,
		"00000000-0000-0000-0000-000000000001", "テスト顧客", "test@example.com",
	)
	require.NoError(t, err)

	// 統合テストは正常系のみ。異常系・網羅は単体テスト（internal/）で担保する。
	t.Run("顧客を1件取得できること", func(t *testing.T) {
		rec := newRequest(t, http.MethodGet, "/customers/00000000-0000-0000-0000-000000000001", nil).do()

		assert.Equal(t, http.StatusOK, rec.Code)

		var body template.GetResponse
		parseBody(t, rec, &body)
		assert.Equal(t, "テスト顧客", body.Name)
	})

	t.Run("認証付きで顧客を作成できること", func(t *testing.T) {
		rec := newRequest(t, http.MethodPost, "/customers", template.CreateRequest{
			Name:   "新規顧客",
			Status: "example",
		}).withAuth().do()

		assert.Equal(t, http.StatusCreated, rec.Code)
	})
}
