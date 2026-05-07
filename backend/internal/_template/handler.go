package template

import (
	"backend/pkg/errors"
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Serviceインターフェース
// - 型は exported、メソッドは unexported
// - パッケージ外からの実装は不可（unexported メソッドのため）
// - DI/ワイヤリング時の引数型として利用する
type Service interface {
	create(ctx context.Context, input CreateInput) error
	getByID(ctx context.Context, id uuid.UUID) (*GetResponse, error)
}

type handler struct {
	svc Service
}

func newHandler(s Service) *handler {
	return &handler{svc: s}
}

// リクエスト/レスポンス構造体
// - 統合テスト（tests/）から参照するため exported
// - json タグを付与する
// - 構造体の定義は使用するメソッドの直上に配置する

type CreateRequest struct {
	Name   string `json:"name"   validate:"required"`
	Status string `json:"status" validate:"required"`
}

func (h *handler) create(c echo.Context) error {
	var req CreateRequest
	if err := c.Bind(&req); err != nil {
		return errors.NewInvalid()
	}
	if err := c.Validate(&req); err != nil {
		return err
	}
	if err := h.svc.create(c.Request().Context(), CreateInput{
		Name:   req.Name,
		Status: req.Status,
	}); err != nil {
		return err
	}
	return c.NoContent(http.StatusCreated)
}

type GetResponse struct {
	ID      string         `json:"id"`
	Name    string         `json:"name"`
	Status  string         `json:"status"`
	Example ExampleSummary `json:"example"` // resolver 経由で取得した他ドメインの情報
}

func (h *handler) getByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errors.NewInvalid()
	}
	res, err := h.svc.getByID(c.Request().Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, res)
}
