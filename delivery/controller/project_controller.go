package controller

import (
	"net/http"
	"strconv"

	"calibration-system.com/delivery/api"
	"calibration-system.com/delivery/api/request"
	"calibration-system.com/delivery/middleware"
	"calibration-system.com/model"
	"calibration-system.com/usecase"
	"calibration-system.com/utils/authenticator"
	"github.com/gin-gonic/gin"
)

type ProjectController struct {
	router       *gin.Engine
	uc           usecase.ProjectUsecase
	tokenService authenticator.AccessToken
	api.BaseApi
}

func (r *ProjectController) listHandler(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		r.NewFailedResponse(c, http.StatusBadRequest, "Invalid page number")
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		r.NewFailedResponse(c, http.StatusBadRequest, "Invalid limit number")
	}

	nameQuery := c.Query("name")
	param := request.PaginationParam{
		Page:   page,
		Limit:  limit,
		Offset: 0,
		Name:   nameQuery,
	}

	projects, pagination, err := r.uc.FindPagination(param)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var newProjects []interface{}
	for _, v := range projects {
		newProjects = append(newProjects, v)
	}

	r.NewSuccesPagedResponse(c, newProjects, "OK", pagination)
}

func (r *ProjectController) getByIdHandler(c *gin.Context) {
	id := c.Param("id")
	projects, err := r.uc.FindById(id)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, projects, "OK")
}

func (r *ProjectController) createHandler(c *gin.Context) {
	var payload model.Project
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

func (r *ProjectController) publishHandler(c *gin.Context) {
	id := c.Param("id")
	if err := r.uc.PublishProject(id); err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	r.NewSuccessSingleResponse(c, "Success Publish", "OK")
}

func (r *ProjectController) updateHandler(c *gin.Context) {
	var payload model.Project

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

func (r *ProjectController) deleteHandler(c *gin.Context) {
	id := c.Param("id")
	if err := r.uc.DeleteData(id); err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusNoContent, "")
}

// func (r *ProjectController) getActiveHandler(c *gin.Context) {
// 	projects, err := r.uc.FindActiveProject()
// 	if err != nil {
// 		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
// 		return
// 	}
// 	r.NewSuccessSingleResponse(c, projects, "OK")
// }

// func (r *ProjectController) getActiveHandlerByID(c *gin.Context) {
// 	id := c.Param("calibratorID")
// 	projects, err := r.uc.FindActiveProjectByCalibratorID(id)
// 	if err != nil {
// 		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
// 		return
// 	}
// 	r.NewSuccessSingleResponse(c, projects, "OK")
// }

func (r *ProjectController) getActiveHandler(c *gin.Context) {
	projects, err := r.uc.FindActiveProject()
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, projects, "OK")
}

func (r *ProjectController) getScoreDistributionHandlerByID(c *gin.Context) {
	id := c.Param("businessUnit")
	projects, err := r.uc.FindScoreDistributionByCalibratorID(id)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, projects, "OK")
}

func (r *ProjectController) getRatingQuotaHandlerByID(c *gin.Context) {
	id := c.Query("calibratorID")
	prevCalibrator := c.Query("prevCalibrator")
	businessUnit := c.Query("businessUnit")
	types := c.Query("type")
	projects, err := r.uc.FindRatingQuotaByCalibratorID(id, prevCalibrator, businessUnit, types)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, projects, "OK")
}

func (r *ProjectController) getTotalActualScoreHandlerByID(c *gin.Context) {
	id := c.Query("calibratorID")
	prevCalibrator := c.Query("prevCalibrator")
	businessUnit := c.Query("businessUnit")
	types := c.Query("type")
	projects, err := r.uc.FindTotalActualScoreByCalibratorID(id, prevCalibrator, businessUnit, types)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, projects, "OK")
}

func (r *ProjectController) getSummaryProjectByCalibratorID(c *gin.Context) {
	id := c.Param("calibratorID")
	projects, err := r.uc.FindSummaryProjectByCalibratorID(id)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, projects, "OK")
}

