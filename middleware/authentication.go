package middleware

import (
	"errors"
	"strings"

	"github.com/parikhrahil/go-micro-starter/dto"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Claims defines the structure of the data inside your JWT
type Claims struct {
	ID    string   `json:"id"`
	Roles []string `json:"roles"`
	jwt.RegisteredClaims
}

type AuthenticationOpts struct {
	JwtSecret string
}

// AuthenticationHandler validates the JWT token and injects claims into the Gin context
func AuthenticationHandler(opts *AuthenticationOpts) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			apperror := dto.ErrUnauthenticated("Authorization header is required")
			c.AbortWithStatusJSON(apperror.Status, gin.H{
				"code":    apperror.Code,
				"message": apperror.Message,
			})
			return
		}

		// Expecting headers format: "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			apperror := dto.ErrUnauthenticated("Authorization header must be Bearer token")
			c.AbortWithStatusJSON(apperror.Status, gin.H{
				"code":    apperror.Code,
				"message": apperror.Message,
			})
			return
		}

		tokenString := parts[1]

		var claims Claims

		// Parse and validate the token
		token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (any, error) {
			// Validate the signing method is HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(opts.JwtSecret), nil
		})

		if err != nil || token == nil || !token.Valid {
			apperror := dto.ErrUnauthenticated("Invalid or expired token")
			c.AbortWithStatusJSON(apperror.Status, gin.H{
				"code":    apperror.Code,
				"message": apperror.Message,
			})
			return
		}

		// Pass user details downstream into the Gin context
		c.Set("userID", claims.ID)
		c.Set("userRoles", claims.Roles)

		c.Next()
	}
}
