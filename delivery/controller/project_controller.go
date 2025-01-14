package controller

import (
	"fmt"
	"net/http"
	"os"
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

func (r *ProjectController) deactivateHandler(c *gin.Context) {
	id := c.Param("id")
	if err := r.uc.DeactivateProject(id); err != nil {
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

func (r *ProjectController) getActiveHandlerByCalibratorID(c *gin.Context) {
	calibratorID := c.Query("calibratorID")
	projects, err := r.uc.FindActiveProjectByCalibratorID(calibratorID)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, projects, "OK")
}

func (r *ProjectController) getActiveHandlerBySpmoID(c *gin.Context) {
	spmoID := c.Query("spmoID")
	projects, err := r.uc.FindActiveProjectBySpmoID(spmoID)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, projects, "OK")
}

func (r *ProjectController) getScoreDistributionHandlerByID(c *gin.Context) {
	id := c.Query("businessUnit")
	projectID := c.Query("projectID")
	projects, err := r.uc.FindScoreDistributionByCalibratorID(id, projectID)
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
	countCurrentUser := c.Query("countCurrentUser")
	projectID := c.Query("projectID")
	countUser, err := strconv.Atoi(countCurrentUser)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	projects, err := r.uc.FindRatingQuotaByCalibratorID(id, prevCalibrator, businessUnit, types, projectID, countUser)
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
	projectID := c.Query("projectID")
	projects, err := r.uc.FindTotalActualScoreByCalibratorID(id, prevCalibrator, businessUnit, types, projectID)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, projects, "OK")
}

func (r *ProjectController) getTotalCalibratedHandlerByID(c *gin.Context) {
	id := c.Query("calibratorID")
	prevCalibrator := c.Query("prevCalibrator")
	businessUnit := c.Query("businessUnit")
	types := c.Query("type")
	projectID := c.Query("projectID")
	projects, err := r.uc.FindTotalCalibratedByCalibratorID(id, prevCalibrator, businessUnit, types, projectID)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, projects, "OK")
}

func (r *ProjectController) getAverageScoreHandlerByID(c *gin.Context) {
	id := c.Query("calibratorID")
	prevCalibrator := c.Query("prevCalibrator")
	businessUnit := c.Query("businessUnit")
	types := c.Query("type")
	projectID := c.Query("projectID")
	projects, err := r.uc.FindAverageScoreByCalibratorID(id, prevCalibrator, businessUnit, types, projectID)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, projects, "OK")
}

func (r *ProjectController) getEmployeeNameHandler(c *gin.Context) {
	id := c.Query("calibratorID")
	prevCalibrator := c.Query("prevCalibrator")
	businessUnit := c.Query("businessUnit")
	types := c.Query("type")
	projectID := c.Query("projectID")
	projects, err := r.uc.FindAllEmployeeName(id, prevCalibrator, businessUnit, types, projectID)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, projects, "OK")
}

func (r *ProjectController) getSupervisorNameHandler(c *gin.Context) {
	id := c.Query("calibratorID")
	prevCalibrator := c.Query("prevCalibrator")
	businessUnit := c.Query("businessUnit")
	types := c.Query("type")
	projectID := c.Query("projectID")
	projects, err := r.uc.FindAllSupervisorName(id, prevCalibrator, businessUnit, types, projectID)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, projects, "OK")
}

func (r *ProjectController) getGradeHandler(c *gin.Context) {
	id := c.Query("calibratorID")
	prevCalibrator := c.Query("prevCalibrator")
	businessUnit := c.Query("businessUnit")
	types := c.Query("type")
	projectID := c.Query("projectID")
	projects, err := r.uc.FindAllGrade(id, prevCalibrator, businessUnit, types, projectID)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, projects, "OK")
}

