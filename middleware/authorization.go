package middleware

import (
	"github.com/parikhrahil/go-micro-starter/dto"

	"github.com/gin-gonic/gin"
)

func AuthorizationHandler(allowedRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		interfaces, exists := ctx.Get("userRoles")
		if !exists {
			apperror := dto.ErrUnauthorized("Authentication role context missing")
			ctx.AbortWithStatusJSON(apperror.Status, gin.H{
				"code":    apperror.Code,
				"message": apperror.Message,
			})
			return
		}

		userRoles, ok := interfaces.([]string)
		if !ok {
			apperror := dto.ErrUnauthorized("Invalid role context format")
			ctx.AbortWithStatusJSON(apperror.Status, gin.H{
				"code":    apperror.Code,
				"message": apperror.Message,
			})
			return
		}

		hasAccess := false
		for _, allowedRole := range allowedRoles {
			for _, userRole := range userRoles {
				if allowedRole == userRole {
					hasAccess = true
					break
				}
			}
		}

		if !hasAccess {
			apperror := dto.ErrUnauthorized("You do not have permission to access this resource")
			ctx.AbortWithStatusJSON(apperror.Status, gin.H{
				"code":    apperror.Code,
				"message": apperror.Message,
			})
			return
		}

		ctx.Next()
	}
}
