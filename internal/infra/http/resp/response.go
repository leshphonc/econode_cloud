package resp

import (
	"econode-cloud/internal/pkg/bizerr"
	"econode-cloud/internal/pkg/ctxx"
	"errors"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Data      any    `json:"data,omitempty"`
	TraceID   string `json:"trace_id,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}

func OK(c *gin.Context, data any) {
	c.JSON(200, Response{
		Code:      0,
		Message:   "ok",
		Data:      data,
		RequestID: ctxx.RequestID(c.Request.Context()),
	})
}

func Fail(c *gin.Context, err error) {
	traceID := ctxx.TraceID(c.Request.Context())
	requestID := ctxx.RequestID(c.Request.Context())

	var bizErr *bizerr.BizError
	if errors.As(err, &bizErr) {
		c.JSON(bizErr.HTTPStatus, Response{
			Code:      bizErr.Code,
			Message:   bizErr.Message,
			TraceID:   traceID,
			RequestID: requestID,
		})
		return
	}

	// 未知错误 → 500
	c.JSON(500, Response{
		Code:      9999,
		Message:   "internal server error",
		TraceID:   traceID,
		RequestID: requestID,
	})
}
