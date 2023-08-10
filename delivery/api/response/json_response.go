package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SendSingleResponse(c *gin.Context, data interface{}, desc string) {
	c.JSON(http.StatusOK, &SingleResponse{
		Status: Status{
			Code:        http.StatusOK,
			Description: desc,
		},
		Data: data,
	})
}

func SendErrorResponse(c *gin.Context, code int, desc string) {
	c.AbortWithStatusJSON(code, &SingleResponse{
		Status: Status{
			Code:        code,
			Description: desc,
		},
	})
}
