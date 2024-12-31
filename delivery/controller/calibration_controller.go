package controller

import (
	"fmt"
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

type CalibrationController struct {
	router       *gin.Engine
	uc           usecase.CalibrationUsecase
	tokenService authenticator.AccessToken
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
	projectID := c.Param("projectID")
	projectPhaseID := c.Param("projectPhaseID")
	employeeID := c.Param("employeeID")
	groupCalibrations, err := r.uc.FindById(projectID, projectPhaseID, employeeID)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, groupCalibrations, "OK")
}

func (r *CalibrationController) getByProjectEmployeeIdHandler(c *gin.Context) {
	projectID := c.Param("projectID")
	employeeID := c.Param("employeeID")
	groupCalibrations, err := r.uc.FindByProjectEmployeeId(projectID, employeeID)
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
	employeeID := c.Param("employeeID")
	if err := r.uc.DeleteData(projectID, employeeID); err != nil {
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
			r.NewFailedResponse(c, http.StatusInternalServerError, strings.Join(logs, ","))
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
			r.NewFailedResponse(c, http.StatusInternalServerError, strings.Join(logs, ","))
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

func (r *CalibrationController) saveCommentCalibrationHandler(c *gin.Context) {
	var payload model.Calibration
	if err := r.ParseRequestBody(c, &payload); err != nil {
		r.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := r.uc.SaveCommentCalibration(&payload); err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	r.NewSuccessSingleResponse(c, payload, "OK")
}

func (r *CalibrationController) submitCalibrationsHandler(c *gin.Context) {
	calibratorID := c.Param("calibratorID")
	projectID := c.Param("projectID")
	var payload request.CalibrationRequest
	if err := r.ParseRequestBody(c, &payload); err != nil {
		r.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := r.uc.SubmitCalibrations(&payload, calibratorID, projectID); err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	r.NewSuccessSingleResponse(c, payload, "OK")
}

func (r *CalibrationController) sendCalibrationToManagerHandler(c *gin.Context) {
	calibratorID := c.Query("calibratorID")
	projectID := c.Query("projectID")
	var payload request.CalibrationRequest
	if err := r.ParseRequestBody(c, &payload); err != nil {
		r.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := r.uc.SendCalibrationsToManager(&payload, calibratorID, projectID); err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	r.NewSuccessSingleResponse(c, payload, "OK")
}

func (r *CalibrationController) sendBackCalibrationsToOnePhaseBeforeHandler(c *gin.Context) {
	calibratorID := c.Query("calibratorID")
	projectID := c.Query("projectID")
	var payload request.CalibrationRequest
	if err := r.ParseRequestBody(c, &payload); err != nil {
		r.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := r.uc.SendBackCalibrationsToOnePhaseBefore(&payload, calibratorID, projectID); err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	r.NewSuccessSingleResponse(c, payload, "OK")
}

func (r *CalibrationController) getSummaryCalibrationsBySPMOIDHandler(c *gin.Context) {
	spmoID := c.Query("spmoID")
	projectID := c.Query("projectID")
	payload, err := r.uc.FindSummaryCalibrationBySPMOID(spmoID, projectID)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	r.NewSuccessSingleResponse(c, payload, "OK")
}

func (r *CalibrationController) getAllActiveCalibrationsBySPMOIDHandler(c *gin.Context) {
	spmoID := c.Param("spmoID")

	payload, err := r.uc.FindActiveUserBySPMOID(spmoID)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	r.NewSuccessSingleResponse(c, payload, "OK")
}

func (r *CalibrationController) spmoAcceptApprovalHandler(c *gin.Context) {
	var payload request.AcceptJustification
	// projectID := c.Query("projectID")
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

func (r *CalibrationController) spmoSubmitHandler(c *gin.Context) {
	var payload request.AcceptMultipleJustification
	if err := r.ParseRequestBody(c, &payload); err != nil {
		r.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := r.uc.SpmoSubmit(&payload); err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, "", "OK")
}

func (r *CalibrationController) getAllDetailActiveCalibrations2BySPMOIDHandler(c *gin.Context) {
	spmoID := c.Param("spmoID")
	calibratorID := c.Param("calibratorID")
	businessUnitID := c.Param("businessUnitID")
	order := c.Param("order")
	projectID := c.Param("projectID")

	intOrder, err := strconv.Atoi(order)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
	}

	payload, err := r.uc.FindAllDetailCalibration2bySPMOID(spmoID, calibratorID, businessUnitID, projectID, intOrder)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	for _, data := range payload {
		for _, score := range data.CalibrationScores {
			for _, topRemark := range score.TopRemarks {
				topRemark.EvidenceLink = fmt.Sprintf("http://%s/view-initiative/%s", c.Request.Host, topRemark.ID)
			}
		}
	}

	r.NewSuccessSingleResponse(c, payload, "OK")
}

func (r *CalibrationController) createByUserHandler(c *gin.Context) {
	var payload request.CalibrationForm
	if err := r.ParseRequestBody(c, &payload); err != nil {
		r.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := r.uc.SaveDataByUser(&payload); err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	r.NewSuccessSingleResponse(c, payload, "OK")
}

func (r *CalibrationController) sendNotificationFirstCalibratorHandler(c *gin.Context) {
	data, err := r.uc.SendNotificationToCurrentCalibrator()
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, data, "OK")
}

func (r *CalibrationController) getRatingQuotaSPMOHandlerByID(c *gin.Context) {
	spmoID := c.Param("spmoID")
	calibratorID := c.Param("calibratorID")
	businessUnitID := c.Param("businessUnitID")
	order := c.Param("order")
	projectID := c.Param("projectID")

	intOrder, err := strconv.Atoi(order)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
	}

	projects, err := r.uc.FindRatingQuotaSPMOByCalibratorID(spmoID, calibratorID, businessUnitID, projectID, intOrder)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, projects, "OK")
}

func NewCalibrationController(r *gin.Engine, tokenService authenticator.AccessToken, uc usecase.CalibrationUsecase) *CalibrationController {
	controller := CalibrationController{
		router:       r,
		tokenService: tokenService,
		uc:           uc,
	}
	auth := r.Group("/auth").Use(middleware.NewTokenValidator(tokenService).RequireToken())
	auth.GET("/calibrations", controller.listHandler)
	auth.GET("/calibrations/:projectID/:projectPhaseID/:employeeID", controller.getByIdHandler)
	auth.GET("/calibrations-project-employee/:projectID/:employeeID", controller.getByProjectEmployeeIdHandler)
	auth.GET("/summary-calibrations/spmo", controller.getSummaryCalibrationsBySPMOIDHandler)
	auth.GET("/calibrations/spmo/:spmoID", controller.getAllActiveCalibrationsBySPMOIDHandler)
	auth.GET("/calibrations/spmo/:spmoID/:calibratorID/:businessUnitID/:order/:projectID", controller.getAllDetailActiveCalibrations2BySPMOIDHandler)
	auth.GET("/calibrations/spmo/rating-quota/:spmoID/:calibratorID/:businessUnitID/:order/:projectID", controller.getRatingQuotaSPMOHandlerByID)
	auth.PUT("/calibrations", controller.updateHandler)
	auth.POST("/calibrations", controller.createHandler)
	auth.POST("/calibrations-user", controller.createByUserHandler)
	auth.POST("/calibrations/upload", controller.uploadHandler)
	auth.POST("/calibrations/upload-employee", controller.uploadNikHandler)
	auth.POST("/calibrations/upload-calibrator", controller.uploadCalibratorHandler)
	auth.POST("/calibrations/save-comment", controller.saveCommentCalibrationHandler)
	auth.POST("/calibrations/save-calibrations", controller.saveCalibrationsHandler)
	auth.POST("/calibrations/submit-calibrations/:calibratorID/:projectID", controller.submitCalibrationsHandler)
	auth.POST("/calibrations/send-calibration-to-manager", controller.sendCalibrationToManagerHandler)
	auth.POST("/calibrations/send-calibrations-back", controller.sendBackCalibrationsToOnePhaseBeforeHandler)
	auth.POST("/calibrations/accept-approval", controller.spmoAcceptApprovalHandler)
	auth.POST("/calibrations/accept-multiple-approval", controller.spmoAcceptMultipleApprovalHandler)
	auth.POST("/calibrations/reject-approval", controller.spmoRejectApprovalHandler)
	auth.POST("/calibrations/spmo/submit", controller.spmoSubmitHandler)
	auth.DELETE("/calibrations/:projectID/:employeeID", controller.deleteHandler)
	// auth.POST("/projects/send-notificaition-first-calibrator", controller.sendNotificationFirstCalibratorHandler)
	return &controller
}
