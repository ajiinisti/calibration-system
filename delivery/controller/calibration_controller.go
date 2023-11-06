package controller

import (
	"net/http"
	"strconv"
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
	projectID := c.Param("projectID")
	projectPhaseID := c.Param("projectPhaseID")
	employeeID := c.Param("employeeID")
	if err := r.uc.DeleteData(projectID, projectPhaseID, employeeID); err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusNoContent, "")
}

func (r *CalibrationController) uploadHandler(c *gin.Context) {
	projectID := c.Request.FormValue("projectID")
	file, err := c.FormFile("excelFile")
	if err != nil {
		r.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = r.uc.BulkInsert(file, projectID)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	r.NewSuccessSingleResponse(c, "", "OK")
}

func (r *CalibrationController) uploadNikHandler(c *gin.Context) {
	projectID := c.Request.FormValue("projectID")
	file, err := c.FormFile("excelFile")
	if err != nil {
		r.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logs, err := r.uc.CheckEmployee(file, projectID)
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
	projectID := c.Request.FormValue("projectID")
	file, err := c.FormFile("excelFile")
	if err != nil {
		r.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logs, err := r.uc.CheckCalibrator(file, projectID)
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

func (r *CalibrationController) getSummaryCalibrationsBySPMOIDHandler(c *gin.Context) {
	spmoID := c.Param("spmoID")
	payload, err := r.uc.FindSummaryCalibrationBySPMOID(spmoID)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	r.NewSuccessSingleResponse(c, payload, "OK")
}

func (r *CalibrationController) getAllActiveCalibrationsBySPMOIDHandler(c *gin.Context) {
	spmoID := c.Param("spmoID")

	payload, err := r.uc.FindActiveBySPMOID(spmoID)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	r.NewSuccessSingleResponse(c, payload, "OK")
}

func (r *CalibrationController) getAllAcceptedCalibrationsBySPMOIDHandler(c *gin.Context) {
	spmoID := c.Param("spmoID")

	payload, err := r.uc.FindAcceptedBySPMOID(spmoID)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	r.NewSuccessSingleResponse(c, payload, "OK")
}

func (r *CalibrationController) getAllRejectdCalibrationsBySPMOIDHandler(c *gin.Context) {
	spmoID := c.Param("spmoID")

	payload, err := r.uc.FindRejectedBySPMOID(spmoID)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	r.NewSuccessSingleResponse(c, payload, "OK")
}

func (r *CalibrationController) spmoAcceptApprovalHandler(c *gin.Context) {
	var payload request.AcceptJustification
	if err := r.ParseRequestBody(c, &payload); err != nil {
		r.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := r.uc.SpmoAcceptApproval(&payload); err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, "", "OK")
}

func (r *CalibrationController) spmoAcceptMultipleApprovalHandler(c *gin.Context) {
	var payload request.AcceptMultipleJustification
	if err := r.ParseRequestBody(c, &payload); err != nil {
		r.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := r.uc.SpmoAcceptMultipleApproval(&payload); err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, "", "OK")
}

func (r *CalibrationController) spmoRejectApprovalHandler(c *gin.Context) {
	var payload request.RejectJustification
	if err := r.ParseRequestBody(c, &payload); err != nil {
		r.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := r.uc.SpmoRejectApproval(&payload); err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, "", "OK")
}

func (r *CalibrationController) getAllDetailActiveCalibrationsBySPMOIDHandler(c *gin.Context) {
	spmoID := c.Param("spmoID")
	calibratorID := c.Param("calibratorID")
	businessUnitID := c.Param("businessUnitID")
	order := c.Param("order")
	department := c.Param("department")

	intOrder, err := strconv.Atoi(order)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
	}

	payload, err := r.uc.FindAllDetailCalibrationbySPMOID(spmoID, calibratorID, businessUnitID, department, intOrder)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	r.NewSuccessSingleResponse(c, payload, "OK")
}

func (r *CalibrationController) getAllDetailActiveCalibrations2BySPMOIDHandler(c *gin.Context) {
	spmoID := c.Param("spmoID")
	calibratorID := c.Param("calibratorID")
	businessUnitID := c.Param("businessUnitID")
	order := c.Param("order")

	intOrder, err := strconv.Atoi(order)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
	}

	payload, err := r.uc.FindAllDetailCalibration2bySPMOID(spmoID, calibratorID, businessUnitID, intOrder)
	if err != nil {
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
	r.GET("/summary-calibrations/spmo/:spmoID", controller.getSummaryCalibrationsBySPMOIDHandler)
	r.GET("/calibrations/spmo/:spmoID", controller.getAllActiveCalibrationsBySPMOIDHandler)
	// r.GET("/calibrations/spmo/:spmoID/:calibratorID/:businessUnitID/:order/:department", controller.getAllDetailActiveCalibrationsBySPMOIDHandler)
	r.GET("/calibrations/spmo/:spmoID/:calibratorID/:businessUnitID/:order", controller.getAllDetailActiveCalibrations2BySPMOIDHandler)
	// r.GET("/calibrations/spmo-accepted/:spmoID", controller.getAllAcceptedCalibrationsBySPMOIDHandler)
	// r.GET("/calibrations/spmo-rejected/:spmoID", controller.getAllRejectdCalibrationsBySPMOIDHandler)
	r.PUT("/calibrations", controller.updateHandler)
	r.POST("/calibrations", controller.createHandler)
	r.POST("/calibrations/upload", controller.uploadHandler)
	r.POST("/calibrations/upload-employee", controller.uploadNikHandler)
	r.POST("/calibrations/upload-calibrator", controller.uploadCalibratorHandler)
	r.POST("/calibrations/save-calibrations", controller.saveCalibrationsHandler)
	r.POST("/calibrations/submit-calibrations/:calibratorID", controller.submitCalibrationsHandler)
	r.POST("/calibrations/accept-approval", controller.spmoAcceptApprovalHandler)
	r.POST("/calibrations/accept-multiple-approval", controller.spmoAcceptMultipleApprovalHandler)
	r.POST("/calibrations/reject-approval", controller.spmoRejectApprovalHandler)
	r.DELETE("/calibrations/:projectID/:projectPhaseID/:employeeID", controller.deleteHandler)
	return &controller
}
