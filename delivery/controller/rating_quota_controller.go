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

type RatingQuotaController struct {
	router *gin.Engine
	uc     usecase.RatingQuotaUsecase
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

func (r *RatingQuotaController) getByIdHandler(c *gin.Context) {
	projectId := c.Param("projectId")
	ratingQuotas, err := r.uc.FindById(projectId)
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
	// Menerima file Excel dari permintaan HTTP POST
	projectId := c.Request.FormValue("projectId")
	file, err := c.FormFile("excelFile")
	if err != nil {
		r.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logs, err := r.uc.BulkInsert(file, projectId)
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

func NewRatingQuotaController(r *gin.Engine, uc usecase.RatingQuotaUsecase) *RatingQuotaController {
	controller := RatingQuotaController{
		router: r,
		uc:     uc,
	}
	r.GET("/rating-quotas", controller.listHandler)
	r.GET("/rating-quotas/:projectId", controller.getByIdHandler)
	r.PUT("/rating-quotas", controller.updateHandler)
	r.POST("/rating-quotas", controller.createHandler)
	r.POST("/rating-quotas/upload", controller.uploadHandler)
	r.DELETE("/rating-quotas/:projectId/:businessUnitId", controller.deleteHandler)
	return &controller
}