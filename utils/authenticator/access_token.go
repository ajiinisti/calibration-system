package authenticator

import (
	"context"
	"errors"
	"fmt"
	"time"

	"calibration-system.com/config"
	"calibration-system.com/model"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type AccessToken interface {
	CreateAccessToken(cred *model.TokenModel) (TokenDetail, error)
	VerifyAccessToken(tokenStr string) (AccessDetail, error)
	StoreAccessToken(username string, tokenDetail TokenDetail) error
	FetchAccessToken(token string) (string, error)
	DeleteAccessToken(accessUUID string) error
}

type accessToken struct {
	Config config.Config
	client *redis.Client
}

func (t *accessToken) CreateAccessToken(cred *model.TokenModel) (TokenDetail, error) {
	tokenDetail := TokenDetail{}
	tokenDetail.AccessUUID = uuid.New().String()
	tokenDetail.AtExpired = time.Now().Add(t.Config.AccessTokenLifeTime).Unix()
	claims := MyClaims{
		StandardClaims: jwt.StandardClaims{Issuer: t.Config.ApplicationName},
		TokenModel: model.TokenModel{
			Username: cred.Username,
			Email:    cred.Email,
			Role:     cred.Role,
			ID:       cred.ID,
		},
		AccessUUID: tokenDetail.AccessUUID,
	}
	now := time.Now().UTC()
	end := now.Add(t.Config.AccessTokenLifeTime)
	claims.IssuedAt = now.Unix()
	claims.ExpiresAt = end.Unix()
	token := jwt.NewWithClaims(
		t.Config.JwtSigningMethod,
		claims,
	)

	newToken, err := token.SignedString([]byte(t.Config.JwtSignatureKey))
	if err != nil {
		return TokenDetail{}, err
	}
	tokenDetail.AccessToken = newToken
	return tokenDetail, nil
}

func (t *accessToken) VerifyAccessToken(tokenString string) (AccessDetail, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Signing method invalid")
		} else if method != t.Config.JwtSigningMethod {
			return nil, fmt.Errorf("Signing method invalid")
		}
		return []byte(t.Config.JwtSignatureKey), nil
	})
	if err != nil {
		return AccessDetail{}, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid || claims["iss"] != t.Config.ApplicationName {
		return AccessDetail{}, err
	}
	email := claims["Email"].(string)
	uuid := claims["AccessUUID"].(string)
	role := claims["Role"].(string)
	id := claims["ID"].(string)
	return AccessDetail{
		AccessUUID: uuid,
		Email:      email,
		Roles:      role,
		ID:         id,
	}, nil
}

func (a *accessToken) DeleteAccessToken(accessUUID string) error {
	rowsAffected, err := a.client.Del(context.Background(), accessUUID).Result()
	if err != nil {
		return err
	} else if rowsAffected == 0 {
		return fmt.Errorf("accessUUID : %s Not Found", accessUUID)
	}
	return nil
}

func (a *accessToken) StoreAccessToken(email string, tokenDetail TokenDetail) error {
	at := time.Unix(tokenDetail.AtExpired, 0)
	err := a.client.Set(context.Background(), tokenDetail.AccessUUID, email, at.Sub(time.Now())).Err()
	if err != nil {
		return err
	}
	return nil
}

func (a *accessToken) FetchAccessToken(token string) (string, error) {
	email, err := a.client.Get(context.Background(), token).Result()
	if err != nil {
		return "", err
	}
	if email == "" {
		return "", errors.New("Invalid Token")
	}
	return email, nil
}

func NewTokenService(config config.Config, client *redis.Client) AccessToken {
	return &accessToken{
		Config: config,
		client: client,
	}
}
