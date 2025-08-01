package response

import "github.com/gin-gonic/gin"

type Result struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(200, Result{Code: 200, Message: "success", Data: data})
}

func Fail(c *gin.Context, msg string) {
	c.JSON(400, Result{Code: 400, Message: msg})
}

func Unauthorized(c *gin.Context, msg string) {
	c.JSON(401, Result{Code: 401, Message: msg})
}

func Forbidden(c *gin.Context, msg string) {
	c.JSON(403, Result{Code: 403, Message: msg})
}

func InternalError(c *gin.Context, msg string) {
	c.JSON(500, Result{Code: 500, Message: msg})
}
