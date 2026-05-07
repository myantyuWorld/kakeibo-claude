package api

import (
	"backend/pkg/errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidator_Validate(t *testing.T) {
	v := newValidator()

	type request struct {
		Name  string `json:"name"  validate:"required"`
		Email string `json:"email" validate:"required,email"`
		Age   int    `json:"age"   validate:"min=0,max=150"`
	}

	tests := []struct {
		name    string
		input   request
		wantErr bool
	}{
		{"正常値", request{Name: "テスト", Email: "test@example.com", Age: 20}, false},
		{"required違反", request{Name: "", Email: "test@example.com", Age: 20}, true},
		{"email形式違反", request{Name: "テスト", Email: "invalid", Age: 20}, true},
		{"min違反", request{Name: "テスト", Email: "test@example.com", Age: -1}, true},
		{"max違反", request{Name: "テスト", Email: "test@example.com", Age: 200}, true},
		{"複数フィールド違反", request{Name: "", Email: "invalid", Age: -1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Validate(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				assert.True(t, errors.IsInvalid(err))
				var appErr *errors.Error
				require.True(t, errors.As(err, &appErr))
				assert.Equal(t, "入力内容に誤りがあります", appErr.Message())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidator_Validate_Unwrap(t *testing.T) {
	v := newValidator()

	type request struct {
		Name string `json:"name" validate:"required"`
	}

	err := v.Validate(request{Name: ""})
	require.Error(t, err)

	var appErr *errors.Error
	require.True(t, errors.As(err, &appErr))
	assert.Error(t, appErr.Unwrap(), "原因エラーとして ValidationErrors を保持する")
}
