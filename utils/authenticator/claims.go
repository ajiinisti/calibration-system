package authenticator

import (
	"calibration-system.com/model"
	"github.com/golang-jwt/jwt"
)

type MyClaims struct {
	jwt.StandardClaims
	model.TokenModel
	AccessUUID string
}