func (r *ProjectController) getSummaryProjectByCalibratorID(c *gin.Context) {
	id := c.Query("calibratorID")
	projectID := c.Query("projectID")
	prevCalibratorIDs := c.Query("prevCalibratorIDs")
	var idStrings []string
	if prevCalibratorIDs != "" {
		idStrings = strings.Split(prevCalibratorIDs, ",")
	}

	projects, err := r.uc.FindSummaryProjectByCalibratorID(id, projectID, idStrings)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, projects, "OK")
}

func (r *ProjectController) getCalibrationsByPrevCalibratorBusinessUnit(c *gin.Context) {
	calibratorID := c.Query("calibratorID")
	prevCalibrator := c.Query("prevCalibrator")
	businessUnit := c.Query("businessUnit")
	projectID := c.Query("projectID")

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		r.NewFailedResponse(c, http.StatusBadRequest, "Invalid page number")
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		r.NewFailedResponse(c, http.StatusBadRequest, "Invalid limit number")
	}

	supervisorNames := c.Query("supervisorNames")
	var supervisorList []string
	if supervisorNames != "" {
		supervisorList = strings.Split(supervisorNames, ";")
	}

	employeeNames := c.Query("employeeNames")
	var employeeList []string
	if employeeNames != "" {
		employeeList = strings.Split(employeeNames, ";")
	}

	grades := c.Query("grades")
	var grade []string
	if grades != "" {
		grade = strings.Split(grades, ";")
	}

	orderEmployeeName := c.Query("orderEmployeeName")
	orderGrade := c.Query("orderGrade")
	orderCalibrationScore := c.Query("orderCalibrationScore")
	orderCalibrationRating := c.Query("orderCalibrationRating")
	filterCalibrationRating := c.Query("filterCalibrationRating")
	calibrationPhaseBefore, err := strconv.Atoi(c.Query("calibrationPhaseBefore"))
	if err != nil {
		r.NewFailedResponse(c, http.StatusBadRequest, "Invalid Phase Before")
	}
	orderCalibrationScoreBefore := c.Query("orderCalibrationScoreBefore")
	orderCalibrationRatingBefore := c.Query("orderCalibrationRatingBefore")
	filterCalibrationRatingBefore := c.Query("filterCalibrationRatingBefore")
	ratingChangedStatus := c.Query("ratingChangedStatus")
	ratingChanged, err := strconv.Atoi(c.Query("ratingChanged"))
	if err != nil {
		r.NewFailedResponse(c, http.StatusBadRequest, "Invalid Rating Changed")
	}

	param := request.PaginationParam{
		Page:                          page,
		Limit:                         limit,
		Offset:                        0,
		Name:                          "",
		SupervisorName:                supervisorList,
		Grade:                         grade,
		EmployeeName:                  employeeList,
		OrderGrade:                    orderGrade,
		OrderEmployeeName:             orderEmployeeName,
		OrderCalibrationScore:         orderCalibrationScore,
		OrderCalibrationRating:        orderCalibrationRating,
		FilterCalibrationRating:       filterCalibrationRating,
		CalibrationPhaseBefore:        calibrationPhaseBefore,
		OrderCalibrationScoreBefore:   orderCalibrationScoreBefore,
		OrderCalibrationRatingBefore:  orderCalibrationRatingBefore,
		FilterCalibrationRatingBefore: filterCalibrationRatingBefore,
		RatingChangedStatus:           ratingChangedStatus,
		RatingChanged:                 ratingChanged,
	}

	projects, pagination, err := r.uc.FindCalibrationsByPrevCalibratorBusinessUnitPaginate(calibratorID, prevCalibrator, businessUnit, projectID, param)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	for _, data := range projects.UserData {
		for _, score := range data.CalibrationScores {
			for _, topRemark := range score.TopRemarks {
				topRemark.EvidenceLink = fmt.Sprintf("http://%s/view-initiative/%s", c.Request.Host, topRemark.ID)
			}
		}
	}
	r.NewSuccesPagedResponse(c, []interface{}{projects}, "OK", pagination)
}

