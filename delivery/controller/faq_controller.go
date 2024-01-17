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

type FaqController struct {
	router       *gin.Engine
	uc           usecase.FaqUsecase
	tokenService authenticator.AccessToken
	api.BaseApi
}

func (r *FaqController) listHandler(c *gin.Context) {
	faqs, err := r.uc.FindAll()
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, faqs, "OK")
}

func (r *FaqController) listActiveHandler(c *gin.Context) {
	faqs, err := r.uc.FindAllActive()
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, faqs, "OK")
}

func (r *FaqController) getByHandler(c *gin.Context) {
	name := c.Param("name")
	faqs, err := r.uc.FindByName(name)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, faqs, "OK")
}

func (r *FaqController) getByIdHandler(c *gin.Context) {
	id := c.Param("id")
	faqs, err := r.uc.FindById(id)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, faqs, "OK")
}

func (r *FaqController) createHandler(c *gin.Context) {
	var payload model.Faq
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

func (r *FaqController) updateHandler(c *gin.Context) {
	var payload model.Faq
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

func (r *FaqController) deleteHandler(c *gin.Context) {
	id := c.Param("id")
	if err := r.uc.DeleteData(id); err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusNoContent, "")
}

func NewFaqController(r *gin.Engine, tokenService authenticator.AccessToken, uc usecase.FaqUsecase) *FaqController {
	controller := FaqController{
		router:       r,
		tokenService: tokenService,
		uc:           uc,
	}

	auth := r.Group("/auth").Use(middleware.NewTokenValidator(tokenService).RequireToken())
	auth.GET("/faqs", controller.listHandler)
	auth.GET("/faqs/active", controller.listActiveHandler)
	// auth.GET("/faqs/:name", controller.getByHandler)
	auth.GET("/faqs/:id", controller.getByIdHandler)
	auth.PUT("/faqs", controller.updateHandler)
	auth.POST("/faqs", controller.createHandler)
	auth.DELETE("/faqs/:id", controller.deleteHandler)
	return &controller
}
