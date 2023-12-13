package controller

import (
	"fmt"
	"net/http"

	"calibration-system.com/config"
	"calibration-system.com/delivery/api"
	"calibration-system.com/delivery/api/request"
	"calibration-system.com/delivery/api/response"
	"calibration-system.com/delivery/middleware"
	"calibration-system.com/model"
	"calibration-system.com/usecase"
	"calibration-system.com/utils"
	"calibration-system.com/utils/authenticator"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

	var roles []string
	for _, v := range user.Roles {
		roles = append(roles, v.Name)
	}

	cred := model.TokenModel{
		Username: "",
		Email:    user.Email,
		Role:     roles,
		ID:       user.ID,
		Name:     user.Name,
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
	code := ctx.Query("code")

	if code == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Authorization code not provided!"})
		return
	}

	tokenRes, err := utils.GetGoogleOauthToken(a.cfg, code)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail token", "message": err.Error()})
		return
	}

	userGoogle, err := utils.GetGoogleUser(tokenRes.Access_token, tokenRes.Id_token)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail google user", "message": err.Error()})
	}
	user, err := a.uc.GetUserByEmail(userGoogle.Email)
	if err == gorm.ErrRecordNotFound {
		a.NewFailedResponse(ctx, http.StatusInternalServerError, fmt.Sprint("Email/Password invalid"))
		return
	}
	if err != nil {
		a.NewFailedResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	var roles []string
	for _, v := range user.Roles {
		roles = append(roles, v.Name)
	}

	cred := model.TokenModel{
		Username: "",
		Email:    user.Email,
		Role:     roles,
		ID:       user.ID,
	}

	tokenDetail, err := a.tokenService.CreateAccessToken(&cred)
	if err != nil {
		a.NewFailedResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	if err := a.tokenService.StoreAccessToken(user.Email, tokenDetail); err != nil {
		a.NewFailedResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	// redis add token
	response := response.LoginResponse{
		AccessToken: tokenDetail.AccessToken,
		TokenModel:  cred,
	}
	// set cookie and redirect to oauth url
	a.NewSuccessSingleResponse(ctx, response, "OK")
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

func (a *AuthController) forgetPassword(ctx *gin.Context) {
	var userCredential *request.ForgotPasswordInput

	if err := ctx.ShouldBindJSON(&userCredential); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	message := "You will receive a reset email if user with that email exist"

	user, err := a.uc.GetUserByEmail(userCredential.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			a.NewFailedResponse(ctx, http.StatusInternalServerError, fmt.Sprint("Your email invalid"))
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error get user", "message": err.Error()})
	}

	// Generate Verification Code
	resetToken, err := utils.GeneratePassword()
	err = a.uc.ForgetPassword(userCredential.Email, resetToken)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			a.NewFailedResponse(ctx, http.StatusInternalServerError, fmt.Sprint("Your email invalid"))
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error forgot password", "message": err.Error()})
	}

	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"status": "success", "message": err.Error()})
		return
	}

	emailData := utils.EmailData{
		URL:       "https://calibration.techconnect.co.id/#/reset-password/" + fmt.Sprintf("%s/%s", user.Email, resetToken),
		FirstName: user.Email,
		Subject:   "Your password reset token (valid for 10min)",
	}

	// utils/templates/resetPassword.html
	err = utils.SendMail([]string{user.Email}, &emailData, "./utils/templates", "resetPassword.html", a.cfg.SMTPConfig)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
}

func (a *AuthController) resetPassword(ctx *gin.Context) {
	resetToken := ctx.Params.ByName("resetToken")
	email := ctx.Params.ByName("email")

	var userCredential *request.ResetPasswordInput

	if err := ctx.ShouldBindJSON(&userCredential); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if userCredential.Password != userCredential.ConfirmPassword {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Passwords do not match"})
		return
	}

	if err := a.uc.ResetPassword(email, resetToken, userCredential.Password, userCredential.ConfirmPassword); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	ctx.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "", -1, "/", "localhost", false, true)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Password data updated successfully"})
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
	r.PATCH("/reset-password/:email/:resetToken", controller.resetPassword)
	auth.POST("/change-password", controller.changePassword)
	auth.POST("/logout", controller.logout)
	return &controller
}
