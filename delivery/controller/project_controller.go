package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"calibration-system.com/delivery/api"
	"calibration-system.com/delivery/api/request"
	"calibration-system.com/model"
	"calibration-system.com/usecase"
	"github.com/gin-gonic/gin"
)

type ProjectController struct {
	router *gin.Engine
	uc     usecase.ProjectUsecase
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

	param := request.PaginationParam{
		Page:   page,
		Limit:  limit,
		Offset: 0,
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
	fmt.Println("========================DATAAA CONTROLLER========================")
	for _, data := range projects {
		fmt.Println(data.Name)
		fmt.Println(data.CalibrationScores)
	}
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

func NewProjectController(r *gin.Engine, uc usecase.ProjectUsecase) *ProjectController {
	controller := ProjectController{
		router: r,
		uc:     uc,
	}
	r.GET("/projects", controller.listHandler)
	r.GET("/projects/active", controller.getActiveHandler)
	// r.GET("/projects/active/:calibratorID", controller.getActiveHandlerByID)
	r.GET("/projects/active/score-distribution/:businessUnit", controller.getScoreDistributionHandlerByID)
	r.GET("/projects/active/rating-quota", controller.getRatingQuotaHandlerByID)
	r.GET("/projects/active/total-actual-score", controller.getTotalActualScoreHandlerByID)
	r.GET("/projects/summary-calibration/:calibratorID", controller.getSummaryProjectByCalibratorID)
	r.GET("/projects/calibrations/:calibratorID/:prevCalibrator/:businessUnit", controller.getCalibrationsByPrevCalibratorBusinessUnit)
	r.GET("/projects/calibrations-one/:calibratorID/:prevCalibrator/:businessUnit", controller.getNumberOneCalibrationsByPrevCalibratorBusinessUnit)
	r.GET("/projects/calibrations-n-minus-one/:calibratorID/:businessUnit", controller.getNMinusOneCalibrationsByPrevCalibratorBusinessUnit)
	r.GET("/projects/project-phase/:calibratorID", controller.getProjectPhaseByCalibratorId)
	r.GET("/projects/:id", controller.getByIdHandler)
	r.PUT("/projects", controller.updateHandler)
	r.POST("/projects", controller.createHandler)
	r.POST("/projects/publish/:id", controller.publishHandler)
	r.DELETE("/projects/:id", controller.deleteHandler)
	return &controller
}
