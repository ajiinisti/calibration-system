package controller

import (
	"net/http"

	"calibration-system.com/delivery/api"
	"calibration-system.com/delivery/middleware"
	"calibration-system.com/model"
	"calibration-system.com/usecase"
	"calibration-system.com/utils/authenticator"
	"github.com/gin-gonic/gin"
)

type RoleController struct {
	router       *gin.Engine
	uc           usecase.RoleUsecase
	tokenService authenticator.AccessToken
	api.BaseApi
}

func (r *RoleController) listHandler(c *gin.Context) {
	roles, err := r.uc.FindAll()
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, roles, "OK")
}

func (r *RoleController) getByHandler(c *gin.Context) {
	name := c.Param("name")
	roles, err := r.uc.FindByName(name)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, roles, "OK")
}

func (r *RoleController) getByIdHandler(c *gin.Context) {
	id := c.Param("id")
	roles, err := r.uc.FindById(id)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, roles, "OK")
}

func (r *RoleController) createHandler(c *gin.Context) {
	var payload model.Role
	if err := r.ParseRequestBody(c, &payload); err != nil {
		r.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := r.uc.SaveData(&payload); err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	r.NewSuccessSingleResponse(c, payload, "OK")
}

func (r *RoleController) updateHandler(c *gin.Context) {
	var payload model.Role

	if err := c.ShouldBind(&payload); err != nil {
		r.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := r.uc.SaveData(&payload); err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	r.NewSuccessSingleResponse(c, payload, "OK")
}

func (r *RoleController) deleteHandler(c *gin.Context) {
	id := c.Param("id")
	if err := r.uc.DeleteData(id); err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusNoContent, "")
}

func NewRoleController(r *gin.Engine, tokenService authenticator.AccessToken, uc usecase.RoleUsecase) *RoleController {
	controller := RoleController{
		router:       r,
		tokenService: tokenService,
		uc:           uc,
	}

	auth := r.Group("/auth").Use(middleware.NewTokenValidator(tokenService).RequireToken())
	r.GET("/roles", controller.listHandler)
	r.POST("/roles-calibration-system", controller.createHandler)
	auth.GET("/roles", controller.listHandler)
	auth.GET("/roles/:name", controller.getByHandler)
	auth.GET("/roles/id/:id", controller.getByIdHandler)
	auth.PUT("/roles", controller.updateHandler)
	auth.POST("/roles", controller.createHandler)
	auth.DELETE("/roles/:id", controller.deleteHandler)
	return &controller
}