func (r *ProjectController) getCalibrationsByBusinessUnit(c *gin.Context) {
	calibratorID := c.Query("calibratorID")
	businessUnit := c.Query("businessUnit")
	projectID := c.Query("projectID")

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		r.NewFailedResponse(c, http.StatusBadRequest, "Invalid page number")
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		r.NewFailedResponse(c, http.StatusBadRequest, "Invalid limit number")
	}

	supervisorNames := c.Query("supervisorNames")
	var supervisorList []string
	if supervisorNames != "" {
		supervisorList = strings.Split(supervisorNames, ";")
	}

	employeeNames := c.Query("employeeNames")
	var employeeList []string
	if employeeNames != "" {
		employeeList = strings.Split(employeeNames, ";")
	}

	grades := c.Query("grades")
	var grade []string
	if grades != "" {
		grade = strings.Split(grades, ";")
	}

	orderEmployeeName := c.Query("orderEmployeeName")
	orderGrade := c.Query("orderGrade")
	orderCalibrationScore := c.Query("orderCalibrationScore")
	orderCalibrationRating := c.Query("orderCalibrationRating")
	filterCalibrationRating := c.Query("filterCalibrationRating")
	calibrationPhaseBefore, err := strconv.Atoi(c.Query("calibrationPhaseBefore"))
	if err != nil {
		r.NewFailedResponse(c, http.StatusBadRequest, "Invalid Phase Before")
	}
	orderCalibrationScoreBefore := c.Query("orderCalibrationScoreBefore")
	orderCalibrationRatingBefore := c.Query("orderCalibrationRatingBefore")
	filterCalibrationRatingBefore := c.Query("filterCalibrationRatingBefore")
	ratingChangedStatus := c.Query("ratingChangedStatus")
	ratingChanged, err := strconv.Atoi(c.Query("ratingChanged"))
	if err != nil {
		r.NewFailedResponse(c, http.StatusBadRequest, "Invalid Rating Changed")
	}

	param := request.PaginationParam{
		Page:                          page,
		Limit:                         limit,
		Offset:                        0,
		Name:                          "",
		SupervisorName:                supervisorList,
		Grade:                         grade,
		EmployeeName:                  employeeList,
		OrderGrade:                    orderGrade,
		OrderEmployeeName:             orderEmployeeName,
		OrderCalibrationScore:         orderCalibrationScore,
		OrderCalibrationRating:        orderCalibrationRating,
		FilterCalibrationRating:       filterCalibrationRating,
		CalibrationPhaseBefore:        calibrationPhaseBefore,
		OrderCalibrationScoreBefore:   orderCalibrationScoreBefore,
		OrderCalibrationRatingBefore:  orderCalibrationRatingBefore,
		FilterCalibrationRatingBefore: filterCalibrationRatingBefore,
		RatingChangedStatus:           ratingChangedStatus,
		RatingChanged:                 ratingChanged,
	}

	projects, pagination, err := r.uc.FindCalibrationsByBusinessUnitPaginate(calibratorID, businessUnit, projectID, param)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	for _, data := range projects.UserData {
		for _, score := range data.CalibrationScores {
			for _, topRemark := range score.TopRemarks {
				topRemark.EvidenceLink = fmt.Sprintf("http://%s/view-initiative/%s", c.Request.Host, topRemark.ID)
			}
		}
	}
	r.NewSuccesPagedResponse(c, []interface{}{projects}, "OK", pagination)
}

func (r *ProjectController) getNumberOneCalibrationsByPrevCalibratorBusinessUnit(c *gin.Context) {
	calibratorID := c.Param("calibratorID")
	prevCalibrator := c.Param("prevCalibrator")
	businessUnit := c.Param("businessUnit")
	projectID := c.Param("projectID")
	projects, err := r.uc.FindNumberOneCalibrationsByPrevCalibratorBusinessUnit(calibratorID, prevCalibrator, businessUnit, projectID)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	for _, data := range projects.UserData {
		for _, score := range data.CalibrationScores {
			for _, topRemark := range score.TopRemarks {
				topRemark.EvidenceLink = fmt.Sprintf("http://%s/view-initiative/%s", c.Request.Host, topRemark.ID)
			}
		}
	}
	r.NewSuccessSingleResponse(c, projects, "OK")
}

