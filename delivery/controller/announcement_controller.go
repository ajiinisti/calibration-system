package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"calibration-system.com/delivery/api"
	"calibration-system.com/delivery/middleware"
	"calibration-system.com/model"
	"calibration-system.com/usecase"
	"calibration-system.com/utils/authenticator"
	"github.com/gin-gonic/gin"
)

type AnnouncementController struct {
	router       *gin.Engine
	uc           usecase.AnnouncementUsecase
	tokenService authenticator.AccessToken
	api.BaseApi
}

func (r *AnnouncementController) listHandler(c *gin.Context) {
	announcements, err := r.uc.FindAll()
	for _, data := range announcements {
		// data.FileLink = fmt.Sprintf("http://%s/announcement/images/%s", c.Request.Host, data.ID)
		data.FileLink = fmt.Sprintf("http://%sannouncement/images/%s", c.Request.Host, data.ID)
	}
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, announcements, "OK")
}

func (r *AnnouncementController) listActiveHandler(c *gin.Context) {
	announcements, err := r.uc.FindAllActive()
	for _, data := range announcements {
		// data.FileLink = fmt.Sprintf("http://%s/announcement/images/%s", c.Request.Host, data.ID)
		data.FileLink = fmt.Sprintf("http://%sannouncement/images/%s", c.Request.Host, data.ID)
	}
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, announcements, "OK")
}

func (r *AnnouncementController) getByHandler(c *gin.Context) {
	name := c.Param("name")
	announcements, err := r.uc.FindByName(name)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, announcements, "OK")
}

func (r *AnnouncementController) getByIdHandler(c *gin.Context) {
	id := c.Param("id")
	announcements, err := r.uc.FindById(id)
	// announcements.FileLink = fmt.Sprintf("http://%s/announcement/images/%s", c.Request.Host, announcements.ID)
	announcements.FileLink = fmt.Sprintf("http://%sannouncement/images/%s", c.Request.Host, announcements.ID)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	r.NewSuccessSingleResponse(c, announcements, "OK")
}

func (r *AnnouncementController) getAnnouncementImage(c *gin.Context) {
	id := c.Param("id")
	announcements, err := r.uc.FindById(id)
	if err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Header("Content-Type", "image/jpeg")                                                      // Adjust content type based on your image type
	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=%s", announcements.FileName)) // Suggest the file name
	c.Data(200, "image/jpeg", announcements.File)
	r.NewSuccessSingleResponse(c, announcements, "OK")
}

func (r *AnnouncementController) createHandler(c *gin.Context) {
	var payload model.Announcement
	payloadJSON := c.Request.FormValue("payload")
	if err := json.Unmarshal([]byte(payloadJSON), &payload); err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	file, header, err := c.Request.FormFile("file")
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
		payload.FileName = header.Filename
		payload.File = fileBytes
	}

	if err := r.uc.SaveData(&payload); err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	r.NewSuccessSingleResponse(c, payload, "OK")
}

func (r *AnnouncementController) updateHandler(c *gin.Context) {
	var payload model.Announcement

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

func (r *AnnouncementController) deleteHandler(c *gin.Context) {
	id := c.Param("id")
	if err := r.uc.DeleteData(id); err != nil {
		r.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusNoContent, "")
}

func NewAnnouncementController(r *gin.Engine, tokenService authenticator.AccessToken, uc usecase.AnnouncementUsecase) *AnnouncementController {
	controller := AnnouncementController{
		router:       r,
		tokenService: tokenService,
		uc:           uc,
	}

	auth := r.Group("/auth").Use(middleware.NewTokenValidator(tokenService).RequireToken())
	auth.GET("/announcements", controller.listHandler)
	auth.GET("/announcements/active", controller.listActiveHandler)
	// auth.GET("/announcements/:name", controller.getByHandler)
	auth.GET("/announcements/:id", controller.getByIdHandler)
	r.GET("/announcement/images/:id", controller.getAnnouncementImage)
	auth.PUT("/announcements", controller.updateHandler)
	auth.POST("/announcements", controller.createHandler)
	auth.DELETE("/announcements/:id", controller.deleteHandler)
	return &controller
}
