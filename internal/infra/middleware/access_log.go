package middleware

import (
	"econode-cloud/internal/pkg/ctxx"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func AccessLog(baseLogger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 生成Trace ID
		traceID := uuid.New().String()
		requestID := uuid.New().String()

		// 创建请求专用的 AccessLog
		requestLogger := baseLogger.With(
			zap.String("trace_id", traceID),
			zap.String("request_id", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
		)
		// 设置到 Context
		reqCtx := c.Request.Context()
		reqCtx = ctxx.WithTraceID(reqCtx, traceID)
		reqCtx = ctxx.WithRequestID(reqCtx, requestID)
		reqCtx = ctxx.WithLogger(reqCtx, requestLogger)

		// 更新请求 Context
		c.Request = c.Request.WithContext(reqCtx)

		// 处理请求
		c.Next()

		// 记录访问日志
		latency := time.Since(start)
		requestLogger.Info("request completed",
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", latency),
			zap.Int("body_size", c.Writer.Size()),
		)
	}
}

//
//func GenerateTraceID() string {
//	// 1. 时间戳（8位十六进制，精确到微秒）
//	timestamp := fmt.Sprintf("%08x", time.Now().UnixMicro())
//
//	// 2. 机器标识（4位）
//	machineID := getMachineID()
//
//	// 3. 进程ID（4位）
//	processID := fmt.Sprintf("%04x", os.Getpid()%0xffff)
//
//	// 4. 序列号（4位，每毫秒重置）
//	seq := getSequence()
//
//	// 5. 随机数（8位，防止冲突）
//	random := fmt.Sprintf("%08x", rand.Uint32())
//
//	// 组合
//	return fmt.Sprintf("%s%s%s%s%s",
//		timestamp, machineID, processID, seq, random)
//}
//
//func GenerateRequestID(serviceName string) string {
//	// 1. 服务名缩写（清晰标识）
//	servicePrefix := getServicePrefix(serviceName)
//
//	// 2. 时间戳（可读格式）
//	timestamp := time.Now().Format("20060102150405.000")
//	timestamp = strings.Replace(timestamp, ".", "", 1)
//
//	// 3. 序列号（3位，每日重置）
//	seq := getDailySequence(serviceName)
//
//	return fmt.Sprintf("%s-%s-%03d", servicePrefix, timestamp, seq)
//}
