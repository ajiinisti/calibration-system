package controller

import (
	"net/http"
	"strconv"
	"strings"

	"calibration-system.com/delivery/api"
	"calibration-system.com/delivery/api/request"
	"calibration-system.com/delivery/middleware"
	"calibration-system.com/model"
	"calibration-system.com/usecase"
	"calibration-system.com/utils/authenticator"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	router       *gin.Engine
	uc           usecase.UserUsecase
	tokenService authenticator.AccessToken
	api.BaseApi
}

func (u *UserController) listHandler(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		u.NewFailedResponse(c, http.StatusBadRequest, "Invalid page number")
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		u.NewFailedResponse(c, http.StatusBadRequest, "Invalid limit number")
	}

	users, pagination, err := u.uc.FindPagination(request.PaginationParam{
		Page:   page,
		Limit:  limit,
		Offset: 0,
	})
	if err != nil {
		u.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var newUsers []interface{}
	for _, v := range users {
		newUsers = append(newUsers, v)
	}

	u.NewSuccesPagedResponse(c, newUsers, "OK", pagination)
}

func (u *UserController) allListHandler(c *gin.Context) {
	userList, err := u.uc.FindAll()
	if err != nil {
		u.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	u.NewSuccessSingleResponse(c, userList, "OK")
}

func (u *UserController) getByIdHandler(c *gin.Context) {
	id := c.Param("id")
	roles, err := u.uc.FindById(id)
	if err != nil {
		u.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	u.NewSuccessSingleResponse(c, roles, "OK")
}

func (u *UserController) getByIdSwitchHandler(c *gin.Context) {
	id := c.Param("id")
	roles, err := u.uc.FindByIdSwitchUser(id)
	if err != nil {
		u.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	u.NewSuccessSingleResponse(c, roles, "OK")
}

func (u *UserController) getByProjectId(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		u.NewFailedResponse(c, http.StatusBadRequest, "Invalid page number")
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		u.NewFailedResponse(c, http.StatusBadRequest, "Invalid limit number")
	}

	projectId := c.Param("projectId")
	users, pagination, err := u.uc.FindByProjectIdPagination(request.PaginationParam{
		Page:   page,
		Limit:  limit,
		Offset: 0,
	}, projectId)

	if err != nil {
		u.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var newUsers []interface{}
	for _, v := range users {
		newUsers = append(newUsers, v)
	}

	u.NewSuccesPagedResponse(c, newUsers, "OK", pagination)
}

func (u *UserController) createHandler(c *gin.Context) {
	var payload request.CreateUser

	if err := c.ShouldBind(&payload); err != nil {
		u.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user := model.User{
		Email:            payload.Email,
		Name:             payload.Name,
		Nik:              payload.Nik,
		DateOfBirth:      payload.DateOfBirth,
		SupervisorNik:    payload.SupervisorNik,
		BusinessUnitId:   &payload.BusinessUnitId,
		OrganizationUnit: payload.OrganizationUnit,
		Division:         payload.Division,
		Department:       payload.Department,
		JoinDate:         payload.JoinDate,
		Grade:            payload.Grade,
		HRBP:             payload.HRBP,
		Position:         payload.Position,
		ScoringMethod:    payload.ScoringMethod,
	}
	if err := u.uc.CreateUser(user, payload.Roles); err != nil {
		u.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	u.NewSuccessSingleResponse(c, "", "OK")
}

func (u *UserController) updateHandler(c *gin.Context) {
	var payload request.UpdateUser

	if err := c.ShouldBind(&payload); err != nil {
		u.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user := model.User{
		BaseModel: model.BaseModel{
			ID: payload.ID,
		},
		Email:            payload.Email,
		Name:             payload.Name,
		Nik:              payload.Nik,
		DateOfBirth:      payload.DateOfBirth,
		SupervisorNik:    payload.SupervisorNik,
		BusinessUnit:     model.BusinessUnit{},
		BusinessUnitId:   &payload.BusinessUnitId,
		OrganizationUnit: payload.OrganizationUnit,
		Division:         payload.Division,
		Department:       payload.Department,
		JoinDate:         payload.JoinDate,
		Grade:            payload.Grade,
		HRBP:             payload.HRBP,
		Position:         payload.Position,
		ScoringMethod:    payload.ScoringMethod,
	}
	if err := u.uc.SaveUser(user, payload.Roles); err != nil {
		u.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	u.NewSuccessSingleResponse(c, payload, "OK")
}

func (u *UserController) deleteHandler(c *gin.Context) {
	id := c.Param("id")
	if err := u.uc.DeleteData(id); err != nil {
		u.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusNoContent, "")
}

func (u *UserController) uploadHandler(c *gin.Context) {
	file, err := c.FormFile("excelFile")
	if err != nil {
		u.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logs, err := u.uc.BulkInsert(file)
	if err != nil {
		u.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if len(logs) > 0 {
		u.NewSuccessSingleResponse(c, strings.Join(logs, "."), "Success with some error")
		return
	}

	u.NewSuccessSingleResponse(c, "", "OK")
}

func (u *UserController) generatePasswordHandler(c *gin.Context) {
	id := c.Param("id")
	if err := u.uc.GeneratePasswordById(id); err != nil {
		u.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	u.NewSuccessSingleResponse(c, "Success Generate Password", "OK")
}

func NewUserController(u *gin.Engine, tokenService authenticator.AccessToken, uc usecase.UserUsecase) *UserController {
	controller := UserController{
		router:       u,
		tokenService: tokenService,
		uc:           uc,
	}

	auth := u.Use(middleware.NewTokenValidator(tokenService).RequireToken())
	auth.GET("/users", controller.listHandler)
	auth.GET("/users/all", controller.allListHandler)
	auth.GET("/users/:id", controller.getByIdHandler)
	auth.GET("/users-switch/:id", controller.getByIdSwitchHandler)
	auth.GET("/users/project/:projectId", controller.getByProjectId)
	auth.PUT("/users", controller.updateHandler)
	auth.POST("/users", controller.createHandler)
	auth.POST("/users/generate-password/:id", controller.generatePasswordHandler)
	auth.POST("/users/upload", controller.uploadHandler)
	auth.DELETE("/users/:id", controller.deleteHandler)
	return &controller
}
