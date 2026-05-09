package context

import (
	"gin-admin-base/internal/infras/query"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AppContext struct {
	*gin.Context
	Query *query.Query
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func (c *AppContext) Success(data interface{}) {
	c.JSON(200, Response{Code: 0, Msg: "success", Data: data})
}

func (c *AppContext) Error(code int, msg string) {
	c.JSON(200, Response{Code: code, Msg: msg})
}

type AppHandler func(ctx *AppContext) error

func HandlerAdapter(h AppHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := &AppContext{
			Context: c,
			Query:   query.Use(c.MustGet("db").(*gorm.DB)), // 假设 db 已注入到 gin context
		}
		if err := h(ctx); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
		}
	}
}
