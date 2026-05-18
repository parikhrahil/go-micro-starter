package middleware

import (
	"errors"
	"net/http"

	"github.com/parikhrahil/go-micro-starter/dto"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if len(ctx.Errors) == 0 {
			return
		}

		err := ctx.Errors.Last().Err

		if appErr, ok := errors.AsType[*dto.AppError](err); ok {
			dto.Fail(ctx, appErr.Status, appErr.Code, appErr.Message)
		} else {
			dto.Fail(ctx, http.StatusInternalServerError, "INTERNAL", "an unexpected error occured")
		}
	}
}