func (r *ProjectController) getCalibrationsByPrevCalibratorBusinessUnit(c *gin.Context) {
	calibratorID := c.Param("calibratorID")
	prevCalibrator := c.Param("prevCalibrator")
	businessUnit := c.Param("businessUnit")
	projects, err := r.uc.FindCalibrationsByPrevCalibratorBusinessUnit(calibratorID, prevCalibrator, businessUnit)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, projects, "OK")
}

func (r *ProjectController) getNumberOneCalibrationsByPrevCalibratorBusinessUnit(c *gin.Context) {
	calibratorID := c.Param("calibratorID")
	prevCalibrator := c.Param("prevCalibrator")
	businessUnit := c.Param("businessUnit")
	projects, err := r.uc.FindNumberOneCalibrationsByPrevCalibratorBusinessUnit(calibratorID, prevCalibrator, businessUnit)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, projects, "OK")
}

func (r *ProjectController) getNMinusOneCalibrationsByPrevCalibratorBusinessUnit(c *gin.Context) {
	calibratorID := c.Param("calibratorID")
	businessUnit := c.Param("businessUnit")
	projects, err := r.uc.FindNMinusOneCalibrationsByPrevCalibratorBusinessUnit(calibratorID, businessUnit)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, projects, "OK")
}

func (r *ProjectController) getProjectPhaseByCalibratorId(c *gin.Context) {
	calibratorID := c.Param("calibratorID")
	projects, err := r.uc.FindCalibratorPhase(calibratorID)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, projects, "OK")
}

func (r *ProjectController) getActiveProjectPhaseHandler(c *gin.Context) {
	projects, err := r.uc.FindActiveProjectPhase()
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, projects, "OK")
}

func (r *ProjectController) getActiveManagerPhaseHandler(c *gin.Context) {
	projects, err := r.uc.FindActiveManagerPhase()
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, projects, "OK")
}

func NewProjectController(r *gin.Engine, tokenService authenticator.AccessToken, uc usecase.ProjectUsecase) *ProjectController {
	controller := ProjectController{
		router:       r,
		tokenService: tokenService,
		uc:           uc,
	}
	auth := r.Group("/auth").Use(middleware.NewTokenValidator(tokenService).RequireToken())
	auth.GET("/projects", controller.listHandler)
	auth.GET("/projects/active", controller.getActiveHandler)
	// auth.GET("/projects/active/:calibratorID", controller.getActiveHandlerByID)
	auth.GET("/projects/active/score-distribution/:businessUnit", controller.getScoreDistributionHandlerByID)
	auth.GET("/projects/active/rating-quota", controller.getRatingQuotaHandlerByID)
	auth.GET("/projects/active/total-actual-score", controller.getTotalActualScoreHandlerByID)
	auth.GET("/projects/summary-calibration/:calibratorID", controller.getSummaryProjectByCalibratorID)
	auth.GET("/projects/calibrations/:calibratorID/:prevCalibrator/:businessUnit", controller.getCalibrationsByPrevCalibratorBusinessUnit)
	auth.GET("/projects/calibrations-one/:calibratorID/:prevCalibrator/:businessUnit", controller.getNumberOneCalibrationsByPrevCalibratorBusinessUnit)
	auth.GET("/projects/calibrations-n-minus-one/:calibratorID/:businessUnit", controller.getNMinusOneCalibrationsByPrevCalibratorBusinessUnit)
	auth.GET("/projects/project-phase/:calibratorID", controller.getProjectPhaseByCalibratorId)
	auth.GET("/projects/project-phase/manager", controller.getActiveManagerPhaseHandler)
	auth.GET("/projects/project-phase/active", controller.getActiveProjectPhaseHandler)
	auth.GET("/projects/:id", controller.getByIdHandler)
	auth.PUT("/projects", controller.updateHandler)
	auth.POST("/projects", controller.createHandler)
	auth.POST("/projects/publish/:id", controller.publishHandler)
	auth.DELETE("/projects/:id", controller.deleteHandler)
	return &controller
}
