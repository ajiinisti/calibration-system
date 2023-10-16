package controller

import (
	"net/http"
	"strings"

	"calibration-system.com/delivery/api"
	"calibration-system.com/delivery/api/request"
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
	projectId := c.Param("projectId")
	projectPhaseId := c.Param("projectPhaseId")
	employeeId := c.Param("employeeId")
	if err := r.uc.DeleteData(projectId, projectPhaseId, employeeId); err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusNoContent, "")
}

func (r *CalibrationController) uploadHandler(c *gin.Context) {
	projectId := c.Request.FormValue("projectId")
	file, err := c.FormFile("excelFile")
	if err != nil {
		r.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = r.uc.BulkInsert(file, projectId)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	r.NewSuccessSingleResponse(c, "", "OK")
}

func (r *CalibrationController) uploadNikHandler(c *gin.Context) {
	projectId := c.Request.FormValue("projectId")
	file, err := c.FormFile("excelFile")
	if err != nil {
		r.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logs, err := r.uc.CheckEmployee(file, projectId)
	if err != nil {
		if len(logs) > 0 {
			r.NewFailedResponse(c, http.StatusInternalServerError, strings.Join(logs, "."))
		} else {
			r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	r.NewSuccessSingleResponse(c, "", "OK")
}

func (r *CalibrationController) uploadCalibratorHandler(c *gin.Context) {
	projectId := c.Request.FormValue("projectId")
	file, err := c.FormFile("excelFile")
	if err != nil {
		r.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logs, err := r.uc.CheckCalibrator(file, projectId)
	if err != nil {
		if len(logs) > 0 {
			r.NewFailedResponse(c, http.StatusInternalServerError, strings.Join(logs, "."))
		} else {
			r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	r.NewSuccessSingleResponse(c, "", "OK")
}

func (r *CalibrationController) saveCalibrationsHandler(c *gin.Context) {
	var payload request.CalibrationRequest
	if err := r.ParseRequestBody(c, &payload); err != nil {
		r.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := r.uc.SaveCalibrations(&payload); err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	r.NewSuccessSingleResponse(c, payload, "OK")
}

func (r *CalibrationController) submitCalibrationsHandler(c *gin.Context) {
	calibratorID := c.Param("calibratorID")
	var payload request.CalibrationRequest
	if err := r.ParseRequestBody(c, &payload); err != nil {
		r.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := r.uc.SubmitCalibrations(&payload, calibratorID); err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	r.NewSuccessSingleResponse(c, payload, "OK")
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
	r.POST("/calibrations/upload", controller.uploadHandler)
	r.POST("/calibrations/upload-employee", controller.uploadNikHandler)
	r.POST("/calibrations/upload-calibrator", controller.uploadCalibratorHandler)
	r.POST("/calibrations/save-calibrations", controller.saveCalibrationsHandler)
	r.POST("/calibrations/submit-calibrations/:calibratorID", controller.submitCalibrationsHandler)
	r.DELETE("/calibrations/:projectId/:projectPhaseId/:employeeId", controller.deleteHandler)
	return &controller
}
