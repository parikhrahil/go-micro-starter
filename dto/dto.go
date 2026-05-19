package dto

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool       `json:"success"`
	Data    any        `json:"data,omitempty"`
	Error   *ErrorInfo `json:"error,omitempty"`
	Meta    *Meta      `json:"meta,omitempty"`
}

type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Meta struct {
	Page       int `json:"page,omitempty"`
	PerPage    int `json:"per_page,omitempty"`
	Total      int `json:"total,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
}

type AppError struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return e.Message
}

func OK(ctx *gin.Context, data any) {
	success(ctx, http.StatusOK, data)
}

func Created(ctx *gin.Context, data any) {
	success(ctx, http.StatusCreated, data)
}

func NoContent(ctx *gin.Context, data any) {
	success(ctx, http.StatusNoContent, data)
}

func success(ctx *gin.Context, status int, data any) {
	ctx.JSON(status, Response{
		Success: true,
		Data:    data,
	})
}

func Fail(ctx *gin.Context, status int, code, message string) {
	ctx.JSON(status, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
		},
	})
}

func Abort(ctx *gin.Context, status int, code, message string) {
	ctx.AbortWithStatusJSON(status, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
		},
	})
}

func ErrNotFound(e string) *AppError {
	return error(http.StatusNotFound, "NOT_FOUND", e)
}

func ErrUnauthenticated(e string) *AppError {
	return error(http.StatusForbidden, "FORBIDDEN", e)
}

func ErrUnauthorized(e string) *AppError {
	return error(http.StatusUnauthorized, "UNAUTHORIZED", e)
}

func ErrBadRequest(e string) *AppError {
	return error(http.StatusBadRequest, "BAD_REQUEST", e)
}

func ErrorInternalServer(e string) *AppError {
	return error(http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", e)
}

func ErrorConflict(e string) *AppError {
	return error(http.StatusConflict, "CONFLICT", e)
}

func error(status int, code, message string) *AppError {
	return &AppError{
		Status:  status,
		Code:    code,
		Message: message,
	}
}