func (r *ProjectController) getNMinusOneCalibrationsByPrevCalibratorBusinessUnit(c *gin.Context) {
	calibratorID := c.Query("calibratorID")
	businessUnit := c.Query("businessUnit")
	projectID := c.Query("projectID")

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		r.NewFailedResponse(c, http.StatusBadRequest, "Invalid page number")
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		r.NewFailedResponse(c, http.StatusBadRequest, "Invalid limit number")
	}

	supervisorNames := c.Query("supervisorNames")
	var supervisorList []string
	if supervisorNames != "" {
		supervisorList = strings.Split(supervisorNames, ";")
	}

	employeeNames := c.Query("employeeNames")
	var employeeList []string
	if employeeNames != "" {
		employeeList = strings.Split(employeeNames, ";")
	}

	grades := c.Query("grades")
	var grade []string
	if grades != "" {
		grade = strings.Split(grades, ";")
	}

	orderEmployeeName := c.Query("orderEmployeeName")
	orderGrade := c.Query("orderGrade")
	orderCalibrationScore := c.Query("orderCalibrationScore")
	orderCalibrationRating := c.Query("orderCalibrationRating")
	filterCalibrationRating := c.Query("filterCalibrationRating")
	calibrationPhaseBefore, err := strconv.Atoi(c.Query("calibrationPhaseBefore"))
	if err != nil {
		r.NewFailedResponse(c, http.StatusBadRequest, "Invalid Phase Before")
	}
	orderCalibrationScoreBefore := c.Query("orderCalibrationScoreBefore")
	orderCalibrationRatingBefore := c.Query("orderCalibrationRatingBefore")
	filterCalibrationRatingBefore := c.Query("filterCalibrationRatingBefore")
	ratingChangedStatus := c.Query("ratingChangedStatus")
	ratingChanged, err := strconv.Atoi(c.Query("ratingChanged"))
	if err != nil {
		r.NewFailedResponse(c, http.StatusBadRequest, "Invalid Rating Changed")
	}

	param := request.PaginationParam{
		Page:                          page,
		Limit:                         limit,
		Offset:                        0,
		Name:                          "",
		SupervisorName:                supervisorList,
		Grade:                         grade,
		EmployeeName:                  employeeList,
		OrderGrade:                    orderGrade,
		OrderEmployeeName:             orderEmployeeName,
		OrderCalibrationScore:         orderCalibrationScore,
		OrderCalibrationRating:        orderCalibrationRating,
		FilterCalibrationRating:       filterCalibrationRating,
		CalibrationPhaseBefore:        calibrationPhaseBefore,
		OrderCalibrationScoreBefore:   orderCalibrationScoreBefore,
		OrderCalibrationRatingBefore:  orderCalibrationRatingBefore,
		FilterCalibrationRatingBefore: filterCalibrationRatingBefore,
		RatingChangedStatus:           ratingChangedStatus,
		RatingChanged:                 ratingChanged,
	}

	projects, pagination, err := r.uc.FindNMinusOneCalibrationsByPrevCalibratorBusinessUnitPaginate(calibratorID, businessUnit, projectID, param)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	for _, data := range projects.UserData {
		for _, score := range data.CalibrationScores {
			for _, topRemark := range score.TopRemarks {
				topRemark.EvidenceLink = fmt.Sprintf("http://%s/view-initiative/%s", c.Request.Host, topRemark.ID)
			}
		}
	}
	r.NewSuccesPagedResponse(c, []interface{}{projects}, "OK", pagination)
}

