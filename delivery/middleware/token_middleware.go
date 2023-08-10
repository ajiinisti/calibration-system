package middleware

import (
	"net/http"

	"calibration-system.com/utils/authenticator"
	"github.com/gin-gonic/gin"
)

type authHeader struct {
	AuthorizationHeader string `header:"Authorization"` // -> key dari header
}

type AuthTokenMiddleware interface {
	RequireToken() gin.HandlerFunc
}

type authTokenMiddleware struct {
	accToken authenticator.AccessToken
}

func (a *authTokenMiddleware) RequireToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString, err := authenticator.BindAuthHeader(ctx)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			ctx.Abort()
			return
		}
		accessDetail, err := a.accToken.VerifyAccessToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			ctx.Abort()
			return
		}
		_, err = a.accToken.FetchAccessToken(accessDetail.AccessUUID)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			ctx.Abort()
			return
		}
		ctx.Set("Roles", accessDetail.Roles)
		ctx.Set("Email", accessDetail.Email)
		ctx.Set("ID", accessDetail.ID)
		ctx.Next()
	}
}

func NewTokenValidator(acctToken authenticator.AccessToken) AuthTokenMiddleware {
	return &authTokenMiddleware{
		accToken: acctToken,
	}
}
