package handlers

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Message string `json:"message"`
	Code    string `json:"ziyad_error_code"`
	TraceID string `json:"trace_id"`
}

func newTraceID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "trace-unavailable"
	}
	return hex.EncodeToString(b)
}

func WriteError(c *gin.Context, status int, msg, code string) {
	c.JSON(status, ErrorResponse{
		Message: msg,
		Code:    code,
		TraceID: newTraceID(),
	})
}
