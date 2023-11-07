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

type ProjectPhaseController struct {
	router       *gin.Engine
	uc           usecase.ProjectPhaseUsecase
	tokenService authenticator.AccessToken
	api.BaseApi
}

func (r *ProjectPhaseController) listHandler(c *gin.Context) {
	projectPhases, err := r.uc.FindAll()
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, projectPhases, "OK")
}

func (r *ProjectPhaseController) getByIdHandler(c *gin.Context) {
	id := c.Param("id")
	projectPhases, err := r.uc.FindById(id)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, projectPhases, "OK")
}

func (r *ProjectPhaseController) createHandler(c *gin.Context) {
	var payload model.ProjectPhase
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

func (r *ProjectPhaseController) updateHandler(c *gin.Context) {
	var payload model.ProjectPhase

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

func (r *ProjectPhaseController) deleteHandler(c *gin.Context) {
	id := c.Param("id")
	if err := r.uc.DeleteData(id); err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusNoContent, "")
}

func NewProjectPhaseController(r *gin.Engine, tokenService authenticator.AccessToken, uc usecase.ProjectPhaseUsecase) *ProjectPhaseController {
	controller := ProjectPhaseController{
		router:       r,
		tokenService: tokenService,
		uc:           uc,
	}
	auth := r.Use(middleware.NewTokenValidator(tokenService).RequireToken())
	auth.GET("/project-phases", controller.listHandler)
	auth.GET("/project-phases/:id", controller.getByIdHandler)
	auth.PUT("/project-phases", controller.updateHandler)
	auth.POST("/project-phases", controller.createHandler)
	auth.DELETE("/project-phases/:id", controller.deleteHandler)
	return &controller
}
