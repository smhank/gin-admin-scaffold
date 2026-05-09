package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	Enabled  bool          // 是否启用
	Rate     int           // 时间窗口内允许的最大请求数
	Window   time.Duration // 时间窗口大小
	Burst    int           // 突发请求数（令牌桶容量）
	PerRoute bool          // 是否按路由分别限流
}

// DefaultRateLimitConfig 默认限流配置
var DefaultRateLimitConfig = RateLimitConfig{
	Enabled:  true,
	Rate:     100,
	Window:   time.Second,
	Burst:    200,
	PerRoute: true,
}

// visitor 访问者记录
type visitor struct {
	tokens    float64   // 当前令牌数
	lastCheck time.Time // 上次检查时间
}

// rateLimiter 限流器
type rateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	config   RateLimitConfig
}

// globalRateLimiter 全局限流器实例
var globalRateLimiter = &rateLimiter{
	visitors: make(map[string]*visitor),
	config:   DefaultRateLimitConfig,
}

// cleanup 定期清理过期记录
func (rl *rateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	now := time.Now()
	for key, v := range rl.visitors {
		if now.Sub(v.lastCheck) > rl.config.Window*2 {
			delete(rl.visitors, key)
		}
	}
}

// allow 检查是否允许请求
func (rl *rateLimiter) allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	v, exists := rl.visitors[key]
	if !exists {
		// 新访问者，初始为 Burst 个令牌
		rl.visitors[key] = &visitor{
			tokens:    float64(rl.config.Burst),
			lastCheck: now,
		}
		v = rl.visitors[key]
	}

	// 计算时间窗口内应补充的令牌数
	elapsed := now.Sub(v.lastCheck).Seconds()
	rate := float64(rl.config.Rate) / rl.config.Window.Seconds()
	v.tokens += elapsed * rate

	// 限制令牌数不超过 Burst
	if v.tokens > float64(rl.config.Burst) {
		v.tokens = float64(rl.config.Burst)
	}

	v.lastCheck = now

	// 检查是否有可用令牌
	if v.tokens >= 1 {
		v.tokens--
		return true
	}

	return false
}

// getKey 生成限流 key
func getKey(c *gin.Context, perRoute bool) string {
	ip := c.ClientIP()
	if perRoute {
		return ip + ":" + c.FullPath()
	}
	return ip
}

// RateLimitMiddleware 基于 IP 和路由的限流中间件
// 使用令牌桶算法，支持按路由分别限流
func RateLimitMiddleware(config ...RateLimitConfig) gin.HandlerFunc {
	rl := globalRateLimiter
	if len(config) > 0 {
		rl = &rateLimiter{
			visitors: make(map[string]*visitor),
			config:   config[0],
		}
	}

	// 启动定期清理（仅全局实例执行一次）
	if rl == globalRateLimiter {
		go func() {
			ticker := time.NewTicker(time.Minute)
			defer ticker.Stop()
			for range ticker.C {
				rl.cleanup()
			}
		}()
	}

	return func(c *gin.Context) {
		if !rl.config.Enabled {
			c.Next()
			return
		}

		key := getKey(c, rl.config.PerRoute)

		if !rl.allow(key) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code":  429,
				"msg":   "请求过于频繁，请稍后再试",
				"retry": rl.config.Window.Seconds(),
			})
			return
		}

		c.Next()
	}
}

// SetRateLimitConfig 设置全局限流配置
func SetRateLimitConfig(config RateLimitConfig) {
	globalRateLimiter.mu.Lock()
	defer globalRateLimiter.mu.Unlock()
	globalRateLimiter.config = config
}
