package model

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type InText struct {
	TextData string `form:"in_text" binding:"required,max=512"`
	NameData string `form:"name" binding:"max=62"`
	CSRF     string `form:"csrf"`
}

type IPAddress struct {
	ID        uint   `gorm:"primary_key" json:"id"`
	IPAddress string `json:"ip_address"`
	CreatedAt string `json:"created_at"`
}

type RateLimiter struct {
	mu          sync.Mutex
	lastRequest map[string]time.Time
	waitTime    time.Duration
	maxWaitTime time.Duration
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		lastRequest: make(map[string]time.Time),
		waitTime:    10 * time.Second,
		maxWaitTime: 120 * time.Second,
	}
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != "POST" {
			c.Next()
			return
		}

		clientIP := c.ClientIP()
		rl.mu.Lock()
		defer rl.mu.Unlock()

		lastTime, exists := rl.lastRequest[clientIP]
		var currentWaitTime time.Duration

		if exists && time.Since(lastTime) < rl.waitTime {
			currentWaitTime = rl.waitTime
			rl.waitTime = time.Duration(2 * rl.waitTime)
			if rl.waitTime > rl.maxWaitTime {
				rl.waitTime = rl.maxWaitTime
			}

			remainingTime := currentWaitTime - time.Since(lastTime)
			c.HTML(http.StatusTooManyRequests, "limiter.html", gin.H{"RemainingTime": remainingTime.String()})
			c.Abort()
			return
		}

		rl.lastRequest[clientIP] = time.Now()
		c.Next()

		rl.waitTime = 10 * time.Second
	}
}

func IsWhitespace(s string) bool {
	return strings.TrimSpace(s) == ""
}
