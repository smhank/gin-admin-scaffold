package middleware

import (
	"bytes"
	"gin-admin-base/internal/infras/global"
	"gin-admin-base/internal/infras/queue"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// OperationLogMiddleware 操作日志中间件，通过消息队列异步记录API请求
func OperationLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 只记录 /api/ 开头的请求，排除操作日志自身的查询和删除
		path := c.Request.URL.Path
		if !strings.HasPrefix(path, "/api/") || strings.HasPrefix(path, "/api/operation-logs") {
			c.Next()
			return
		}

		// 只记录写操作（POST, PUT, DELETE, PATCH）
		method := c.Request.Method
		if method == "GET" {
			c.Next()
			return
		}

		startTime := time.Now()

		// 读取请求体
		var requestBody string
		if c.Request.Body != nil {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err != nil {
				global.Logger.Warn("读取请求体失败", zap.Error(err))
			} else {
				if len(bodyBytes) > 2048 {
					requestBody = string(bodyBytes[:2048]) + "..."
				} else {
					requestBody = string(bodyBytes)
				}
			}
			// 始终重新设置请求体，以便后续 handler 可以读取
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// 创建自定义 ResponseWriter 来捕获响应
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// 处理请求
		c.Next()

		// 计算耗时
		duration := time.Since(startTime).Milliseconds()

		// 获取操作结果
		result := "success"
		if c.Writer.Status() >= 400 {
			result = "fail"
		}

		// 获取操作人
		username := "unknown"
		if u, exists := c.Get("username"); exists {
			username = u.(string)
		} else {
			token := c.GetHeader("Authorization")
			if token != "" {
				if strings.HasPrefix(token, "mock-jwt-token-") {
					username = strings.TrimPrefix(token, "mock-jwt-token-")
				}
			}
		}

		// 生成操作动作描述
		action := generateAction(method, path)

		// 通过消息队列异步写入操作日志
		if global.MsgRouter != nil {
			payload := map[string]interface{}{
				"username":   username,
				"action":     action,
				"method":     method,
				"path":       path,
				"params":     requestBody,
				"result":     result,
				"duration":   duration,
				"ip":         c.ClientIP(),
				"user_agent": c.Request.UserAgent(),
			}
			if err := global.MsgRouter.Publish(queue.TopicOperationLog, payload); err != nil {
				global.Logger.Warn("发布操作日志消息失败", zap.Error(err))
			}
		} else {
			global.Logger.Warn("消息队列未初始化，操作日志丢失")
		}
	}
}

// bodyLogWriter 自定义 ResponseWriter，用于捕获响应体
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// generateAction 根据方法和路径生成操作描述
func generateAction(method, path string) string {
	// 提取路径的最后一段作为操作对象
	parts := strings.Split(strings.TrimRight(path, "/"), "/")
	resource := ""
	if len(parts) > 0 {
		resource = parts[len(parts)-1]
		// 如果最后一段是数字ID，取前一段
		if isNumeric(resource) && len(parts) > 1 {
			resource = parts[len(parts)-2]
		}
	}

	actionMap := map[string]string{
		"POST":   "创建",
		"PUT":    "更新",
		"DELETE": "删除",
		"PATCH":  "修改",
	}

	action := actionMap[method]
	if action == "" {
		action = method
	}

	// 根据路径生成更友好的描述
	nameMap := map[string]string{
		"users":       "用户",
		"roles":       "角色",
		"permissions": "权限",
		"menus":       "菜单",
		"paths":       "API路径",
		"login":       "登录",
	}

	if name, ok := nameMap[resource]; ok {
		return action + name
	}

	return action + resource
}

func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return len(s) > 0
}
