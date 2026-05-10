package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger はリクエストの情報を slog で JSON 出力する gin ミドルウェア。
// gin.Default() に含まれるデフォルトのテキストロガーの代わりに使う。
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 後続のハンドラを実行
		c.Next()

		slog.Info("request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency", time.Since(start).String(),
			"ip", c.ClientIP(),
		)
	}
}
