package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandlerFuncResError to make router handler what return Result and error
type HandlerFuncResError func(*gin.Context) (Result, error)

// GinHandler wrap HandlerFuncResError to easily handle and display errors nicely
func GinHandler(myHandler HandlerFuncResError) gin.HandlerFunc {
	return func(c *gin.Context) {
		if res, err := myHandler(c); err != nil {
			res.UpdateAllFields(http.StatusInternalServerError, err.Error(), nil)
			c.AbortWithStatusJSON(http.StatusInternalServerError, res)
		} else {
			c.JSON(http.StatusOK, res)
		}
	}
}
