package controller

import (
	"net/http"

	"calibration-system.com/delivery/api"
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
	projects, err := r.uc.FindAll()
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, projects, "OK")
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

func NewProjectController(r *gin.Engine, uc usecase.ProjectUsecase) *ProjectController {
	controller := ProjectController{
		router: r,
		uc:     uc,
	}
	r.GET("/projects", controller.listHandler)
	r.GET("/projects/:id", controller.getByIdHandler)
	r.PUT("/projects", controller.updateHandler)
	r.POST("/projects", controller.createHandler)
	r.DELETE("/projects/:id", controller.deleteHandler)
	return &controller
}
