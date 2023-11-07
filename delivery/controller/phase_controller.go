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

type PhaseController struct {
	router       *gin.Engine
	uc           usecase.PhaseUsecase
	tokenService authenticator.AccessToken
	api.BaseApi
}

func (r *PhaseController) listHandler(c *gin.Context) {
	phases, err := r.uc.FindAll()
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, phases, "OK")
}

func (r *PhaseController) getByIdHandler(c *gin.Context) {
	id := c.Param("id")
	phases, err := r.uc.FindById(id)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, phases, "OK")
}

func (r *PhaseController) createHandler(c *gin.Context) {
	var payload model.Phase
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

func (r *PhaseController) updateHandler(c *gin.Context) {
	var payload model.Phase

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

func (r *PhaseController) deleteHandler(c *gin.Context) {
	id := c.Param("id")
	if err := r.uc.DeleteData(id); err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusNoContent, "")
}

func NewPhaseController(r *gin.Engine, tokenService authenticator.AccessToken, uc usecase.PhaseUsecase) *PhaseController {
	controller := PhaseController{
		router:       r,
		tokenService: tokenService,
		uc:           uc,
	}
	auth := r.Use(middleware.NewTokenValidator(tokenService).RequireToken())
	auth.GET("/phases", controller.listHandler)
	auth.GET("/phases/:id", controller.getByIdHandler)
	auth.PUT("/phases", controller.updateHandler)
	auth.POST("/phases", controller.createHandler)
	auth.DELETE("/phases/:id", controller.deleteHandler)
	return &controller
}
