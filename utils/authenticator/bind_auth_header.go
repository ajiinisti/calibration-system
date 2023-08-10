package authenticator

import (
	"errors"
	"net/http"
	"strings"

	request "calibration-system.com/delivery/api/request"
	"github.com/gin-gonic/gin"
)

func BindAuthHeader(ctx *gin.Context) (string, error) {
	var h request.AuthHeader

	if err := ctx.ShouldBindHeader(&h); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
		ctx.Abort()
		return "", err
	}

	tokenString := strings.Replace(h.AuthorizationHeader, "Bearer ", "", -1)
	if tokenString == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		ctx.Abort()
		return "", errors.New("token emmpty")
	}
	return tokenString, nil
}
