package controller

import (
	"net/http"

	"calibration-system.com/delivery/api"
	"calibration-system.com/delivery/api/request"
	"calibration-system.com/delivery/middleware"
	"calibration-system.com/model"
	"calibration-system.com/usecase"
	"calibration-system.com/utils/authenticator"
	"github.com/gin-gonic/gin"
)

type TopRemarkController struct {
	router       *gin.Engine
	uc           usecase.TopRemarkUsecase
	tokenService authenticator.AccessToken
	api.BaseApi
}

func (r *TopRemarkController) getByIdHandler(c *gin.Context) {
	projectID := c.Param("projectID")
	employeeID := c.Param("employeeID")
	projectPhaseID := c.Param("projectPhaseID")
	TopRemarks, err := r.uc.FindByForeignKeyID(projectID, employeeID, projectPhaseID)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, TopRemarks, "OK")
}

func (r *TopRemarkController) createHandler(c *gin.Context) {
	var payload model.TopRemark
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

func (r *TopRemarkController) createHandlerByProject(c *gin.Context) {
	var payload []*model.TopRemark

	if err := c.ShouldBind(&payload); err != nil {
		r.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := r.uc.SaveDataByProject(payload); err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	r.NewSuccessSingleResponse(c, payload, "OK")
}

func (r *TopRemarkController) updateHandler(c *gin.Context) {
	var payload model.TopRemark

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

func (r *TopRemarkController) deleteHandler(c *gin.Context) {
	projectID := c.Param("projectID")
	employeeID := c.Param("employeeID")
	projectPhaseID := c.Param("projectPhaseID")
	if err := r.uc.DeleteData(projectID, employeeID, projectPhaseID); err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusNoContent, "")
}

func (r *TopRemarkController) deleteHandlerByProject(c *gin.Context) {
	var payload request.DeleteTopRemarks
	if err := c.ShouldBind(&payload); err != nil {
		r.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := r.uc.BulkDeleteData(payload); err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusNoContent, "")
}

func NewTopRemarkController(r *gin.Engine, tokenService authenticator.AccessToken, uc usecase.TopRemarkUsecase) *TopRemarkController {
	controller := TopRemarkController{
		router:       r,
		tokenService: tokenService,
		uc:           uc,
	}

	auth := r.Use(middleware.NewTokenValidator(tokenService).RequireToken())
	auth.GET("/top-remark/:projectID/:employeeID/:projectPhaseID", controller.getByIdHandler)
	auth.PUT("/top-remark", controller.updateHandler)
	auth.POST("/top-remark", controller.createHandler)
	auth.POST("/top-remark/project", controller.createHandlerByProject)
	auth.POST("/top-remark/delete", controller.deleteHandlerByProject)
	auth.DELETE("/top-remark/:id", controller.deleteHandler)
	return &controller
}
