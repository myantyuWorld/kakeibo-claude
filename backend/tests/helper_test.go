package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func truncateTables(t *testing.T) {
	t.Helper()

	rows, err := db.QueryContext(context.Background(),
		`SELECT tablename FROM pg_tables WHERE schemaname = 'public'`,
	)
	require.NoError(t, err)
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var name string
		require.NoError(t, rows.Scan(&name))
		tables = append(tables, name)
	}
	require.NoError(t, rows.Err())

	if len(tables) == 0 {
		return
	}

	query := "TRUNCATE TABLE " + strings.Join(tables, ", ") + " CASCADE"

	_, err = db.ExecContext(context.Background(), query)
	require.NoError(t, err)
}

type request struct {
	t   *testing.T
	req *http.Request
}

// newRequest は server に対する HTTP リクエストを構築する。
// body が nil でない場合は JSON エンコードしてリクエストボディに設定する。
func newRequest(t *testing.T, method, path string, body any) *request {
	t.Helper()

	var req *http.Request
	if body != nil {
		b, err := json.Marshal(body)
		require.NoError(t, err)
		req = httptest.NewRequestWithContext(t.Context(), method, path, bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequestWithContext(t.Context(), method, path, nil)
	}

	return &request{t: t, req: req}
}

// withAuth は Authorization ヘッダにダミートークンを付与する。
// 認証基盤の導入後に実装を調整する。
func (r *request) withAuth() *request {
	r.req.Header.Set("Authorization", "Bearer test-dummy-token")
	return r
}

func (r *request) do() *httptest.ResponseRecorder {
	r.t.Helper()
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, r.req)
	return rec
}

func parseBody(t *testing.T, rec *httptest.ResponseRecorder, dest any) {
	t.Helper()
	require.NoError(t, json.NewDecoder(rec.Body).Decode(dest))
}
