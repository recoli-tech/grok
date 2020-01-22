package grok

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

//APIController ...
type APIController interface {
	RegisterRoutes(*gin.RouterGroup)
}

//BindingError ...
func BindingError(context *gin.Context, err error) {
	context.Error(err)
	context.JSON(http.StatusBadRequest, NewError(http.StatusBadRequest, err.Error()))
}

//ResolveError ...
func ResolveError(context *gin.Context, err error) {
	context.Error(err)

	if DefaultErrorMapping.Exists(err) {
		err = DefaultErrorMapping.Get(err)
	}

	if reflect.TypeOf(err) != reflect.TypeOf(&Error{}) {
		context.Status(http.StatusInternalServerError)
		return
	}

	status := http.StatusBadRequest
	message := err.(*Error)

	if message.Code != 0 {
		status = message.Code
	}

	context.JSON(status, message)
}
