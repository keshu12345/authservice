package utils

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type SuccessResponse struct {
	Timestamp string      `json:"timestamp"`
	Code      int         `json:"code"`
	Status    string      `json:"status"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Timestamp string `json:"timestamp"`
	Code      int    `json:"code"`
	Status    string `json:"status"`
	Error     string `json:"error"`
}

func RespondSuccess(c *gin.Context, statusCode int, message string, data interface{}) {
	resp := SuccessResponse{
		Timestamp: time.Now().Format(time.RFC3339),
		Code:      statusCode,
		Status:    http.StatusText(statusCode),
		Message:   message,
		Data:      data,
	}
	c.JSON(statusCode, resp)
}

func RespondError(c *gin.Context, statusCode int, err error) {
	resp := ErrorResponse{
		Timestamp: time.Now().Format(time.RFC3339),
		Code:      statusCode,
		Status:    http.StatusText(statusCode),
		Error:     err.Error(),
	}
	c.JSON(statusCode, resp)
}
