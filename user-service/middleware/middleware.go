package middleware

import (
	"bookstore-framework/configs"
	"bookstore-framework/pkg"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth() gin.HandlerFunc {
	cfg, err := configs.LoadConfig()
	if err != nil {
		panic("error when load config")
	}
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			pkg.UnauthorizedResponse(ctx)
			ctx.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			pkg.ErrorResponse(ctx, http.StatusUnauthorized, "invalid authorization format", nil)
			ctx.Abort()
			return
		}

		tokenString := parts[1]

		token, err := jwt.ParseWithClaims(
			tokenString,
			&pkg.Claims{},
			func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(cfg.SecretKey), nil
			},
		)

		if err != nil {
			pkg.ErrorResponse(ctx, http.StatusUnauthorized, "invalid or expired token", err.Error())
			ctx.Abort()
			return
		}

		if claims, ok := token.Claims.(*pkg.Claims); ok && token.Valid {
			ctx.Set("userID", claims.UserID)
			ctx.Set("username", claims.Username)
			ctx.Set("email", claims.Email)
			ctx.Next()
		} else {
			pkg.UnauthorizedResponse(ctx)
			ctx.Abort()
			return
		}
	}
}
