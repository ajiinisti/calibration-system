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

type GroupBusinessUnitController struct {
	router       *gin.Engine
	uc           usecase.GroupBusinessUnitUsecase
	tokenService authenticator.AccessToken
	api.BaseApi
}

func (r *GroupBusinessUnitController) listHandler(c *gin.Context) {
	groupBusinessUnits, err := r.uc.FindAll()
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, groupBusinessUnits, "OK")
}

func (r *GroupBusinessUnitController) getByIdHandler(c *gin.Context) {
	id := c.Param("id")
	groupBusinessUnits, err := r.uc.FindById(id)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, groupBusinessUnits, "OK")
}

func (r *GroupBusinessUnitController) createHandler(c *gin.Context) {
	var payload model.GroupBusinessUnit
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

func (r *GroupBusinessUnitController) updateHandler(c *gin.Context) {
	var payload model.GroupBusinessUnit

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

func (r *GroupBusinessUnitController) deleteHandler(c *gin.Context) {
	id := c.Param("id")
	if err := r.uc.DeleteData(id); err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusNoContent, "")
}

func NewGroupBusinessUnitController(r *gin.Engine, tokenService authenticator.AccessToken, uc usecase.GroupBusinessUnitUsecase) *GroupBusinessUnitController {
	controller := GroupBusinessUnitController{
		router:       r,
		tokenService: tokenService,
		uc:           uc,
	}
	auth := r.Group("/auth").Use(middleware.NewTokenValidator(tokenService).RequireToken())
	auth.GET("/group-business-units", controller.listHandler)
	auth.GET("/group-business-units/:id", controller.getByIdHandler)
	auth.PUT("/group-business-units", controller.updateHandler)
	auth.POST("/group-business-units", controller.createHandler)
	auth.DELETE("/group-business-units/:id", controller.deleteHandler)
	return &controller
}
