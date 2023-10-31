package transportutil

import (
	"net/http"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/apperror"

	"github.com/gin-gonic/gin"
)

type BaseResponse struct {
	Code         int32       `json:"code"`
	Data         interface{} `json:"data,omitempty"`
	Message      string      `json:"message,omitempty"`
	DebugMessage string      `json:"debugMessage,omitempty"`
	Timestamp    time.Time   `json:"timestamp"`
}

func SendJSONResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &BaseResponse{
		Data:      data,
		Timestamp: time.Now(),
	})
}

func Error(c *gin.Context, err error) {
	_ = c.Error(err)
}

const headerLang = "Accept-Language"

func HandleError(c *gin.Context) {
	c.Next()

	for _, err := range c.Errors {
		switch e := err.Err.(type) {
		case *apperror.AppError:
			c.AbortWithStatusJSON(http.StatusOK, &BaseResponse{
				Code:         int32(e.Code),
				Message:      e.Translate(c.GetHeader(headerLang)),
				DebugMessage: e.Error(),
				Timestamp:    time.Now(),
			})
		default:
			c.AbortWithStatusJSON(http.StatusOK, &BaseResponse{
				Code:         int32(apperror.ErrUnknown.Code),
				Message:      apperror.ErrUnknown.Translate(c.GetHeader(headerLang)),
				DebugMessage: e.Error(),
				Timestamp:    time.Now(),
			})
		}

		return
	}
}
