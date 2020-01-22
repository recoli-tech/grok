package grok

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

//LogMiddleware ...
func LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer recovery()
		defer c.Request.Body.Close()

		requestID := uuid.New()

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		blw.Header().Set("Request-Id", requestID.String())
		c.Writer = blw

		now := time.Now()
		req := request(c)

		c.Next()

		elapsed := time.Since(now)
		fields := make(map[string]interface{})

		fields["request"] = req
		fields["claims"] = c.Keys
		fields["errors"] = c.Errors
		fields["ip"] = c.ClientIP()
		fields["latency"] = elapsed.Seconds()
		fields["request_id"] = requestID.String()
		fields["response"] = response(blw)

		logrus.WithFields(fields).Infof(
			"Request incoming from %s elapsed %s completed with %d",
			c.ClientIP(),
			elapsed.String(),
			c.Writer.Status(),
		)
	}
}

func request(context *gin.Context) interface{} {
	r := make(map[string]interface{})

	bodyCopy := new(bytes.Buffer)
	io.Copy(bodyCopy, context.Request.Body)
	bodyData := bodyCopy.Bytes()

	var body map[string]interface{}
	json.Unmarshal(bodyData, &body)

	r["body"] = body
	r["host"] = context.Request.Host
	r["form"] = context.Request.Form
	r["path"] = context.Request.URL.Path
	r["method"] = context.Request.Method
	r["headers"] = context.Request.Header
	r["url"] = context.Request.URL.String()
	r["post_form"] = context.Request.PostForm
	r["remote_addr"] = context.Request.RemoteAddr
	r["query_string"] = context.Request.URL.Query()

	context.Request.Body = ioutil.NopCloser(bytes.NewReader(bodyData))

	return r
}

func response(writer *bodyLogWriter) interface{} {
	r := make(map[string]interface{})

	var body map[string]interface{}
	json.Unmarshal(writer.body.Bytes(), &body)

	r["body"] = body
	r["status"] = writer.Status()
	r["headers"] = writer.Header()

	return r
}

func recovery() {
	if err := recover(); err != nil {
		logrus.WithField("error", err).Error("Error on logging middleware")
	}
}
