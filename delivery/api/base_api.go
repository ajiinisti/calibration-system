package api

import (
	"calibration-system.com/delivery/api/response"
	"github.com/gin-gonic/gin"
)

type BaseApi struct{}

func (b *BaseApi) ParseRequestBody(c *gin.Context, payload interface{}) error {
	if err := c.ShouldBind(payload); err != nil {
		return err
	}
	return nil
}

func (b *BaseApi) NewSuccessSingleResponse(c *gin.Context, data interface{}, desc string) {
	response.SendSingleResponse(c, data, desc)
}

func (b *BaseApi) NewFailedResponse(c *gin.Context, code int, desc string) {
	response.SendErrorResponse(c, code, desc)
}
