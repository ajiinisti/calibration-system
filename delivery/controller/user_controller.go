package controller

import (
	"net/http"

	"calibration-system.com/delivery/api"
	"calibration-system.com/delivery/api/request"
	"calibration-system.com/model"
	"calibration-system.com/usecase"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	router *gin.Engine
	uc     usecase.UserUsecase
	api.BaseApi
}

func (u *UserController) listHandler(c *gin.Context) {
	user, err := u.uc.FindAll()
	if err != nil {
		u.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	u.NewSuccessSingleResponse(c, user, "OK")
}

func (u *UserController) createHandler(c *gin.Context) {
	var payload request.CreateUser

	if err := c.ShouldBind(&payload); err != nil {
		u.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := u.uc.CreateUser(payload.Email, payload.Role); err != nil {
		u.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	u.NewSuccessSingleResponse(c, "", "OK")
}

func (u *UserController) updateHandler(c *gin.Context) {
	var payload model.User

	if err := c.ShouldBind(&payload); err != nil {
		u.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := u.uc.SaveData(&payload); err != nil {
		u.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	u.NewSuccessSingleResponse(c, payload, "OK")
}

func (u *UserController) deleteHandler(c *gin.Context) {
	id := c.Param("id")
	if err := u.uc.DeleteData(id); err != nil {
		u.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusNoContent, "")
}

func NewUserController(u *gin.Engine, uc usecase.UserUsecase) *UserController {
	controller := UserController{
		router: u,
		uc:     uc,
	}
	u.GET("/users", controller.listHandler)
	u.PUT("/users", controller.updateHandler)
	u.POST("/users", controller.createHandler)
	u.DELETE("/users/:id", controller.deleteHandler)
	return &controller
}
