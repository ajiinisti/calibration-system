package controller

import (
	"net/http"

	"calibration-system.com/config"
	"calibration-system.com/delivery/api"
	"calibration-system.com/delivery/api/request"
	"calibration-system.com/delivery/api/response"
	"calibration-system.com/delivery/middleware"
	"calibration-system.com/model"
	"calibration-system.com/usecase"
	"calibration-system.com/utils/authenticator"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	router       *gin.Engine
	uc           usecase.AuthUsecase
	tokenService authenticator.AccessToken
	cfg          config.Config
	api.BaseApi
}

func (a *AuthController) login(c *gin.Context) {
	var payload request.Login

	if err := a.ParseRequestBody(c, &payload); err != nil {
		a.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	user, err := a.uc.Login(payload)
	if err != nil {
		a.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	cred := model.TokenModel{
		Email: user.Email,
		Role:  user.Role.Name,
		ID:    user.ID,
	}
	tokenDetail, err := a.tokenService.CreateAccessToken(&cred)
	if err != nil {
		a.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if err := a.tokenService.StoreAccessToken(user.Email, tokenDetail); err != nil {
		a.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	// redis add token
	response := response.LoginResponse{
		AccessToken: tokenDetail.AccessToken,
		TokenModel:  cred,
	}
	a.NewSuccessSingleResponse(c, response, "OK")
}

func (a *AuthController) redirectGoogle(c *gin.Context) {
	c.JSON(http.StatusOK, "ok")
}

func (a *AuthController) loginGoogle(ctx *gin.Context) {
	panic("unimplemented")
}

func (a *AuthController) logout(c *gin.Context) {
	token, err := authenticator.BindAuthHeader(c)
	if err != nil {
		c.AbortWithStatus(401)
	}

	accountDetail, err := a.tokenService.VerifyAccessToken(token)
	if err != nil {
		a.NewFailedResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	if err = a.tokenService.DeleteAccessToken(accountDetail.AccessUUID); err != nil {
		a.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(200, gin.H{
		"message": "Success Logout",
	})
}

func (a *AuthController) forgetPassword(c *gin.Context) {
	panic("unimplemented")
}

func (a *AuthController) changePassword(c *gin.Context) {
	panic("unimplemented")
}

func NewAuthController(r *gin.Engine, uc usecase.AuthUsecase, tokenService authenticator.AccessToken, cfg config.Config) *AuthController {
	controller := AuthController{
		router:       r,
		uc:           uc,
		tokenService: tokenService,
		cfg:          cfg,
	}

	auth := r.Group("/auth").Use(middleware.NewTokenValidator(tokenService).RequireToken())
	r.POST("/login", controller.login)
	r.POST("/forget-password", controller.forgetPassword)
	r.GET("/sessions/oauth/google", controller.redirectGoogle)
	r.GET("/sessions/oauth", controller.loginGoogle)
	auth.POST("/change-password", controller.changePassword)
	auth.POST("/logout", controller.logout)
	return &controller
}
