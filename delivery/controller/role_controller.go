package controller

import (
	"net/http"

	"calibration-system.com/delivery/api"
	"calibration-system.com/model"
	"calibration-system.com/usecase"
	"github.com/gin-gonic/gin"
)

type RoleController struct {
	router *gin.Engine
	uc     usecase.RoleUsecase
	api.BaseApi
}

func (r *RoleController) listHandler(c *gin.Context) {
	roles, err := r.uc.FindAll()
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, roles, "OK")
}

func (r *RoleController) getByHandler(c *gin.Context) {
	name := c.Param("name")
	roles, err := r.uc.FindByName(name)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, roles, "OK")
}

func (r *RoleController) createHandler(c *gin.Context) {
	var payload model.Role
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

func (r *RoleController) updateHandler(c *gin.Context) {
	var payload model.Role

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

func (r *RoleController) deleteHandler(c *gin.Context) {
	id := c.Param("id")
	if err := r.uc.DeleteData(id); err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusNoContent, "")
}

func NewRoleController(r *gin.Engine, uc usecase.RoleUsecase) *RoleController {
	controller := RoleController{
		router: r,
		uc:     uc,
	}
	r.GET("/roles", controller.listHandler)
	r.GET("/roles/:name", controller.getByHandler)
	r.PUT("/roles", controller.updateHandler)
	r.POST("/roles", controller.createHandler)
	r.DELETE("/roles/:id", controller.deleteHandler)
	return &controller
}
