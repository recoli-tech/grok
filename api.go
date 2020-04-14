package grok

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HTTPCreated ...
func HTTPCreated(ctx *gin.Context, id string, content interface{}) {
	ctx.Writer.Header().Set("Location", fmt.Sprintf("%s/%s", ctx.Request.URL.Path, id))
	ctx.JSON(http.StatusCreated, content)
}
