package grok

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Swagger ...
func Swagger(file string) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := os.Open(file)

		if err != nil {
			logrus.WithError(err).
				Error("swagger file error")

			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		defer file.Close()

		content, _ := ioutil.ReadAll(file)

		c.Writer.Write(content)
		c.Writer.Header().Set("content-type", "application/json")
	}
}
