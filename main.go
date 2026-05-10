package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/yoshioka0101/go-todo-sample/internal/middleware"
)

func main() {
	// slog を JSON 形式で出力するよう初期化する
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	// gin.Default() はテキストロガーを含むので、gin.New() で素の状態から組み立てる
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Logger()) // slog でリクエストを JSON 出力するミドルウェア

	// ルーティング定義
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// gin.Engine を http.Server に渡す
	// graceful shutdown のために r.Run() は使わない
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// ListenAndServe はブロッキングなので goroutine で起動し、メインの処理を止めない
	go func() {
		slog.Info("server started", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// ErrServerClosed は Shutdown() による正常終了なので無視する
			slog.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	// SIGINT（Ctrl+C）または SIGTERM を受け取るまでここでブロック
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down...")

	// 5秒以内に処理中のリクエストが完了しなければ強制終了
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("forced shutdown", "err", err)
		os.Exit(1)
	}

	slog.Info("server stopped")
}
