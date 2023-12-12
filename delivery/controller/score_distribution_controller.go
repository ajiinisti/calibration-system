package controller

import (
	"net/http"

	"calibration-system.com/delivery/api"
	"calibration-system.com/delivery/middleware"
	"calibration-system.com/model"
	"calibration-system.com/usecase"
	"calibration-system.com/utils/authenticator"
	"github.com/gin-gonic/gin"
)

type ScoreDistributionController struct {
	router       *gin.Engine
	uc           usecase.ScoreDistributionUsecase
	tokenService authenticator.AccessToken
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
	projectId := c.Param("projectId")
	groupBusinessUnitId := c.Param("groupBusinessUnitId")
	if err := r.uc.DeleteData(projectId, groupBusinessUnitId); err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusNoContent, "")
}

func NewScoreDistributionController(r *gin.Engine, tokenService authenticator.AccessToken, uc usecase.ScoreDistributionUsecase) *ScoreDistributionController {
	controller := ScoreDistributionController{
		router:       r,
		tokenService: tokenService,
		uc:           uc,
	}

	auth := r.Group("/auth").Use(middleware.NewTokenValidator(tokenService).RequireToken())
	auth.GET("/score-distribution", controller.listHandler)
	auth.GET("/score-distribution/:id", controller.getByIdHandler)
	auth.PUT("/score-distribution", controller.updateHandler)
	auth.POST("/score-distribution", controller.createHandler)
	auth.DELETE("/score-distribution/:projectId/:groupBusinessUnitId", controller.deleteHandler)
	return &controller
}
