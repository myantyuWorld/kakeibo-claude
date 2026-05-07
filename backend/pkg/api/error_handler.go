package api

import (
	"backend/pkg/errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type errorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func errorHandler(err error, c echo.Context) {
	var appErr *errors.Error
	if errors.As(err, &appErr) {
		_ = c.JSON(httpStatus(appErr.Code()), errorResponse{
			Code:    string(appErr.Code()),
			Message: appErr.Message(),
		})
		return
	}
	_ = c.JSON(http.StatusInternalServerError, errorResponse{
		Code:    string(errors.ErrInternal),
		Message: "an unexpected error occurred",
	})
}

func httpStatus(code errors.ErrCode) int {
	switch code {
	case errors.ErrNotFound:
		return http.StatusNotFound
	case errors.ErrUnauthorized:
		return http.StatusUnauthorized
	case errors.ErrForbidden:
		return http.StatusForbidden
	default:
		return http.StatusBadRequest
	}
}
