package utils

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "sync"
    "time"
)

type RateLimiter struct {
    MaxRequests  int           // 允许的最大请求数
    Interval     time.Duration // 限制的时间间隔
    LastRequest  time.Time     // 上一次请求的时间
    Mutex        sync.Mutex    // 用于保护数据的互斥锁
    RequestCount int           // 记录当前时间窗口内的请求次数
}

func RequestRateLimit(rateLimiter *RateLimiter) gin.HandlerFunc {
    return func(c *gin.Context) {
        rateLimiter.Mutex.Lock()
        defer rateLimiter.Mutex.Unlock()

        // 检查是否超过了请求数限制
        if rateLimiter.RequestCount >= rateLimiter.MaxRequests {
            c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
            c.Abort()
            return
        }

        // 检查是否超过了时间间隔，如果超过了就重置请求计数
        if time.Since(rateLimiter.LastRequest) >= rateLimiter.Interval {
            rateLimiter.LastRequest = time.Now()
            rateLimiter.RequestCount = 1
        } else {
            rateLimiter.RequestCount++
        }

        c.Next()
    }
}

func InitRequestRateLimit(r *gin.Engine, maxRequests int) {
	rateLimiter := &RateLimiter{
        MaxRequests:  maxRequests,
        Interval:     time.Second,
        LastRequest:  time.Now(),
    }

	r.Use(RequestRateLimit(rateLimiter))
}

