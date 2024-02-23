package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

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
	topRemarks, err := r.uc.FindByForeignKeyID(projectID, employeeID, projectPhaseID)

	for _, data := range topRemarks {
		if data.EvidenceName != "" {
			// data.EvidenceLink = fmt.Sprintf("http://%s/view-initiative/%s", c.Request.Host, data.ID)
			data.EvidenceLink = fmt.Sprintf("http://%s/view-initiative/%s", c.Request.Host, data.ID)
		}
	}

	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, topRemarks, "OK")
}

func (r *TopRemarkController) viewFileHandler(c *gin.Context) {
	id := c.Param("id")
	topRemarks, err := r.uc.FindById(id)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var contentType string
	var fileExtension = strings.ToLower(strings.TrimPrefix(filepath.Ext(topRemarks.EvidenceName), "."))
	switch fileExtension {
	case "pdf":
		contentType = "application/pdf"
	case "doc", "docx":
		contentType = "application/msword"
	case "xlsx", "xls":
		contentType = "application/vnd.ms-excel"
	case "jpg", "jpeg":
		contentType = "image/jpeg"
	case "png":
		contentType = "image/png"
	// Add cases for other file types as needed
	default:
		// Default to octet-stream for unknown file types
		contentType = "application/octet-stream"
	}

	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=%s", topRemarks.EvidenceName))

	c.Data(http.StatusOK, contentType, topRemarks.Evidence)
	r.NewSuccessSingleResponse(c, topRemarks, "OK")
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
	payloadsJSON := c.Request.FormValue("payloads")
	var payloads []*model.TopRemark
	if err := json.Unmarshal([]byte(payloadsJSON), &payloads); err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	for index, payloadData := range payloads {
		file, header, err := c.Request.FormFile(fmt.Sprintf("Evidence_%d", index))
		if err != nil && err.Error() != "http: no such file" {
			// Handle file upload error for a specific payload
			r.NewFailedResponse(c, http.StatusInternalServerError, fmt.Sprintf("file open err: %s %t", err.Error(), file == nil))
			return
		}
		if file != nil {
			var fileBuffer bytes.Buffer
			_, err = io.Copy(&fileBuffer, file)
			if err != nil {
				// Handle error while copying the file data
				r.NewFailedResponse(c, http.StatusInternalServerError, fmt.Sprintf("file copy err: %s", err.Error()))
				return
			}
			fileBytes := fileBuffer.Bytes()
			payloadData.EvidenceName = header.Filename
			payloadData.Evidence = fileBytes
		} else {
			payloadData.EvidenceName = ""
		}

	}

	if err := r.uc.SaveDataByProject(payloads); err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, fmt.Sprintf("save data error err: %s", err.Error()))
		return
	}

	r.NewSuccessSingleResponse(c, payloads, "OK")
}

func (r *TopRemarkController) updateHandler(c *gin.Context) {
	var payload model.TopRemark

	if err := c.ShouldBind(&payload); err != nil {
		r.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		r.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		r.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	payload.EvidenceName = header.Filename
	payload.Evidence = fileBytes
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

	auth := r.Group("/auth").Use(middleware.NewTokenValidator(tokenService).RequireToken())
	auth.GET("/top-remark/:projectID/:employeeID/:projectPhaseID", controller.getByIdHandler)
	r.GET("/view-initiative/:id", controller.viewFileHandler)
	auth.PUT("/top-remark", controller.updateHandler)
	auth.POST("/top-remark", controller.createHandler)
	auth.POST("/top-remark/project", controller.createHandlerByProject)
	auth.POST("/top-remark/delete", controller.deleteHandlerByProject)
	auth.DELETE("/top-remark/:id", controller.deleteHandler)
	return &controller
}
