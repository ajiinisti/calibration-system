package controller

import (
	"net/http"

	"calibration-system.com/delivery/api"
	"calibration-system.com/model"
	"calibration-system.com/usecase"
	"github.com/gin-gonic/gin"
)

type CalibrationController struct {
	router *gin.Engine
	uc     usecase.CalibrationUsecase
	api.BaseApi
}

func (r *CalibrationController) listHandler(c *gin.Context) {
	calibrations, err := r.uc.FindAll()
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, calibrations, "OK")
}

func (r *CalibrationController) getByIdHandler(c *gin.Context) {
	id := c.Param("id")
	groupCalibrations, err := r.uc.FindById(id)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, groupCalibrations, "OK")
}

func (r *CalibrationController) createHandler(c *gin.Context) {
	var payload model.Calibration
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

func (r *CalibrationController) updateHandler(c *gin.Context) {
	var payload model.Calibration

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

func (r *CalibrationController) deleteHandler(c *gin.Context) {
	id := c.Param("id")
	if err := r.uc.DeleteData(id); err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusNoContent, "")
}

func NewCalibrationController(r *gin.Engine, uc usecase.CalibrationUsecase) *CalibrationController {
	controller := CalibrationController{
		router: r,
		uc:     uc,
	}
	r.GET("/calibrations", controller.listHandler)
	r.GET("/calibrations/:id", controller.getByIdHandler)
	r.PUT("/calibrations", controller.updateHandler)
	r.POST("/calibrations", controller.createHandler)
	r.DELETE("/calibrations/:id", controller.deleteHandler)
	return &controller
}
