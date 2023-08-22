package controller

import (
	"net/http"

	"calibration-system.com/delivery/api"
	"calibration-system.com/model"
	"calibration-system.com/usecase"
	"github.com/gin-gonic/gin"
)

type EmployeeController struct {
	router *gin.Engine
	uc     usecase.EmployeeUsecase
	api.BaseApi
}

func (e *EmployeeController) listHandler(c *gin.Context) {
	employees, err := e.uc.FindAll()
	if err != nil {
		e.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	e.NewSuccessSingleResponse(c, employees, "OK")
}

func (e *EmployeeController) getByHandler(c *gin.Context) {
	name := c.Param("name")
	employees, err := e.uc.FindByEmail(name)
	if err != nil {
		e.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	e.NewSuccessSingleResponse(c, employees, "OK")
}

func (e *EmployeeController) createHandler(c *gin.Context) {
	var payload model.Employee
	if err := e.ParseRequestBody(c, &payload); err != nil {
		e.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := e.uc.SaveData(&payload); err != nil {
		e.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	e.NewSuccessSingleResponse(c, payload, "OK")
}

func (e *EmployeeController) updateHandler(c *gin.Context) {
	var payload model.Employee

	if err := c.ShouldBind(&payload); err != nil {
		e.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := e.uc.SaveData(&payload); err != nil {
		e.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	e.NewSuccessSingleResponse(c, payload, "OK")
}

func (e *EmployeeController) deleteHandler(c *gin.Context) {
	id := c.Param("id")
	if err := e.uc.DeleteData(id); err != nil {
		e.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusNoContent, "")
}

func NewEmployeeController(e *gin.Engine, uc usecase.EmployeeUsecase) *EmployeeController {
	controller := EmployeeController{
		router: e,
		uc:     uc,
	}
	e.GET("/employees", controller.listHandler)
	e.GET("/employees/:email", controller.getByHandler)
	e.PUT("/employees", controller.updateHandler)
	e.POST("/employees", controller.createHandler)
	e.DELETE("/employees/:id", controller.deleteHandler)
	return &controller
}
