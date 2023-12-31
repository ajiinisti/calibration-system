package controller

import (
	"net/http"
	"strconv"
	"strings"

	"calibration-system.com/delivery/api"
	"calibration-system.com/delivery/api/request"
	"calibration-system.com/delivery/middleware"
	"calibration-system.com/model"
	"calibration-system.com/usecase"
	"calibration-system.com/utils/authenticator"
	"github.com/gin-gonic/gin"
)

type BusinessUnitController struct {
	router       *gin.Engine
	uc           usecase.BusinessUnitUsecase
	tokenService authenticator.AccessToken
	api.BaseApi
}

func (r *BusinessUnitController) listHandler(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		r.NewFailedResponse(c, http.StatusBadRequest, "Invalid page number")
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		r.NewFailedResponse(c, http.StatusBadRequest, "Invalid limit number")
	}

	nameQuery := c.Query("name")
	param := request.PaginationParam{
		Page:   page,
		Limit:  limit,
		Offset: 0,
		Name:   nameQuery,
	}

	businessUnits, pagination, err := r.uc.FindPagination(param)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var newBu []interface{}
	for _, v := range businessUnits {
		newBu = append(newBu, v)
	}

	r.NewSuccesPagedResponse(c, newBu, "OK", pagination)
}

func (r *BusinessUnitController) getByIdHandler(c *gin.Context) {
	id := c.Param("id")
	groupBusinessUnits, err := r.uc.FindById(id)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, groupBusinessUnits, "OK")
}

func (r *BusinessUnitController) getAllHandler(c *gin.Context) {
	groupBusinessUnits, err := r.uc.FindAll()
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, groupBusinessUnits, "OK")
}

func (r *BusinessUnitController) createHandler(c *gin.Context) {
	var payload model.BusinessUnit
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

func (r *BusinessUnitController) updateHandler(c *gin.Context) {
	var payload model.BusinessUnit

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

func (r *BusinessUnitController) deleteHandler(c *gin.Context) {
	id := c.Param("id")
	if err := r.uc.DeleteData(id); err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusNoContent, "")
}

func (r *BusinessUnitController) uploadHandler(c *gin.Context) {
	// Menerima file Excel dari permintaan HTTP POST
	file, err := c.FormFile("excelFile")
	if err != nil {
		r.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logs, err := r.uc.BulkInsert(file)
	if err != nil {
		if len(logs) > 0 {
			r.NewFailedResponse(c, http.StatusInternalServerError, strings.Join(logs, "."))
		} else {
			r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	r.NewSuccessSingleResponse(c, "", "OK")
}

func NewBusinessUnitController(r *gin.Engine, tokenService authenticator.AccessToken, uc usecase.BusinessUnitUsecase) *BusinessUnitController {
	controller := BusinessUnitController{
		router:       r,
		tokenService: tokenService,
		uc:           uc,
	}
	auth := r.Group("/auth").Use(middleware.NewTokenValidator(tokenService).RequireToken())
	auth.GET("/business-units", controller.listHandler)
	auth.GET("/business-units/all", controller.getAllHandler)
	auth.GET("/business-units/:id", controller.getByIdHandler)
	auth.PUT("/business-units", controller.updateHandler)
	auth.POST("/business-units", controller.createHandler)
	auth.POST("/business-units/upload", controller.uploadHandler)
	auth.DELETE("/business-units/:id", controller.deleteHandler)
	return &controller
}