func (r *ProjectController) getCalibrationsByPrevCalibratorBusinessUnitAndRating(c *gin.Context) {
	calibratorID := c.Query("calibratorID")
	prevCalibrator := c.Query("prevCalibrator")
	businessUnit := c.Query("businessUnit")
	rating := c.Query("rating")
	projectID := c.Query("projectID")

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
		Name:   "",
	}

	projects, pagination, err := r.uc.FindCalibrationsByPrevCalibratorBusinessUnitAndRating(calibratorID, prevCalibrator, businessUnit, rating, projectID, param)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	for _, data := range projects.UserData {
		for _, score := range data.CalibrationScores {
			for _, topRemark := range score.TopRemarks {
				topRemark.EvidenceLink = fmt.Sprintf("http://%s/view-initiative/%s", c.Request.Host, topRemark.ID)
			}
		}
	}
	r.NewSuccesPagedResponse(c, []interface{}{projects}, "OK", pagination)
}

func (r *ProjectController) getCalibrationsByBusinessUnitAndRating(c *gin.Context) {
	calibratorID := c.Query("calibratorID")
	businessUnit := c.Query("businessUnit")
	projectID := c.Query("projectID")
	rating := c.Query("rating")

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
		Name:   "",
	}

	projects, pagination, err := r.uc.FindCalibrationsByBusinessUnitAndRating(calibratorID, businessUnit, rating, projectID, param)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	for _, data := range projects.UserData {
		for _, score := range data.CalibrationScores {
			for _, topRemark := range score.TopRemarks {
				topRemark.EvidenceLink = fmt.Sprintf("http://%s/view-initiative/%s", c.Request.Host, topRemark.ID)
			}
		}
	}

	r.NewSuccesPagedResponse(c, []interface{}{projects}, "OK", pagination)
}

func (r *ProjectController) getCalibrationsByRating(c *gin.Context) {
	calibratorID := c.Query("calibratorID")
	projectID := c.Query("projectID")
	rating := c.Query("rating")

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
		Name:   "",
	}

	projects, pagination, err := r.uc.FindCalibrationsByRating(calibratorID, rating, projectID, param)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	for _, data := range projects.UserData {
		for _, score := range data.CalibrationScores {
			for _, topRemark := range score.TopRemarks {
				topRemark.EvidenceLink = fmt.Sprintf("http://%s/view-initiative/%s", c.Request.Host, topRemark.ID)
			}
		}
	}

	r.NewSuccesPagedResponse(c, []interface{}{projects}, "OK", pagination)
}

func (r *ProjectController) getSummaryTotalProjectByCalibrator(c *gin.Context) {
	calibratorID := c.Query("calibratorID")
	projectID := c.Query("projectID")
	projects, err := r.uc.FindSummaryProjectTotalByCalibratorID(calibratorID, projectID)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, projects, "OK")
}

func (r *ProjectController) getProjectPhaseByCalibratorId(c *gin.Context) {
	calibratorID := c.Query("calibratorID")
	projectID := c.Query("projectID")
	projects, err := r.uc.FindCalibratorPhase(calibratorID, projectID)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, projects, "OK")
}

