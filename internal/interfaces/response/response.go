package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`               // 业务码：0 成功，非 0 失败
	Message string      `json:"message"`            // 提示信息
	Data    interface{} `json:"data,omitempty"`     // 数据（成功时）
	TraceID string      `json:"trace_id,omitempty"` // 请求追踪 ID
}

// PageData 分页数据
type PageData struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMessage 成功响应（自定义消息）
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: message,
		Data:    data,
	})
}

// SuccessPage 分页成功响应
func SuccessPage(c *gin.Context, list interface{}, total int64, page int, pageSize int) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data: PageData{
			List:     list,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
	})
}

// Error 错误响应
func Error(c *gin.Context, httpStatus int, code int, message string) {
	c.AbortWithStatusJSON(httpStatus, Response{
		Code:    code,
		Message: message,
	})
}

// BadRequest 参数错误 400
func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, 1, message)
}

// Unauthorized 未授权 401
func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, 1, message)
}

// Forbidden 权限不足 403
func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, 1, message)
}

// NotFound 资源不存在 404
func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, 1, message)
}

// TooManyRequests 请求过于频繁 429
func TooManyRequests(c *gin.Context, message string) {
	Error(c, http.StatusTooManyRequests, 429, message)
}

// InternalError 服务器内部错误 500
func InternalError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, 1, message)
}
