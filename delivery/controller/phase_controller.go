package controller

import (
	"net/http"

	"calibration-system.com/delivery/api"
	"calibration-system.com/model"
	"calibration-system.com/usecase"
	"github.com/gin-gonic/gin"
)

type PhaseController struct {
	router *gin.Engine
	uc     usecase.PhaseUsecase
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

func NewPhaseController(r *gin.Engine, uc usecase.PhaseUsecase) *PhaseController {
	controller := PhaseController{
		router: r,
		uc:     uc,
	}
	r.GET("/phases", controller.listHandler)
	r.GET("/phases/:id", controller.getByIdHandler)
	r.PUT("/phases", controller.updateHandler)
	r.POST("/phases", controller.createHandler)
	r.DELETE("/phases/:id", controller.deleteHandler)
	return &controller
}
