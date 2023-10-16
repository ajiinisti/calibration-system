package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"calibration-system.com/delivery/api"
	"calibration-system.com/delivery/api/request"
	"calibration-system.com/model"
	"calibration-system.com/usecase"
	"github.com/gin-gonic/gin"
)

type RemarkSettingController struct {
	router *gin.Engine
	uc     usecase.RemarkSettingUsecase
	api.BaseApi
}

func (r *RemarkSettingController) listHandler(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		r.NewFailedResponse(c, http.StatusBadRequest, "Invalid page number")
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		r.NewFailedResponse(c, http.StatusBadRequest, "Invalid limit number")
	}

	projectId := c.Query("id")
	param := request.PaginationParam{
		Page:   page,
		Limit:  limit,
		Offset: 0,
	}

	RemarkSettings, pagination, err := r.uc.FindPagination(param, projectId)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var newRatings []interface{}
	for _, v := range RemarkSettings {
		newRatings = append(newRatings, v)
	}

	r.NewSuccesPagedResponse(c, newRatings, "OK", pagination)
}

func (r *RemarkSettingController) getByIdHandler(c *gin.Context) {
	projectId := c.Param("projectId")
	RemarkSettings, err := r.uc.FindById(projectId)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, RemarkSettings, "OK")
}

func (r *RemarkSettingController) createHandler(c *gin.Context) {
	var payload model.RemarkSetting
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

func (r *RemarkSettingController) createHandlerByProject(c *gin.Context) {
	var payload []*model.RemarkSetting

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

func (r *RemarkSettingController) updateHandler(c *gin.Context) {
	var payload model.RemarkSetting

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

func (r *RemarkSettingController) deleteHandler(c *gin.Context) {
	id := c.Param("id")
	if err := r.uc.DeleteData(id); err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusNoContent, "")
}

func (r *RemarkSettingController) deleteHandlerByProject(c *gin.Context) {
	var payload request.DeleteRemark
	if err := c.ShouldBind(&payload); err != nil {
		r.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	fmt.Println("PAYLOAD := ", payload.IDs)

	if err := r.uc.BulkDeleteData(payload); err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusNoContent, "")
}

func NewRemarkSettingController(r *gin.Engine, uc usecase.RemarkSettingUsecase) *RemarkSettingController {
	controller := RemarkSettingController{
		router: r,
		uc:     uc,
	}
	r.GET("/remark-settings", controller.listHandler)
	r.GET("/remark-settings/:projectId", controller.getByIdHandler)
	r.PUT("/remark-settings", controller.updateHandler)
	r.POST("/remark-settings", controller.createHandler)
	r.POST("/remark-settings/project", controller.createHandlerByProject)
	r.POST("/remark-settings/delete", controller.deleteHandlerByProject)
	r.DELETE("/remark-settings/:id", controller.deleteHandler)
	return &controller
}
