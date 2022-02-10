package controller

import (
	"admin/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type baseController struct {

}

func (bcl baseController) GetUserInfo(c *gin.Context) models.Admins{
	token := c.Request.Header.Get("X-Token")
	if token == "" || len(token) == 0 {
		token = c.Query("token")
	}

	if UserAuth, ok := models.UserAuthSession[token]; ok {
		return UserAuth.AdminInfo
	}

	return models.Admins{}
}


func (bcl baseController) buildSuccessData(message string, data interface{}) map[string]interface{} {
	return gin.H{
		"code" : 1000,
		"message" : message,
		"data" : data,
	}
}

func (bcl baseController) buildErrorData(message string, data interface{}) map[string]interface{} {
	return gin.H{
		"code" : 0,
		"message" : message,
		"data" : data,
	}
}

func (bcl baseController) SuccessResponse(c *gin.Context,message string, data interface{}) {
	c.JSON(http.StatusOK, bcl.buildSuccessData(message, data))
}

func (bcl baseController) ErrorResponse(c *gin.Context,  message string, data interface{}) {
	c.JSON(http.StatusOK, bcl.buildErrorData(message, data))
}

func (bcl baseController) SuccessRedirect(c *gin.Context, message, url string) {
	c.Redirect(http.StatusFound, "/redirect/success?message=" + message + "&url=" + url )
}

func (bcl baseController) ErrorRedirect(c *gin.Context, message, url string) {
	c.Redirect(http.StatusFound, "/redirect/error?message=" + message + "&url=" + url )
}
