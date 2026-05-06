package middleware

import (
	"api-gateway/internal/logger"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// LoggerMiddleware Gin日志中间件
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 生成trace_id
		traceID := uuid.New().String()
		c.Set("trace_id", traceID)

		// 记录请求开始时间
		startTime := time.Now()

		// 记录请求信息
		logger.Info("Incoming request",
			zap.String("trace_id", traceID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("client_ip", c.ClientIP()),
		)

		// 处理请求
		c.Next()

		// 计算请求耗时
		duration := time.Since(startTime)

		// 记录响应信息
		if c.Writer.Status() >= 400 {
			logger.Error("Request completed with error",
				zap.String("trace_id", traceID),
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.Int("status", c.Writer.Status()),
				zap.Duration("duration", duration),
			)
		} else {
			logger.Info("Request completed",
				zap.String("trace_id", traceID),
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.Int("status", c.Writer.Status()),
				zap.Duration("duration", duration),
			)
		}
	}
}
