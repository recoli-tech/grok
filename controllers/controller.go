package controllers

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/raafvargas/grok/utils"
)

//Controller ...
type Controller interface {
	RegisterRoutes(*gin.RouterGroup)
}

//BindingError ...
func BindingError(context *gin.Context, err error) {
	context.Error(err)
	context.JSON(http.StatusBadRequest, utils.NewError(http.StatusBadRequest, err.Error()))
}

//ResolveError ...
func ResolveError(context *gin.Context, err error) {
	context.Error(err)

	if utils.DefaultErrorMapping.Exists(err) {
		err = utils.DefaultErrorMapping.Get(err)
	}

	if reflect.TypeOf(err) != reflect.TypeOf(&utils.Error{}) {
		context.Status(http.StatusInternalServerError)
		return
	}

	status := http.StatusBadRequest
	message := err.(*utils.Error)

	if message.Code != 0 {
		status = message.Code
	}

	context.JSON(status, message)
}
