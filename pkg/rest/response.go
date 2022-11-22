package rest

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"top-ping/pkg/baseerr"
	"top-ping/pkg/utils"
)

var R = NewResponse()

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Details []string    `json:"details"`
	TraceID string      `json:"traceId"`
}

func NewResponse() *Response {
	return &Response{}
}

func (r *Response) Success(c *gin.Context, data interface{}) {
	if data == nil {
		data = gin.H{}
	}

	c.JSON(http.StatusOK, &Response{
		Code:    baseerr.Success.Code(),
		Message: baseerr.Success.Msg(),
		Data:    data,
		Details: []string{},
		TraceID: getTraceId(c.Request.Context()),
	})
}

func (r *Response) Error(c *gin.Context, error error) {
	if error != nil {
		if v, ok := error.(*baseerr.Error); ok {
			response := &Response{
				Code:    v.Code(),
				Message: v.Msg(),
				Data:    gin.H{},
				Details: []string{},
				TraceID: getTraceId(c.Request.Context()),
			}

			details := v.Details()
			if len(details) > 0 {
				response.Details = details
			}
			c.JSON(v.StatusCode(), response)
			return
		}
	}

	c.JSON(http.StatusOK, &Response{
		Code:    baseerr.Success.Code(),
		Message: baseerr.Success.Msg(),
		Data:    gin.H{},
		TraceID: getTraceId(c.Request.Context()),
	})
}

func getTraceId(c context.Context) string {
	var traceID string

	if v := c.Value(utils.TraceKey); v != nil {
		if t, ok := v.(string); ok {
			traceID = t
		}
	}

	return traceID
}
