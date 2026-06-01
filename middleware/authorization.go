package middleware

import (
	"slices"

	"github.com/parikhrahil/go-micro-starter/dto"

	"github.com/gin-gonic/gin"
)

func AuthorizationHandler(allowedRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		interfaces, exists := ctx.Get("userRoles")
		if !exists {
			apperror := dto.ErrUnauthorized("Authentication role context missing")
			dto.Abort(ctx, apperror.Status, apperror.Code, apperror.Message)
			return
		}

		userRoles, ok := interfaces.([]string)
		if !ok {
			apperror := dto.ErrUnauthorized("Invalid role context format")
			dto.Abort(ctx, apperror.Status, apperror.Code, apperror.Message)
			return
		}

		hasAccess := false
		for _, allowedRole := range allowedRoles {
			if slices.Contains(userRoles, allowedRole) {
				hasAccess = true
				break
			}
		}

		if !hasAccess {
			apperror := dto.ErrUnauthorized("You do not have permission to access this resource")
			dto.Abort(ctx, apperror.Status, apperror.Code, apperror.Message)
			return
		}

		ctx.Next()
	}
}
