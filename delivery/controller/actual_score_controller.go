package controller

import (
	"net/http"
	"strings"

	"calibration-system.com/delivery/api"
	"calibration-system.com/model"
	"calibration-system.com/usecase"
	"github.com/gin-gonic/gin"
)

type ActualScoreController struct {
	router *gin.Engine
	uc     usecase.ActualScoreUsecase
	api.BaseApi
}

func (r *ActualScoreController) listHandler(c *gin.Context) {
	actualScores, err := r.uc.FindAll()
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, actualScores, "OK")
}

func (r *ActualScoreController) getByIdHandler(c *gin.Context) {
	id := c.Param("id")
	groupActualScores, err := r.uc.FindById(id)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, groupActualScores, "OK")
}

func (r *ActualScoreController) createHandler(c *gin.Context) {
	var payload model.ActualScore
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

func (r *ActualScoreController) updateHandler(c *gin.Context) {
	var payload model.ActualScore

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

func (r *ActualScoreController) deleteHandler(c *gin.Context) {
	projectId := c.Param("projectId")
	employeeId := c.Param("employeeId")
	if err := r.uc.DeleteData(projectId, employeeId); err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusNoContent, "")
}

func (r *ActualScoreController) uploadHandler(c *gin.Context) {
	// Menerima file Excel dari permintaan HTTP POST
	projectId := c.Request.FormValue("projectId")
	file, err := c.FormFile("excelFile")
	if err != nil {
		r.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logs, err := r.uc.BulkInsert(file, projectId)
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

func NewActualScoreController(r *gin.Engine, uc usecase.ActualScoreUsecase) *ActualScoreController {
	controller := ActualScoreController{
		router: r,
		uc:     uc,
	}
	r.GET("/actual-scores", controller.listHandler)
	r.GET("/actual-scores/:id", controller.getByIdHandler)
	r.PUT("/actual-scores", controller.updateHandler)
	r.POST("/actual-scores", controller.createHandler)
	r.POST("/actual-scores/upload", controller.uploadHandler)
	r.DELETE("/actual-scores/:projectId/:employeeId", controller.deleteHandler)
	return &controller
}