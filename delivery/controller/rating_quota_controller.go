package controller

import (
	"net/http"

	"calibration-system.com/delivery/api"
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
	ratingQuotas, err := r.uc.FindAll()
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, ratingQuotas, "OK")
}

func (r *RatingQuotaController) getByIdHandler(c *gin.Context) {
	id := c.Param("id")
	ratingQuotas, err := r.uc.FindById(id)
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
	id := c.Param("id")
	if err := r.uc.DeleteData(id); err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusNoContent, "")
}

func NewRatingQuotaController(r *gin.Engine, uc usecase.RatingQuotaUsecase) *RatingQuotaController {
	controller := RatingQuotaController{
		router: r,
		uc:     uc,
	}
	r.GET("/rating-quotas", controller.listHandler)
	r.GET("/rating-quotas/:id", controller.getByIdHandler)
	r.PUT("/rating-quotas", controller.updateHandler)
	r.POST("/rating-quotas", controller.createHandler)
	r.DELETE("/rating-quotas/:id", controller.deleteHandler)
	return &controller
}
