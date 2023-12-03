package controller

import (
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

type RatingQuotaController struct {
	router       *gin.Engine
	uc           usecase.RatingQuotaUsecase
	tokenService authenticator.AccessToken
	api.BaseApi
}

func (r *RatingQuotaController) listHandler(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		r.NewFailedResponse(c, http.StatusBadRequest, "Invalid page number")
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		r.NewFailedResponse(c, http.StatusBadRequest, "Invalid limit number")
	}

	projectId := c.Query("id")
	param := request.PaginationParam{
		Page:   page,
		Limit:  limit,
		Offset: 0,
	}

	ratingQuotas, pagination, err := r.uc.FindPagination(param, projectId)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var newRatings []interface{}
	for _, v := range ratingQuotas {
		newRatings = append(newRatings, v)
	}

	r.NewSuccesPagedResponse(c, newRatings, "OK", pagination)
}

func (r *RatingQuotaController) getByProjectHandler(c *gin.Context) {
	projectId := c.Param("projectId")
	ratingQuotas, err := r.uc.FindByProject(projectId)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, ratingQuotas, "OK")
}

func (r *RatingQuotaController) getByIDHandler(c *gin.Context) {
	projectId := c.Param("projectId")
	businessUnitId := c.Param("businessUnitId")
	ratingQuotas, err := r.uc.FindById(projectId, businessUnitId)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, ratingQuotas, "OK")
}

func (r *RatingQuotaController) createHandler(c *gin.Context) {
	var payload model.RatingQuota
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

func (r *RatingQuotaController) updateHandler(c *gin.Context) {
	var payload model.RatingQuota

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

func (r *RatingQuotaController) deleteHandler(c *gin.Context) {
	projectId := c.Param("projectId")
	businessUnitId := c.Param("businessUnitId")
	if err := r.uc.DeleteData(projectId, businessUnitId); err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusNoContent, "")
}

func (r *RatingQuotaController) uploadHandler(c *gin.Context) {
	projectId := c.Request.FormValue("projectID")
	file, err := c.FormFile("excelFile")
	if err != nil {
		r.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logs, err := r.uc.BulkInsert(file, projectId)
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

func NewRatingQuotaController(r *gin.Engine, tokenService authenticator.AccessToken, uc usecase.RatingQuotaUsecase) *RatingQuotaController {
	controller := RatingQuotaController{
		router:       r,
		tokenService: tokenService,
		uc:           uc,
	}
	auth := r.Group("/auth").Use(middleware.NewTokenValidator(tokenService).RequireToken())
	auth.GET("/rating-quotas", controller.listHandler)
	auth.GET("/rating-quotas/:projectId/:businessUnitId", controller.getByIDHandler)
	auth.GET("/rating-quotas/:projectId", controller.getByProjectHandler)
	auth.PUT("/rating-quotas", controller.updateHandler)
	auth.POST("/rating-quotas", controller.createHandler)
	auth.POST("/rating-quotas/upload", controller.uploadHandler)
	auth.DELETE("/rating-quotas/:projectId/:businessUnitId", controller.deleteHandler)
	return &controller
}
