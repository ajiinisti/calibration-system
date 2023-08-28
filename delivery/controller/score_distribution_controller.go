package controller

import (
	"net/http"

	"calibration-system.com/delivery/api"
	"calibration-system.com/model"
	"calibration-system.com/usecase"
	"github.com/gin-gonic/gin"
)

type ScoreDistributionController struct {
	router *gin.Engine
	uc     usecase.ScoreDistributionUsecase
	api.BaseApi
}

func (r *ScoreDistributionController) listHandler(c *gin.Context) {
	scoreDistribution, err := r.uc.FindAll()
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, scoreDistribution, "OK")
}

func (r *ScoreDistributionController) getByIdHandler(c *gin.Context) {
	id := c.Param("id")
	scoreDistribution, err := r.uc.FindById(id)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, scoreDistribution, "OK")
}

func (r *ScoreDistributionController) createHandler(c *gin.Context) {
	var payload model.ScoreDistribution
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

func (r *ScoreDistributionController) updateHandler(c *gin.Context) {
	var payload model.ScoreDistribution

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

func (r *ScoreDistributionController) deleteHandler(c *gin.Context) {
	id := c.Param("id")
	if err := r.uc.DeleteData(id); err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusNoContent, "")
}

func NewScoreDistributionController(r *gin.Engine, uc usecase.ScoreDistributionUsecase) *ScoreDistributionController {
	controller := ScoreDistributionController{
		router: r,
		uc:     uc,
	}
	r.GET("/score-distribution", controller.listHandler)
	r.GET("/score-distribution/:id", controller.getByIdHandler)
	r.PUT("/score-distribution", controller.updateHandler)
	r.POST("/score-distribution", controller.createHandler)
	r.DELETE("/score-distribution/:id", controller.deleteHandler)
	return &controller
}
