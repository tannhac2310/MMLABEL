package interceptor

import (
	"bytes"
	"io/ioutil"
	"strings"
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Ginzap returns a gin.HandlerFunc (middleware) that logs requests using uber-go/zap.
//
// Requests with errors are logged using zap.Error().
// Requests without errors are logged using zap.Info().
//
func Ginzap(logger *zap.Logger, decider func(c *gin.Context) bool, logReq, logResp bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !decider(c) {
			c.Next()
			return
		}

		// some evil middlewares modify this values
		path := c.Request.URL.Path

		fields := []zapcore.Field{
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("requestID", requestid.Get(c)),
		}

		if logReq {
			if !strings.Contains(c.GetHeader("Content-Type"), "boundary=----WebKitFormBoundary") {
				body, _ := ioutil.ReadAll(c.Request.Body)
				c.Request.Body = ioutil.NopCloser(bytes.NewReader(body))
				fields = append(fields, zap.String("request", string(body)))
			}
		}

		blw := &bodyLogWriter{}
		if logResp {
			// wrap log response
			blw = &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
			c.Writer = blw
		}

		c.Request = c.Request.WithContext(ctxzap.ToContext(c.Request.Context(), logger.With(fields...)))

		start := time.Now()
		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		newLogger := ctxzap.Extract(c.Request.Context()).With(
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", latency),
		)

		if userID := UserIDFromCtx(c); userID != "" {
			newLogger = newLogger.With(zap.String("userID", userID))
		}

		if logResp {
			newLogger = newLogger.With(zap.String("response", blw.body.String()))
		}

		if len(c.Errors) > 0 {
			// Append error field if this is an erroneous request.
			for _, e := range c.Errors.Errors() {
				newLogger.Error(e)
			}
		} else {
			newLogger.Info("request finished")
		}
	}
}
