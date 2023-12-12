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

type BottomRemarkController struct {
	router       *gin.Engine
	uc           usecase.BottomRemarkUsecase
	tokenService authenticator.AccessToken
	api.BaseApi
}

func (r *BottomRemarkController) getByIdHandler(c *gin.Context) {
	projectID := c.Param("projectID")
	employeeID := c.Param("employeeID")
	projectPhaseID := c.Param("projectPhaseID")
	BottomRemarks, err := r.uc.FindByForeignKeyID(projectID, employeeID, projectPhaseID)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, BottomRemarks, "OK")
}

func (r *BottomRemarkController) createHandler(c *gin.Context) {
	var payload model.BottomRemark
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

func (r *BottomRemarkController) updateHandler(c *gin.Context) {
	var payload model.BottomRemark

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

func (r *BottomRemarkController) deleteHandler(c *gin.Context) {
	projectID := c.Param("projectID")
	employeeID := c.Param("employeeID")
	projectPhaseID := c.Param("projectPhaseID")
	if err := r.uc.DeleteData(projectID, employeeID, projectPhaseID); err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusNoContent, "")
}

func NewBottomRemarkController(r *gin.Engine, tokenService authenticator.AccessToken, uc usecase.BottomRemarkUsecase) *BottomRemarkController {
	controller := BottomRemarkController{
		router:       r,
		tokenService: tokenService,
		uc:           uc,
	}
	auth := r.Group("/auth").Use(middleware.NewTokenValidator(tokenService).RequireToken())
	auth.GET("/bottom-remark/:projectID/:employeeID/:projectPhaseID", controller.getByIdHandler)
	auth.PUT("/bottom-remark", controller.updateHandler)
	auth.POST("/bottom-remark", controller.createHandler)
	auth.DELETE("/bottom-remark/:projectID/:employeeID/:projectPhaseID", controller.deleteHandler)
	return &controller
}