func (r *ProjectController) getProjectPhaseHandler(c *gin.Context) {
	projectID := c.Query("projectID")
	projects, err := r.uc.FindActiveProjectPhase(projectID)
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

func (r *ProjectController) getReportCalibrations(c *gin.Context) {
	types := c.Param("type")
	calibratorID := c.Param("calibratorID")
	businessUnit := c.Param("businessUnit")
	prevCalibrator := c.Param("prevCalibrator")
	projectID := c.Param("projectID")
	file, err := r.uc.ReportCalibrations(types, calibratorID, businessUnit, prevCalibrator, projectID, c)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	defer func() {
		// Clean up: Remove the file after it has been served
		err := os.Remove(file)
		if err != nil {
			fmt.Println("Error removing file:", err)
		}
	}()

	// Set the response headers for downloading
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename="+file)

	// Serve the file
	c.File(file)
}

func (r *ProjectController) getSummaryReportCalibrations(c *gin.Context) {
	calibratorID := c.Query("calibratorID")
	projectID := c.Query("projectID")
	file, err := r.uc.SummaryReportCalibrations(calibratorID, projectID, c)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	defer func() {
		// Clean up: Remove the file after it has been served
		err := os.Remove(file)
		if err != nil {
			fmt.Println("Error removing file:", err)
		}
	}()

	// Set the response headers for downloading
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename="+file)

	// Serve the file
	c.File(file)
}

func (r *ProjectController) getReportAllCalibrations(c *gin.Context) {
	types := c.Query("type")
	calibratorID := c.Query("calibratorID")
	businessUnit := c.Query("businessUnit")
	prevCalibrator := c.Query("prevCalibrator")
	projectID := c.Query("projectID")
	file, err := r.uc.ReportCalibrations(types, calibratorID, businessUnit, prevCalibrator, projectID, c)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	defer func() {
		// Clean up: Remove the file after it has been served
		err := os.Remove(file)
		if err != nil {
			fmt.Println("Error removing file:", err)
		}
	}()

	// Set the response headers for downloading
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename="+file)

	// Serve the file
	c.File(file)
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
	auth.GET("/projects/calibrator", controller.getActiveHandlerByCalibratorID)
	auth.GET("/projects/spmo", controller.getActiveHandlerBySpmoID)
	auth.GET("/projects/score-distribution", controller.getScoreDistributionHandlerByID) // BELUM
	auth.GET("/projects/rating-quota", controller.getRatingQuotaHandlerByID)
	auth.GET("/projects/total-actual-score", controller.getTotalActualScoreHandlerByID)
	auth.GET("/projects/total-calibrated", controller.getTotalCalibratedHandlerByID)
	auth.GET("/projects/average-score", controller.getAverageScoreHandlerByID)
	auth.GET("/projects/employee-name-filter", controller.getEmployeeNameHandler)
	auth.GET("/projects/supervisor-name-filter", controller.getSupervisorNameHandler)
	auth.GET("/projects/grade-filter", controller.getGradeHandler)
	auth.GET("/projects/summary-calibration", controller.getSummaryProjectByCalibratorID)
	auth.GET("/projects/calibrations", controller.getCalibrationsByPrevCalibratorBusinessUnit)
	auth.GET("/projects/calibrations-all-bu", controller.getCalibrationsByBusinessUnit)
	auth.GET("/projects/calibrations-one", controller.getNumberOneCalibrationsByPrevCalibratorBusinessUnit)
	auth.GET("/projects/calibrations-n-minus-one", controller.getNMinusOneCalibrationsByPrevCalibratorBusinessUnit)
	auth.GET("/projects/calibrations-score", controller.getCalibrationsByPrevCalibratorBusinessUnitAndRating)
	auth.GET("/projects/calibrations-score-all-bu", controller.getCalibrationsByBusinessUnitAndRating)
	auth.GET("/projects/calibrations-score-all", controller.getCalibrationsByRating)
	auth.GET("/projects/summary-calibration-total", controller.getSummaryTotalProjectByCalibrator)
	auth.GET("/projects/project-phase/calibrator", controller.getProjectPhaseByCalibratorId)
	auth.GET("/projects/project-phase/manager", controller.getActiveManagerPhaseHandler) //BELUM
	auth.GET("/projects/project-phase", controller.getProjectPhaseHandler)
	auth.GET("/projects/:id", controller.getByIdHandler)
	auth.PUT("/projects", controller.updateHandler)
	auth.POST("/projects", controller.createHandler)
	auth.POST("/projects/publish/:id", controller.publishHandler)
	auth.POST("/projects/deactive/:id", controller.deactivateHandler)
	auth.DELETE("/projects/:id", controller.deleteHandler)
	auth.GET("/projects-report/:type/:calibratorID/:businessUnit/:prevCalibrator/:projectID", controller.getReportCalibrations)
	auth.GET("/projects-report-summary", controller.getSummaryReportCalibrations)
	return &controller
}
